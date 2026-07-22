package part

// Service бизнес-логика каталога деталей.
type Service struct {
	repo                 PartRepository
	tx                   TxManager
	compatibilityChecker CompatibilityChecker
}

// NewService создаёт сервис каталога.
func NewService(repo PartRepository, tx TxManager, compatibilityChecker CompatibilityChecker) *Service {
	return &Service{
		repo:                 repo,
		tx:                   tx,
		compatibilityChecker: compatibilityChecker,
	}
}
