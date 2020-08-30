package entity

type Role struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	AuthList []*Auth `json:"auth_list"`
}

type CreateRoleRequest struct {
	Name    string `json:"name"`
	AuthIds []int  `json:"auth_ids"`
}
