package entity

import "time"

const (
	StudentConflictStatusUnprocessed = 1
	StudentConflictStatusProcessed   = 2
)

type CreateStudentConflictRequest struct {
	Telephone string `json:"telephone"`
	AuthorID  int    `json:"author_id"`
	Total     int    `json:"total"`
}

type HandleStudentConflictRequest struct {
	RecordID        int `json:"record_id"`
	SelectStudentID int `json:"select_student_id"`
}
type HandleStudentConflictStatusRequest struct {
	RecordID int `json:"record_id"`
	Status   int `json:"status"`
}
type HandleUpdateConflictStudentStatusRequest struct {
	StudentID int `json:"student_id"`
	Status    int `json:"status"`
}

type StudentConflictRecord struct {
	ID        int        `json:"id"`
	Status    int        `json:"status"`
	Telephone string     `json:"telephone"`
	AuthorID  int        `json:"author_id"`
	Total     int        `json:"total"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type StudentConflictRecordDetails struct {
	ID         int        `json:"id"`
	Status     int        `json:"status"`
	Telephone  string     `json:"telephone"`
	AuthorID   int        `json:"author_id"`
	Total      int        `json:"total"`
	UpdatedAt  *time.Time `json:"updated_at"`
	AuthorName string     `json:"author_name"`
}

type StudentConflictsInfoList struct {
	Total   int                             `json:"total"`
	Records []*StudentConflictRecordDetails `json:"records"`
}
