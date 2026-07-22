package assembly

// Service — application-сервис сборки корабля.
type Service struct {
	cfg AssemblerConfig
}

// NewService создаёт сервис сборки.
func NewService(cfg AssemblerConfig) *Service {
	return &Service{cfg: cfg}
}
