package entity

import "time"

type CreateSettlementRequest struct {
	ID            int       `json:"id"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"ent_at"`
	SuccessOrders []string  `json:"success_orders"`
	FailedOrders  []string  `json:"failed_orders"`
	Amount        float64   `json:"amount"`
	Status        int       `json:"status"`
	Invoice       int       `json:"invoice"`
}

type SettlementData struct {
	ID            int        `json:"id"`
	StartAt       time.Time  `json:"start_at"`
	EndAt         time.Time  `json:"ent_at"`
	SuccessOrders []string   `json:"success_orders"`
	FailedOrders  []string   `json:"failed_orders"`
	Amount        float64    `json:"amount"`
	Status        int        `json:"status"`
	Invoice       int        `json:"invoice"`
	AuthorID      int        `json:"author_id"`
	AuthorName    string     `json:"author_name"`
	UpdatedAt     *time.Time `json:"updated_at"`
	CreatedAt     *time.Time `json:"created_at"`
}
