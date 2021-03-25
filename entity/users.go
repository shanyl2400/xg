package entity

type JWTUser struct {
	UserId int
	OrgId  int
	RoleId int
}

type UserLoginResponse struct {
	UserDetailsInfo
	Token string `json:"token"`
}
type UserDetailsInfo struct {
	UserId   int     `json:"user_id"`
	RoleId   int     `json:"role_id"`
	OrgId    int     `json:"org_id"`
	RoleName string  `json:"role_name"`
	OrgName  string  `json:"org_name"`
	Auths    []*Auth `json:"auths"`
	Avatar   string  `json:"avatar"`
}

type UserInfo struct {
	UserId   int    `json:"user_id"`
	RoleId   int    `json:"role_id"`
	OrgId    int    `json:"org_id"`
	RoleName string `json:"role_name"`
	OrgName  string `json:"org_name"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

type UserLoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Name   string `json:"name"`
	OrgId  int    `json:"org_id"`
	RoleId int    `json:"role_id"`
}

type UserUpdatePasswordRequest struct {
	NewPassword string `json:"new_password"`
}
