package record

import (
	"time"

	"github.com/google/uuid"
)

// Order — запись заказа в PostgreSQL.
type Order struct {
	UUID            uuid.UUID  `db:"uuid"`
	Status          string     `db:"status"`
	TransactionUUID *uuid.UUID `db:"transaction_uuid"`
	PaymentMethod   *string    `db:"payment_method"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
}

// OrderItem — запись позиции заказа в PostgreSQL.
type OrderItem struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
	PartType  string    `db:"part_type"`
	Price     int64     `db:"price"`
}
