package model

// Part — деталь из каталога inventory, нужная для расчёта заказа.
type Part struct {
	UUID          string
	Price         int64
	StockQuantity int64
}
