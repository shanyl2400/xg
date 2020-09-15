package entity

import (
	"strconv"
	"strings"
	"xg/log"
)

const (
	OrgStatusCreated = iota + 1
	OrgStatusCertified
	OrgStatusRejected
	OrgStatusRevoked

	RootOrgId = 1
)

type Org struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Subjects  []string `json:"subjects"`
	Address   string   `json:"address"`
	ParentID  int      `json:"parent_id"`
	Telephone string   `json:"telephone"`
	SupportRoleID []int `json:"support_role_id"`

	Status int `json:"status"`

	SubOrgs []*Org `json:"sub_orgs"`
}

type CreateOrgRequest struct {
	Name      string   `json:"name"`
	Subjects  []string `json:"subjects"`
	Address   string   `json:"address"`
	Telephone string   `json:"telephone"`
	SupportRoleID	[]int `json:"support_role_id"`

	Status   int `json:"status"`
	ParentID int `json:"parent_id"`
}

type CreateOrgWithSubOrgsRequest struct {
	OrgData CreateOrgRequest    `json:"org"`
	SubOrgs []*CreateOrgRequest `json:"sub_orgs"`
}

type UpdateOrgRequest struct {
	ID        int      `json:"id"`
	Subjects  []string `json:"subjects"`
	Address   string   `json:"address"`
	Telephone string   `json:"telephone"`

	Status int `json:"status"`
}

func IntArrayToString(a []int) string{
	if len(a) < 1 {
		return ""
	}
	ids := make([]string, len(a))
	for i := range a {
		ids[i] = strconv.Itoa(a[i])
	}
	return strings.Join(ids, ",")
}
func StringToIntArray(s string) []int {
	if s == "" {
		return nil
	}
	ret := make([]int, 0)
	parts := strings.Split(s, ",")
	for i := range parts {
		id, err := strconv.Atoi(parts[i])
		if err != nil{
			log.Warning.Println("Can't convert ids, str: ", s, ", part: ", parts[i])
			continue
		}
		ret = append(ret, id)
	}
	return ret
}