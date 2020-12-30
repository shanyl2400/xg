package entity

import "time"

const (
	StudentCreated = iota + 1
	StudentConflictFailed
	StudentConflictSuccess
	StudentExceed

	StudentApply
	StudentApplyWithFee
)

type CreateStudentRequest struct {
	Name          string   `json:"name"`
	Gender        bool     `json:"gender"`
	Telephone     string   `json:"telephone"`
	Address       string   `json:"address"`
	AddressExt    string   `json:"address_ext"`
	Email         string   `json:"email"`
	IntentSubject []string `json:"intent_subject"`
	Note          string   `json:"note"`
	OrderSourceID int      `json:"order_source_id"`

	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type CreateStudentResponse struct {
	ID     int `json:"id"`
	Status int `json:"status"`
}

type UpdateStudentRequest struct {
	Name          string   `json:"name"`
	Gender        bool     `json:"gender"`
	Telephone     string   `json:"telephone"`
	Email         string   `json:"email"`
	Address       string   `json:"address"`
	AddressExt    string   `json:"address_ext"`
	IntentSubject []string `json:"intent_subject"`
	OrderSourceID int      `json:"order_source_id"`

	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type SearchStudentRequest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Address   string `json:"address"`
	Keywords string `json:"keywords"`

	IntentSubject string `json:"intent_subject"`
	AuthorIDList  []int  `json:"author_id_list"`
	Status        []int    `json:"status"`

	NoDispatchOrder bool `json:"no_dispatch_order"`

	OrderBy  string `json:"order_by"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

type AddStudentNoteRequest struct {
	ID   int    `json:"id"`
	Note string `json:"note"`
}

type StudentInfo struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Gender          bool       `json:"gender"`
	Telephone       string     `json:"telephone"`
	Address         string     `json:"address"`
	AddressExt      string     `json:"address_ext"`
	Email           string     `json:"email"`
	IntentSubject   []string   `json:"intent_subject"`
	AuthorID        int        `json:"authorID"`
	AuthorName      string     `json:"authorName"`
	Status          int        `json:"status"`
	OrderCount      int        `json:"order_count"`
	Note            string     `json:"note"`
	OrderSourceID   int        `json:"order_source_id"`
	OrderSourceName string     `json:"order_source_name"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type StudentSummaryInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Gender     bool   `json:"gender"`
	Telephone  string `json:"telephone"`
	Address    string `json:"address"`
	Note       string `json:"note"`
	AddressExt string `json:"address_ext"`

	AuthorId  int        `json:"author_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type StudentInfoList struct {
	Total    int            `json:"total"`
	Students []*StudentInfo `json:"students"`
}

type StudentInfosWithOrders struct {
	StudentInfo
	Orders []*OrderInfoDetails `json:"orders"`
}
