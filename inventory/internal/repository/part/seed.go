package part

import (
	"time"

	"github.com/google/uuid"

	"github.com/Anton119/rocket-service-/inventory/internal/model"
)

// SeedParts возвращает предзаполненный каталог деталей для cmd и тестов.
func SeedParts() map[uuid.UUID]model.Part {
	now := time.Now()

	return map[uuid.UUID]model.Part{
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Алюминиевый корпус",
			Description:   "Лёгкий корпус для небольших кораблей",
			Price:         500000,
			PartType:      model.PartTypeHull,
			StockQuantity: 10,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Титановый корпус",
			Description:   "Прочный корпус для средних кораблей",
			Price:         1500000,
			PartType:      model.PartTypeHull,
			StockQuantity: 5,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Ионный двигатель C",
			Description:   "Базовый ионный двигатель класса C",
			Price:         300000,
			PartType:      model.PartTypeEngine,
			StockQuantity: 8,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Ионный двигатель B",
			Description:   "Улучшенный ионный двигатель класса B",
			Price:         800000,
			PartType:      model.PartTypeEngine,
			StockQuantity: 3,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Энергетический щит",
			Description:   "Стандартный энергетический щит",
			Price:         400000,
			PartType:      model.PartTypeShield,
			StockQuantity: 6,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440006",
			Name:          "Лазерная пушка",
			Description:   "Точная лазерная пушка",
			Price:         250000,
			PartType:      model.PartTypeWeapon,
			StockQuantity: 7,
			CreatedAt:     now,
		},
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
			UUID:          "550e8400-e29b-41d4-a716-446655440007",
			Name:          "Корпус без наличия",
			Description:   "Дорогой корпус, нет на складе",
			Price:         2000000,
			PartType:      model.PartTypeHull,
			StockQuantity: 0,
			CreatedAt:     now,
		},
	}
}
