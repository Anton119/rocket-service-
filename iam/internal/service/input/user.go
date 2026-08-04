package input

// RegisterInput — вход use case регистрации пользователя.
type RegisterInput struct {
	Login    string
	Email    string
	Password string
}

// LoginInput — вход use case входа в систему.
type LoginInput struct {
	Login    string
	Password string
}
