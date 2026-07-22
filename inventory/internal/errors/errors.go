package errs

import "errors"

var (
	// ErrPartNotFound возвращается, когда деталь с указанным UUID отсутствует.
	ErrPartNotFound = errors.New("деталь не найдена")
	// ErrEmptyUUID — пустой uuid в запросе.
	ErrEmptyUUID = errors.New("uuid не может быть пустым")
	// ErrInvalidUUID — uuid не является UUID.
	ErrInvalidUUID = errors.New("неверный формат uuid")
	// ErrOutOfStock — нет доступного остатка (stock - reserved < 1).
	ErrOutOfStock = errors.New("деталь отсутствует на складе")
	// ErrNothingToCommit — нет резерва для списания (reserved = 0).
	ErrNothingToCommit = errors.New("нечего списывать")
	// ErrNothingToRelease — нет резерва для освобождения (reserved = 0).
	ErrNothingToRelease = errors.New("нечего освобождать")
	// ErrIncompatibleParts — детали несовместимы по бизнес-правилам.
	ErrIncompatibleParts = errors.New("детали несовместимы")
	// ErrPartTypeMismatch — тип детали не соответствует ожидаемому набору.
	ErrPartTypeMismatch = errors.New("тип детали не соответствует слоту корабля")
	// ErrInvalidProperties — некорректные свойства детали.
	ErrInvalidProperties = errors.New("некорректные свойства детали")
)

// IsInvalidUUID возвращает true для ошибок разбора UUID.
func IsInvalidUUID(err error) bool {
	return errors.Is(err, ErrInvalidUUID) || errors.Is(err, ErrEmptyUUID)
}
