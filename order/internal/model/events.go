package model

import "time"

// OrderPaidEvent — событие оплаты заказа для Kafka.
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	TransactionUUID string
	PaymentMethod   string
	UserUUID        string
	PaidAt          time.Time
}

// ShipAssembledEvent — событие завершения сборки из Kafka.
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
