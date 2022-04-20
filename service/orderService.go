package service

type OrderCreate struct {
	Timestamp   string  `json:"timestamp"`
	Operation   string  `json:"operation"`
	IssuerName  string  `json:"issuer_name"`
	TotalShares int     `json:"total_shares"`
	SharePrice  float32 `json:"total_price"`
}

type Issuer struct {
	IssuerName     string   `json:"issuer_name"`
	TotalShares    int      `json:"total_shares"`
	SharePrice     float32  `json:"share_price"`
	BusinessErrors []string `json:"business_erros"`
}

type CurrentBalance struct {
	Cash    float32  `json:"cash"`
	Issuers []Issuer `json:"issuers"`
}

type OrderCreateRequest struct {
	AccountID int           `json:"account_id"`
	Orders    []OrderCreate `json:"orders"`
}

type OrderCreateResponse struct {
	CBalance CurrentBalance
}

type OrderService interface {
	CreateOrder(*OrderCreateRequest) *OrderCreateResponse
}
