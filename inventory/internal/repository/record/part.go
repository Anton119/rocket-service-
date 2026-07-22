package record

import (
	"time"

	"github.com/google/uuid"
)

// Part — запись детали в PostgreSQL.
type Part struct {
	UUID          uuid.UUID  `db:"uuid"`
	Name          string     `db:"name"`
	Description   string     `db:"description"`
	PartType      string     `db:"part_type"`
	Price         int64      `db:"price"`
	StockQuantity int32      `db:"stock_quantity"`
	Reserved      int32      `db:"reserved"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}
