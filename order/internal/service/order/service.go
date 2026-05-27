package order

// Service инкапсулирует бизнес-логику заказов.
type Service struct {
	repo OrderRepository
	inv  InventoryClient
	pay  PaymentClient
}

// NewService создаёт сервис заказов.
func NewService(repo OrderRepository, inv InventoryClient, pay PaymentClient) *Service {
	return &Service{
		repo: repo,
		inv:  inv,
		pay:  pay,
	}
}
