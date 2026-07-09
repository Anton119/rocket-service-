package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus — статус заказа в доменной модели (строки совпадают с OpenAPI enum).
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// OrderItem — позиция заказа (snapshot детали на момент создания).
type OrderItem struct {
	PartUUID uuid.UUID
	PartType string
	Price    int64
}

// Order — заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	Items           []OrderItem
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID
	WeaponUUID      *uuid.UUID
	TotalPrice      int64
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
