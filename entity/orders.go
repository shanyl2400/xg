package entity

import "time"

const (
	OrderStatusCreated = iota + 1
	OrderStatusSigned
	OrderStatusRevoked

	OrderPayStatusPending  = 1
	OrderPayStatusChecked  = 2
	OrderPayStatusRejected = 3

	OrderPayModePay     = 1
	OrderPayModePayback = 2

	OrderRemarkModeServer = 1
	OrderRemarkModeClient = 2
)

type CreateOrderRequest struct {
	StudentID      int      `json:"student_id"`
	ToOrgID        int      `json:"to_org_id"`
	IntentSubjects []string `json:"intent_subjects"`
}

type SearchOrderCondition struct {
	StudentIDList  []int  `json:"student_id_list"`
	ToOrgIDList    []int  `json:"to_org_id_list"`
	IntentSubjects string `json:"intent_subjects"`
	PublisherID    int    `json:"publisher_id"`

	Status  []int    `json:"status"`
	OrderBy string `json:"order_by"`

	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

type SearchPayRecordCondition struct {
	PayRecordIDList []int `json:"pay_record_id_list"`
	OrderIDList     []int `json:"order_id_list"`
	AuthorIDList    []int `json:"author_id_list"`
	Mode            int   `json:"mode"`
	StatusList      []int `json:"status_list"`

	OrderBy string `json:"order_by"`

	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

type OrderInfoList struct {
	Total  int                 `json:"total"`
	Orders []*OrderInfoDetails `json:"orders"`
}

type PayRecordInfoList struct {
	Total   int              `json:"total"`
	Records []*PayRecordInfo `json:"records"`
}

type PayRecordInfo struct {
	ID      int    `json:"id"`
	OrderID int    `json:"order_id"`
	Mode    int    `json:"mode"`
	Title   string `json:"title"`
	Amount  int    `json:"amount"`

	Status int `json:"status"`
}

type OrderInfo struct {
	ID            int      `json:"id"`
	StudentID     int      `json:"student_id"`
	ToOrgID       int      `json:"to_org_id"`
	IntentSubject []string `json:"intent_subject"`
	PublisherID   int      `json:"publisher_id"`

	Status int `json:"status"`
}

type OrderInfoDetails struct {
	OrderInfo
	StudentName      string `json:"student_name"`
	StudentTelephone string `json:"student_telephone"`
	OrgName          string `json:"org_name"`
	PublisherName    string `json:"publisher_name"`
}

type OrderInfoWithRecords struct {
	OrderInfo
	StudentSummary *StudentSummaryInfo `json:"student_summary"`
	OrgName        string              `json:"org_name"`
	PublisherName  string              `json:"publisher_name"`
	PaymentInfo    []*OrderPayRecord
	RemarkInfo     []*OrderRemarkRecord
}

type OrderPayRecord struct {
	ID      int `json:"id"`
	OrderID int `json:"order_id"`
	Mode    int `json:"mode"`
	Amount  int `json:"amount"`

	Status int `json:"status"`

	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}

type OrderRemarkRecord struct {
	ID      int    `json:"id"`
	OrderID int    `json:"order_id"`
	Author  int    `json:"author"`
	Mode    int    `json:"mode"`
	Content string `json:"content"`

	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}

type OrderPayRequest struct {
	OrderID int    `json:"order_id"`
	Amount  int    `json:"amount"`
	Title   string `json:"title"`
}

type OrderMarkRequest struct {
	OrderID int `json:"order_id"`
	Content string `json:"content"`
}