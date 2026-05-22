package errs

import "errors"

var (
	// ErrPartNotFound возвращается, когда деталь с указанным UUID отсутствует.
	ErrPartNotFound = errors.New("деталь не найдена")
	// ErrEmptyUUID — пустой uuid в запросе.
	ErrEmptyUUID = errors.New("uuid не может быть пустым")
	// ErrInvalidUUID — uuid не является UUID.
	ErrInvalidUUID = errors.New("неверный формат uuid")
)

// IsInvalidUUID возвращает true для ошибок разбора UUID.
func IsInvalidUUID(err error) bool {
	return errors.Is(err, ErrInvalidUUID) || errors.Is(err, ErrEmptyUUID)
}
