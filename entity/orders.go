package entity

import "time"

const (
	OrderStatusCreated  = 1
	OrderStatusSigned   = 2
	OrderStatusRevoked  = 3
	OrderStatusInvalid  = 4
	OrderStatusDeposit  = 5
	OrderStatusConsider = 6

	OrderPayStatusPending  = 1
	OrderPayStatusChecked  = 2
	OrderPayStatusRejected = 3

	OrderPayModePay     = 1
	OrderPayModePayback = 2

	OrderRemarkModeServer = 1
	OrderRemarkModeClient = 2

	OrderRemarkInfoTypeNormal   = 1
	OrderRemarkInfoTypeSignup   = 2
	OrderRemarkInfoTypeDeposit  = 3
	OrderRemarkInfoTypeRevoke   = 4
	OrderRemarkInfoTypeInvalid  = 5
	OrderRemarkInfoTypeConsider = 6
	OrderRemarkInfoTypePay      = 6
	OrderRemarkInfoTypePayback  = 6

	OrderRemarkUnread = 1
	OrderRemarkRead   = 2

	OrderNotifyStatusUnread = 1
	OrderNotifyStatusRead   = 2

	OrderNotifyClassifySignup  = 1
	OrderNotifyClassifyDeposit = 2
	OrderNotifyClassifyRevoke  = 3
	OrderNotifyClassifyInvalid = 4
)

type CreateOrderRequest struct {
	StudentID      int      `json:"student_id"`
	ToOrgID        int      `json:"to_org_id"`
	IntentSubjects []string `json:"intent_subjects"`
}

type SearchOrderCondition struct {
	IDs              []int  `json:"i_ds"`
	StudentIDList    []int  `json:"student_id_list"`
	ToOrgIDList      []int  `json:"to_org_id_list"`
	IntentSubjects   string `json:"intent_subjects"`
	PublisherID      []int  `json:"publisher_id"`
	OrderSourceList  []int  `json:"order_source"`
	StudentsKeywords string `json:"students_keywords"`
	Keywords         string `json:"keywords"`
	Address          string `json:"address"`
	AuthorID         []int  `json:"author_id"`

	CreateStartAt *time.Time `json:"create_start_at"`
	CreateEndAt   *time.Time `json:"create_end_at"`
	UpdateStartAt *time.Time `json:"update_start_at"`
	UpdateEndAt   *time.Time `json:"update_end_at"`

	Status  []int  `json:"status"`
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

type OrderRemarkList struct {
	Total   int                  `json:"total"`
	Records []*OrderRemarkRecord `json:"records"`
}

type PayRecordInfoList struct {
	Total   int              `json:"total"`
	Records []*PayRecordInfo `json:"records"`
}

type PayRecordInfo struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	Mode    int     `json:"mode"`
	Title   string  `json:"title"`
	Amount  float64 `json:"amount"`
	Content string  `json:"content"`

	StudentID     int      `json:"student_id"`
	ToOrgID       int      `json:"to_org_id"`
	IntentSubject []string `json:"intent_subject"`
	PublisherID   int      `json:"publisher_id"`
	StudentName   string   `json:"student_name"`
	OrgName       string   `json:"org_name"`
	PublisherName string   `json:"publisher_name"`

	Status int `json:"status"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type OrderInfo struct {
	ID            int      `json:"id"`
	StudentID     int      `json:"student_id"`
	ToOrgID       int      `json:"to_org_id"`
	IntentSubject []string `json:"intent_subject"`
	PublisherID   int      `json:"publisher_id"`
	OrderSource   int      `json:"order_source"`

	Status int `json:"status"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type OrderInfoDetails struct {
	OrderInfo
	StudentName      string `json:"student_name"`
	StudentTelephone string `json:"student_telephone"`
	OrgName          string `json:"org_name"`
	PublisherName    string `json:"publisher_name"`

	AuthorID   int    `json:"author_id"`
	AuthorName string `json:"author_name"`

	OrderSourceName string `json:"order_source_name"`
}

type OrderInfoWithRecords struct {
	OrderInfo
	StudentSummary *StudentSummaryInfo `json:"student_summary"`
	OrgName        string              `json:"org_name"`
	PublisherName  string              `json:"publisher_name"`
	AuthorName     string              `json:"author_name"`
	PaymentInfo    []*OrderPayRecord
	RemarkInfo     []*OrderRemarkRecord

	OrderSourceName string `json:"order_source_name"`
}

type OrderPayRecord struct {
	ID      int     `json:"id"`
	OrderID int     `json:"order_id"`
	Mode    int     `json:"mode"`
	Amount  float64 `json:"amount"`
	Title   string  `json:"title"`
	Content string  `json:"content"`

	Status int `json:"status"`

	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}

type OrderNotify struct {
	ID        int               `json:"id"`
	OrderID   int               `json:"order_id"`
	Classify  int               `json:"classify"`
	Content   string            `json:"content"`
	Author    int               `json:"author"`
	OrderInfo *OrderInfoDetails `json:"order_info"`

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
	Status  int    `json:"status"`

	Info     string `json:"info"`
	InfoType int    `json:"info_type"`

	UpdatedAt *time.Time `json:"updated_at"`
	CreatedAt *time.Time `json:"created_at"`
}

type OrderRemarkRequest struct {
	OrderID  int    `json:"order_id"`
	InfoType int    `json:"info_type"`
	Info     string `json:"info"`
	Content  string `json:"content"`
}

type OrderPayRequest struct {
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
}

type UpdateOrderPayAmountRequest struct {
	Record int     `json:"id"`
	Amount float64 `json:"amount"`
}

type OrderUpdateStatusRequest struct {
	Status  int     `json:"status"`
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
}

type OrderMarkRequest struct {
	OrderID int    `json:"order_id"`
	Content string `json:"content"`
}

type MarkOrderRemarkRequest struct {
	Status int   `json:"status"`
	IDs    []int `json:"ids"`
}
