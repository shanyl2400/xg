package entity

import (
	"strconv"
	"strings"
	"time"
	"xg/log"
)

const (
	OrgStatusCreated   = 1
	OrgStatusCertified = 2
	OrgStatusRejected  = 3
	OrgStatusRevoked   = 4
	OrgStatusOverDue   = 5

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

	BusinessLicense       string     `json:"business_license"`
	CorporateIdentity     string     `json:"corporate_identity"`
	SchoolPermission      string     `json:"school_permission"`
	SettlementInstruction string     `json:"settlement_instruction"`
	Extra                 string     `json:"extra"`
	ExpiredAt             *time.Time `json:"expired_at"`
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

	BusinessLicense       string     `json:"business_license"`
	CorporateIdentity     string     `json:"corporate_identity"`
	SchoolPermission      string     `json:"school_permission"`
	SettlementInstruction string     `json:"settlement_instruction"`
	Extra                 string     `json:"extra"`
	ExpiredAt             *time.Time `json:"expired_at"`
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

	BusinessLicense       string `json:"business_license"`
	CorporateIdentity     string `json:"corporate_identity"`
	SchoolPermission      string `json:"school_permission"`
	SettlementInstruction string `json:"settlement_instruction"`

	Extra      string `json:"extra"`
	ValidMonth int    `json:"valid_month"`
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

	BusinessLicense       string `json:"business_license"`
	CorporateIdentity     string `json:"corporate_identity"`
	SchoolPermission      string `json:"school_permission"`
	Extra                 string `json:"extra"`
	SettlementInstruction string `json:"settlement_instruction"`
}

type UpdateSubOrgsEntity struct {
	UpdateOrgReq   CreateOrUpdateOrgRequest    `json:"update_org_req"`
	InsertOrgList  []*CreateOrUpdateOrgRequest `json:"insert_org_list"`
	UpdateOrgsList []*CreateOrUpdateOrgRequest `json:"update_orgs_list"`
	DeletedIds     []int                       `json:"deleted_ids"`
	OrgInfo        *Org                        `json:"org_info"`
}

type UpdateOrgRequest struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	Subjects              []string `json:"subjects"`
	Address               string   `json:"address"`
	AddressExt            string   `json:"address_ext"`
	Telephone             string   `json:"telephone"`
	Longitude             float64  `json:"longitude"`
	Latitude              float64  `json:"latitude"`
	BusinessLicense       string   `json:"business_license"`
	CorporateIdentity     string   `json:"corporate_identity"`
	SchoolPermission      string   `json:"school_permission"`
	SettlementInstruction string   `json:"settlement_instruction"`
	Extra                 string   `json:"extra"`

	Status int `json:"status"`
}

type RenewOrgRequest struct {
	ID         int `json:"id"`
	ValidMonth int `json:"valid_month"`
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
