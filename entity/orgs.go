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
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Subjects      []string `json:"subjects"`
	Address       string   `json:"address"`
	AddressExt    string   `json:"address_ext"`
	ParentID      int      `json:"parent_id"`
	Telephone     string   `json:"telephone"`
	SupportRoleID []int    `json:"support_role_id"`

	Status int `json:"status"`

	SubOrgs []*Org `json:"sub_orgs"`
}

type SubOrgWithDistance struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Subjects      []string `json:"subjects"`
	Address       string   `json:"address"`
	AddressExt    string   `json:"address_ext"`
	ParentID      int      `json:"parent_id"`
	Telephone     string   `json:"telephone"`
	SupportRoleID []int    `json:"support_role_id"`

	Status   int     `json:"status"`
	Distance float64 `json:"distance"`
}

type CreateOrgRequest struct {
	Name          string   `json:"name"`
	Subjects      []string `json:"subjects"`
	Address       string   `json:"address"`
	AddressExt    string   `json:"address_ext"`
	Telephone     string   `json:"telephone"`
	SupportRoleID []int    `json:"support_role_id"`
	Longitude     float64  `json:"longitude"`
	Latitude      float64  `json:"latitude"`

	Status   int `json:"status"`
	ParentID int `json:"parent_id"`
}

type CreateOrgWithSubOrgsRequest struct {
	OrgData CreateOrgRequest    `json:"org"`
	SubOrgs []*CreateOrgRequest `json:"sub_orgs"`
}

type UpdateOrgWithSubOrgsRequest struct {
	OrgData CreateOrUpdateOrgRequest    `json:"org"`
	SubOrgs []*CreateOrUpdateOrgRequest `json:"sub_orgs"`
}

type CreateOrUpdateOrgRequest struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Subjects   []string `json:"subjects"`
	Address    string   `json:"address"`
	AddressExt string   `json:"address_ext"`
	Telephone  string   `json:"telephone"`

	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type UpdateSubOrgsEntity struct {
	UpdateOrgReq   CreateOrUpdateOrgRequest    `json:"update_org_req"`
	InsertOrgList  []*CreateOrUpdateOrgRequest `json:"insert_org_list"`
	UpdateOrgsList []*CreateOrUpdateOrgRequest `json:"update_orgs_list"`
	DeletedIds     []int                       `json:"deleted_ids"`
	OrgInfo        *Org                        `json:"org_info"`
}

type UpdateOrgRequest struct {
	ID         int      `json:"id"`
	Name string `json:"name"`
	Subjects   []string `json:"subjects"`
	Address    string   `json:"address"`
	AddressExt string   `json:"address_ext"`
	Telephone  string   `json:"telephone"`
	Longitude  float64  `json:"longitude"`
	Latitude   float64  `json:"latitude"`

	Status int `json:"status"`
}

func IntArrayToString(a []int) string {
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
		if err != nil {
			log.Warning.Println("Can't convert ids, str: ", s, ", part: ", parts[i])
			continue
		}
		ret = append(ret, id)
	}
	return ret
}

type Coordinate struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
