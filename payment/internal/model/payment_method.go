package model

// PaymentMethod — способ оплаты в доменной модели (значения совпадают с payment.v1.PaymentMethod).
type PaymentMethod int32

const (
	PaymentMethodUnspecified   PaymentMethod = 0
	PaymentMethodCard          PaymentMethod = 1
	PaymentMethodSBP           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)
