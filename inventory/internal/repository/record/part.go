package record

import "time"

// Part — запись детали в хранилище.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      int32 // persistence: числовой код типа (см. model.PartType)
	StockQuantity int64
	CreatedAt     time.Time
}
