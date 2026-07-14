package model

import (
	"time"

	"github.com/google/uuid"

	errs "github.com/Anton119/rocket-service-/inventory/internal/errors"
)

// Part — доменная сущность детали космического корабля.
type Part struct {
	uuid          uuid.UUID
	name          string
	description   string
	partType      PartType
	price         int64
	stockQuantity int
	reserved      int
	properties    PartProperties
	createdAt     time.Time
}

// RestorePart восстанавливает сущность из БД (без валидации — данные уже проверены).
func RestorePart(
	partUUID uuid.UUID,
	name, description string,
	partType PartType,
	price int64,
	stockQuantity, reserved int,
	properties PartProperties,
	createdAt time.Time,
) Part {
	return Part{
		uuid:          partUUID,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}

// Reserve резервирует одну единицу детали.
func (p *Part) Reserve() error {
	if p.Available() <= 0 {
		return errs.ErrOutOfStock
	}

	p.reserved++

	return nil
}

// Release освобождает одну зарезервированную единицу.
func (p *Part) Release() error {
	if p.reserved <= 0 {
		return errs.ErrNothingToRelease
	}

	p.reserved--

	return nil
}

func (p Part) UUID() uuid.UUID            { return p.uuid }
func (p Part) Name() string               { return p.name }
func (p Part) Description() string        { return p.description }
func (p Part) PartType() PartType         { return p.partType }
func (p Part) Price() int64               { return p.price }
func (p Part) StockQuantity() int         { return p.stockQuantity }
func (p Part) Reserved() int              { return p.reserved }
func (p Part) Properties() PartProperties { return p.properties }
func (p Part) CreatedAt() time.Time       { return p.createdAt }
func (p Part) Available() int             { return p.stockQuantity - p.reserved }
