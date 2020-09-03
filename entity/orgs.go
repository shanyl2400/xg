package entity

const (
	OrgStatusCreated = iota + 1
	OrgStatusCertified
	OrgStatusRevoked

	RootOrgId = 1
)

type Org struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Subjects []string `json:"subjects"`
	Address string `json:"address"`
	ParentID int `json:"parent_id"`

	Status int `json:"status"`

	SubOrgs []*Org `json:"sub_orgs"`
}

type CreateOrgRequest struct {
	Name     string   `json:"name"`
	Subjects []string `json:"subjects"`
	Address string `json:"address"`

	Status   int `json:"status"`
	ParentID int `json:"parent_id"`
}

type UpdateOrgRequest struct {
	ID       int      `json:"id"`
	Subjects []string `json:"subjects"`
	Address string `json:"address"`

	Status int `json:"status"`
}
