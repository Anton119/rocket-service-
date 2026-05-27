package model

// PaymentMethod — способ оплаты в домене (до вызова payment gRPC).
type PaymentMethod int

const (
	PaymentMethodInvalid PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
