package errs

import "errors"

// IsUserNotFound возвращает true, если ошибка означает отсутствие пользователя.
func IsUserNotFound(err error) bool {
	return errors.Is(err, ErrUserNotFound)
}
