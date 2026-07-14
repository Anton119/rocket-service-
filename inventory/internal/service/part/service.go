package part

// Service — application service каталога деталей.
type Service struct {
	repo                 PartRepository
	txManager            TxManager
	compatibilityChecker CompatibilityChecker
}

// NewService создаёт сервис каталога.
func NewService(
	repo PartRepository,
	txManager TxManager,
	compatibilityChecker CompatibilityChecker,
) *Service {
	return &Service{
		repo:                 repo,
		txManager:            txManager,
		compatibilityChecker: compatibilityChecker,
	}
}
