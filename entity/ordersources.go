package entity

const (
	OtherOrderSource = 11
)

type OrderSource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateOrderSourceRequest struct {
	Name string `json:"name"`
}
