package model

import "time"

// OrderPaidEvent — доменное представление события оплаты заказа.
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	TransactionUUID string
	PaymentMethod   string
	UserUUID        string
	PaidAt          time.Time
}

// ShipAssembledEvent — доменное представление события завершения сборки.
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
