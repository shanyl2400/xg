package entity

const (
	AuthEnterStudent      = 1
	AuthDispatchSelfOrder = 2
	AuthDispatchOrder     = 3
	AuthCheckOrder        = 4
	AuthListAllOrder      = 5
	AuthListOrgOrder      = 6
	AuthManageOrderSource = 7
	AuthManageSubject     = 8
	AuthManageOrg         = 9
	AuthCheckOrg          = 10
	AuthManageUser        = 11
	AuthManageRole        = 12
	AuthManageSelfOrg     = 13

	AuthManageSettlement = 14

	AuthModeRootAuth = 1
	AuthModeOrgAuth  = 2
)

type Auth struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Mode int    `json:"mode"`
}
