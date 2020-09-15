package entity

const (
	RoleSuperAdmin = 1
	RoleEnterWorker = 2
	RoleDispatchWorker = 3
	RoleUserManager = 4
	RoleOrgManager = 5
	RoleChecker = 6
	RoleOutOrg = 7
)

type Role struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	AuthList []*Auth `json:"auth_list"`
}

type CreateRoleRequest struct {
	Name    string `json:"name"`
	AuthIds []int  `json:"auth_ids"`
}
