package service

type OrderCreate struct {
	AccountID   int
	Timestamp   string
	Operation   string
	IssuerName  string
	TotalShares int
	SharePrice  float32
}

type Issuer struct {
	IssuerName     string
	TotalShares    int
	SharePrice     float32
	BusinessErrors []string
}

type CurrentBalance struct {
	Cash    float32
	Issuers []Issuer
}

type OrderCreateRequest struct {
	Orders []OrderCreate
}

type OrderCreateResponse struct {
	CBalance CurrentBalance
}

type OrderService interface {
	CreateOrder(*OrderCreateRequest) *OrderCreateResponse
}
