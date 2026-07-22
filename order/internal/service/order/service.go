package order

// Service инкапсулирует бизнес-логику заказов.
type Service struct {
	repo      OrderRepository
	inv       InventoryClient
	pay       PaymentClient
	orderPaid OrderPaidProducer
	txManager TxManager
}

// NewService создаёт сервис заказов.
func NewService(
	repo OrderRepository,
	inv InventoryClient,
	pay PaymentClient,
	orderPaid OrderPaidProducer,
	txManager TxManager,
) *Service {
	return &Service{
		repo:      repo,
		inv:       inv,
		pay:       pay,
		orderPaid: orderPaid,
		txManager: txManager,
	}
}
