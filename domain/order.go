package domain

type OrderIn struct {
	AccountID   int
	Timestamp   string
	Operation   string
	IssuerName  string
	TotalShares int
	SharePrice  float32
}

type Order struct {
	AccountID     int
	OrderID       int
	Timestamp     string
	Operation     string
	IssuerName    string
	TotalShares   int
	SharePrice    float32
	Balance       float32
	BusinessError []string
}

type OrderRepository interface {
	CreateOrder(OrderIn) (*Order, string)
}
