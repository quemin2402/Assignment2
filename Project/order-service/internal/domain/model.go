package domain

type OrderItem struct {
	ProductID string
	Quantity  int32
}

type Order struct {
	ID     string
	Items  []OrderItem
	Status string
}
