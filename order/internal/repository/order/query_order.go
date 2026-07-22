package order

// Общие фрагменты SQL для заказов (checklist: query_order.go).
const (
	orderColumns = `uuid, user_uuid, status, transaction_uuid, payment_method, created_at, updated_at`
	itemColumns  = `order_uuid, part_uuid, part_type, price`
)
