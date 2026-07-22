package model

import "time"

// Part представляет деталь космического корабля.
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // в копейках
	PartType      PartType
	StockQuantity int64
	Reserved      int64
	Properties    PartProperties
	CreatedAt     time.Time
}

// Available — доступный остаток (stock - reserved).
func (p Part) Available() int64 {
	return p.StockQuantity - p.Reserved
}
