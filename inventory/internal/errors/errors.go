package errs

import "errors"

// ErrPartNotFound возвращается репозиторием, когда деталь с указанным UUID отсутствует.
var ErrPartNotFound = errors.New("деталь не найдена")
