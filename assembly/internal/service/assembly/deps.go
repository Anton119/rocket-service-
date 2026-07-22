package assembly

// AssemblerConfig — параметры эмуляции сборки.
type AssemblerConfig interface {
	MinSec() int
	MaxSec() int
}
