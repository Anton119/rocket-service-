package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/Anton119/rocket-service-/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Anton119/rocket-service-/shared/pkg/proto/payment/v1"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// OrderHandler реализует интерфейс orderv1.Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	hull := req.GetHullUUID()
	engine := req.GetEngineUUID()

	uuids := []string{hull.String(), engine.String()}

	// если щит есть добавляем, если нет - игнорируем
	if shield, ok := req.GetShieldUUID().Get(); ok {
		uuids = append(uuids, shield.String())
	}
	// также с оружием
	if weapon, ok := req.WeaponUUID.Get(); ok {
		uuids = append(uuids, weapon.String())
	}

	listResp, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	// обработка вариаций ошибок после grpc вызова
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: "деталь не найдена",
			}, nil
		}
		if ok && st.Code() == codes.InvalidArgument {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: st.Message(),
			}, nil
		}
		return nil, err
	}
	// словарь деталей для быстрого доступа: uuid + вся ифна о детали
	byUUID := make(map[string]*inventoryv1.Part, len(listResp.GetParts()))
	for _, p := range listResp.GetParts() {
		byUUID[p.GetUuid()] = p
	}

	// проверка складских заказов для каждой детали в заказе
	for _, id := range uuids {
		p := byUUID[id]
		// проверяем stock_quantity > 0
		if p.GetStockQuantity() <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "деталь отсутствует на складе",
			}, nil
		}
	}
	// подсчет итоговой суммы заказа в копейках
	var totalSum int64
	for _, id := range uuids {
		totalSum += byUUID[id].GetPrice()
	}

	orderUUID := uuid.New()
	now := time.Now()

	order := Order{
		OrderUUID:  orderUUID,
		HullUUID:   hull,
		EngineUUID: engine,
		TotalPrice: totalSum,
		Status:     string(orderv1.OrderStatusPENDINGPAYMENT),
		CreatedAt:  now,
	}

	// добавляем опциональные поля щита и оружия (если были переданы)
	if shield, ok := req.GetShieldUUID().Get(); ok {
		sh := shield
		order.ShieldUUID = &sh
	}
	if weapon, ok := req.GetWeaponUUID().Get(); ok {
		w := weapon
		order.WeaponUUID = &w
	}

	// потокобезопасная запись в in-memory map
	h.store.mu.Lock()
	h.store.orders[orderUUID] = order
	h.store.mu.Unlock()

	resp := orderv1.CreateOrderResponse{}
	resp.SetOrderUUID(orderUUID)
	resp.SetTotalPrice(totalSum)

	return &resp, nil
}

// toProtoPaymentMethod переводит способ оплаты из HTTP (OpenAPI) в protobuf enum payment‑сервиса.
func toProtoPaymentMethod(m orderv1.PaymentMethod) (paymentv1.PaymentMethod, bool) {
	switch m {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, true
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, true
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, true
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, true
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, false
	}
}

func (h *OrderHandler) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	// Первая секция: читаем заказ из store под mutex и проверяем допустимость оплаты по статусу.
	h.store.mu.Lock()
	order, exists := h.store.orders[params.OrderUUID]
	if !exists {
		h.store.mu.Unlock()
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// Оплата разрешена только в PENDING_PAYMENT; иначе возвращаем конфликт жизненного цикла.
	switch order.Status {
	case "PAID", "CANCELLED":
		h.store.mu.Unlock()
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплата невозможна в текущем статусе",
		}, nil
	case "PENDING_PAYMENT":
		// Статус ожидает оплаты — продолжаем.
	default:
		h.store.mu.Unlock()
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплата невозможна в текущем статусе",
		}, nil
	}
	h.store.mu.Unlock()

	// Приводим payment_method из HTTP к protobuf enum для вызова payment‑сервиса.
	pm, ok := toProtoPaymentMethod(req.GetPaymentMethod())
	if !ok {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неизвестный способ оплаты",
		}, nil
	}

	// В gRPC payment.order_uuid имеет тип string, поэтому конвертируем uuid.UUID из path.
	orderUUID := params.OrderUUID.String()

	// Вызываем payment‑микросервис: списание/фиксация платежа в его контуре.
	payResp, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: pm,
	})
	// Ошибки payment: InvalidArgument мапим в HTTP 400, прочее пробрасываем как внутреннюю ошибку.
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			return &orderv1.PayOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: st.Message(),
			}, nil
		}
		return nil, err
	}

	// Разбираем transaction_uuid из ответа payment для сохранения в модели заказа.
	txID, err := uuid.Parse(payResp.GetTransactionUuid())
	if err != nil {
		return nil, err
	}

	// Сохраняем метод оплаты в том виде, как пришёл в HTTP.
	methodStr := string(req.GetPaymentMethod())

	// Вторая секция: после успешного RPC перепроверяем заказ и атомарно переводим в PAID.
	h.store.mu.Lock()
	stored, stillExists := h.store.orders[params.OrderUUID]
	if !stillExists {
		h.store.mu.Unlock()
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}
	// Пока выполнялся payment, параллельный запрос мог сменить статус — защищаемся от гонки.
	switch stored.Status {
	case "PENDING_PAYMENT":
		stored.Status = "PAID"
		stored.TransactionUUID = &txID
		stored.PaymentMethod = &methodStr
		h.store.orders[params.OrderUUID] = stored
	case "PAID", "CANCELLED":
		h.store.mu.Unlock()
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплата невозможна в текущем статусе",
		}, nil
	default:
		h.store.mu.Unlock()
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплата невозможна в текущем статусе",
		}, nil
	}
	h.store.mu.Unlock()

	// Успешный ответ OpenAPI: отдаём transaction_uuid клиенту.
	out := orderv1.PayOrderResponse{}
	out.SetTransactionUUID(txID)
	return &out, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	_ = ctx

	// Блокируем хранилище: отмена меняет заказ и должна быть атомарна.
	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	order, exists := h.store.orders[params.OrderUUID]
	if !exists {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// Отменить можно только пока заказ ждёт оплаты.
	if order.Status != "PENDING_PAYMENT" {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "отмена невозможна в текущем статусе",
		}, nil
	}

	order.Status = "CANCELLED"
	h.store.orders[params.OrderUUID] = order

	// Успешная отмена по OpenAPI: пустой JSON объект {}.
	return &orderv1.CancelOrderResponse{}, nil
}
