package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"xg/entity"
)

func TestUpdateOrg(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	createOrgRes, err := client.CreateOrg(ctx, orgs[0], superToken)
	if !assert.NoError(t, err) {
		return
	}

	id := createOrgRes.ID
	t.Log(id)

	err = client.UpdateOrgById(ctx, id, entity.UpdateOrgWithSubOrgsRequest{
		OrgData: entity.CreateOrUpdateOrgRequest{
			Name:       "新测试机构000",
			Address:    "上海市闸北区",
			AddressExt: "古北路101弄",
			Telephone:  "18800000000",
		},
		SubOrgs: []*entity.CreateOrUpdateOrgRequest{
			{
				ID:         id + 1,
				Name:       "修改机构",
				Subjects:   []string{"科目二"},
				Address:    "上海市松江区",
				AddressExt: "西子湖同",
				Telephone:  "15855552222",
			},
			{
				Name:       "修改机构222",
				Subjects:   []string{"科目三"},
				Address:    "上海市闵行区",
				AddressExt: "西子湖同222",
				Telephone:  "18898562222",
			},
		},
	}, superToken)
	if !assert.NoError(t, err) {
		return
	}

	orgRes, err := client.GetOrgById(ctx, id, superToken)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("ORG: %#v", orgRes.Org)
	for i := range orgRes.Org.SubOrgs {
		t.Logf("SUB ORG [%v]: %#v", i, orgRes.Org.SubOrgs[i])
	}
}

func TestUpdateSelfOrg(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	createOrgRes, err := client.CreateOrg(ctx, orgs[0], superToken)
	if !assert.NoError(t, err) {
		return
	}

	id := createOrgRes.ID
	t.Log(id)
	name := RandString(10)
	_, err = client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   name,
		OrgId:  id,
		RoleId: entity.RoleOutOrg,
	}, superToken)

	logRes, err := client.Login(ctx, "HelloTest", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := logRes.Data.Token


	err = client.UpdateSelfOrgById(ctx, entity.UpdateOrgWithSubOrgsRequest{
		OrgData: entity.CreateOrUpdateOrgRequest{
			Name:       "新测试机构000",
			Address:    "上海市闸北区",
			AddressExt: "古北路101弄",
			Telephone:  "18800000000",
		},
		SubOrgs: []*entity.CreateOrUpdateOrgRequest{
			{
				ID:         id + 1,
				Name:       "修改机构",
				Subjects:   []string{"科目二"},
				Address:    "上海市松江区",
				AddressExt: "西子湖同",
				Telephone:  "15855552222",
			},
			{
				Name:       "修改机构222",
				Subjects:   []string{"科目三"},
				Address:    "上海市闵行区",
				AddressExt: "西子湖同222",
				Telephone:  "18898562222",
			},
		},
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	orgRes, err := client.GetOrgById(ctx, id, superToken)
	if !assert.NoError(t, err) {
		return
	}
	t.Logf("ORG: %#v", orgRes.Org)
	for i := range orgRes.Org.SubOrgs {
		t.Logf("SUB ORG [%v]: %#v", i, orgRes.Org.SubOrgs[i])
	}
}