package entity

import "time"

const (
	CommissionSettlementStatusUncreated = 1
	CommissionSettlementStatusCreated   = 2

	SettlementStatusUnsettled = 1
	SettlementStatusSettled   = 2
)

type CreateSettlementRequest struct {
	StartAt       int     `json:"start_at"`
	EndAt         int     `json:"end_at"`
	SuccessOrders []int   `json:"success_orders"`
	FailedOrders  []int   `json:"failed_orders"`
	Amount        float64 `json:"amount"`
	Status        int     `json:"status"`
	Commission    float64 `json:"commission"`
	Invoice       string  `json:"invoice"`
}
type CreateCommissionSettlementRequest struct {
	PaymentID      int     `json:"payment_id"`
	Commission     float64 `json:"commission"`
	SettlementNote string  `json:"settlement_note"`
	Note           string  `json:"note"`
}

type UpdateCommissionSettlementRequest struct {
	ID             int     `json:"id"`
	Commission     float64 `json:"commission"`
	SettlementNote string  `json:"settlement_note"`
	Note           string  `json:"note"`
}

type SettlementData struct {
	ID            int        `json:"id"`
	StartAt       time.Time  `json:"start_at"`
	EndAt         time.Time  `json:"end_at"`
	SuccessOrders []int      `json:"success_orders"`
	FailedOrders  []int      `json:"failed_orders"`
	Amount        float64    `json:"amount"`
	Status        int        `json:"status"`
	Commission    float64    `json:"commission"`
	Invoice       string     `json:"invoice"`
	AuthorID      int        `json:"author_id"`
	AuthorName    string     `json:"author_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
	CreatedAt     *time.Time `json:"created_at"`
}
type CommissionSettlementData struct {
	ID             int     `json:"id"`
	OrderID        int     `json:"order_id"`
	PaymentID      int     `json:"payment_id"`
	Amount         float64 `json:"amount"`
	Commission     float64 `json:"commission"`
	SettlementNote string  `json:"settlement_note"`
	Status         int     `json:"status"`
	AuthorID       int     `json:"author_id"`
	AuthorName     string  `json:"author_name"`
	Note           string  `json:"note"`

	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}
