package entity

type Subject struct {
	ID       int      `json:"id"`
	Level    int      `json:"level"`
	ParentId int      `json:"parent_id"`
	Name     string   `json:"name"`
	Parent   *Subject `json:"parent"`
}

type CreateSubjectRequest struct {
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
}

type BatchCreateSubjectRequest struct {
	Data []*CreateSubjectRequest `json:"data"`
}

type SubjectTreeNode struct {
	ID       int                `json:"id"`
	Level    int                `json:"level"`
	Name     string             `json:"name"`
	Title    string             `json:"title"`
	Value    string             `json:"value"`
	Key      int                `json:"key"`
	Children []*SubjectTreeNode `json:"children"`
}
