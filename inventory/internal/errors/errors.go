package errs

import "errors"

var (
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrEmptyUUID         = errors.New("uuid не может быть пустым")
	ErrInvalidUUID       = errors.New("неверный формат uuid")
	ErrOutOfStock        = errors.New("деталь отсутствует на складе")
	ErrNothingToRelease  = errors.New("нечего освобождать")
	ErrIncompatibleParts = errors.New("детали несовместимы")
	ErrPartTypeMismatch  = errors.New("тип детали не соответствует слоту корабля")
	ErrInvalidProperties = errors.New("некорректные свойства детали")
)

// IsInvalidUUID возвращает true для ошибок разбора UUID.
func IsInvalidUUID(err error) bool {
	return errors.Is(err, ErrInvalidUUID) || errors.Is(err, ErrEmptyUUID)
}
