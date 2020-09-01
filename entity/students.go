package entity

const (
	StudentCreated = iota + 1
	StudentConflictFailed
	StudentConflictSuccess

	StudentApply
	StudentApplyWithFee
)

type CreateStudentRequest struct {
	Name          string   `json:"name"`
	Gender        bool     `json:"gender"`
	Telephone     string   `json:"telephone"`
	Address       string   `json:"address"`
	Email         string   `json:"email"`
	IntentSubject []string `json:"intent_subject"`
	Note          string   `json:"note"`
	OrderSourceID int      `json:"order_source_id"`
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
	IntentSubject []string `json:"intent_subject"`
	OrderSourceID int      `json:"order_source_id"`
}

type SearchStudentRequest struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Address   string `json:"address"`

	AuthorIDList []int `json:"author_id_list"`

	OrderBy  string `json:"order_by"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

type AddStudentNoteRequest struct {
	ID   int    `json:"id"`
	Note string `json:"note"`
}

type StudentInfo struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Gender          bool     `json:"gender"`
	Telephone       string   `json:"telephone"`
	Address         string   `json:"address"`
	Email           string   `json:"email"`
	IntentSubject   []string `json:"intent_subject"`
	AuthorID        int      `json:"authorID"`
	AuthorName      string   `json:"authorName"`
	Status          int      `json:"status"`
	Note            string   `json:"note"`
	OrderSourceID   int      `json:"order_source_id"`
	OrderSourceName string   `json:"order_source_name"`
}

type StudentSummaryInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Gender    bool   `json:"gender"`
	Telephone string `json:"telephone"`
	Address   string `json:"address"`
	Note      string `json:"note"`
}

type StudentInfoList struct {
	Total    int            `json:"total"`
	Students []*StudentInfo `json:"students"`
}

type StudentInfosWithOrders struct {
	StudentInfo
	Orders []*OrderInfoDetails `json:"orders"`
}