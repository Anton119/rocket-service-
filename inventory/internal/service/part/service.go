package part

// Service бизнес-логика каталога деталей.
type Service struct {
	repo PartRepository
}

// NewService создаёт сервис каталога.
func NewService(repo PartRepository) *Service {
	return &Service{repo: repo}
}
