package entity

type Subject struct {
	ID       int    `json:"id"`
	Level    int    `json:"level"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
}

type CreateSubjectRequest struct {
	Level    int    `json:"level"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
}
