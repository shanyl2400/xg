package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"xg/entity"
)

func Test_006(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateUserManagerUser(t)

	//1.登录派单员
	t.Log("Log in super user")
	res, err := client.Login(ctx, "user0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token

	cRes, err := client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "TestUser",
		OrgId:  1,
		RoleId: 4,
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	userResp, err := client.ListUsers(ctx, token)
	if !assert.NoError(t, err) {
		return
	}

	flag := false
	for i := range userResp.Users {
		if cRes.ID == userResp.Users[i].UserId{
			flag = true
		}
	}

	assert.Equal(t, true, flag)
	t.Log("Done")
}

func Test_007(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateOrgManagerUser(t)

	//1.登录机构管理员
	t.Log("Log in super user")
	res, err := client.Login(ctx, "morg0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token

	auths, err := client.ListUserAuthority(ctx, token)
	if !assert.NoError(t, err) {
		return
	}
	for i := range auths.AuthList {
		t.Logf("%#v\n", auths.AuthList[i])
	}

	//2.创建机构A，B
	org0Resp, err := client.CreateOrg(ctx, orgs[1], token)
	if !assert.NoError(t, err) {
		return
	}
	org1Resp, err := client.CreateOrg(ctx, orgs[2], token)
	if !assert.NoError(t, err) {
		return
	}
	//3.查看机构列表
	orgs, err := client.ListOrgs(ctx, token)
	if !assert.NoError(t, err) {
		return
	}
	err = containsOrgs2([]int{org0Resp.ID, org1Resp.ID}, orgs.Data.Orgs)
	if !assert.Error(t, err) {
		return
	}
	//4.登录超级管理员，审核机构A通过，机构B驳回
	superToken := getSuperToken(t)
	err = client.ApprovePendingOrgById(ctx, org0Resp.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.RejectPendingOrgById(ctx, org1Resp.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}

	//5.登录机构管理员，查看机构订单
	//3.查看机构列表
	orgs, err = client.ListOrgs(ctx, token)
	if !assert.NoError(t, err) {
		return
	}
	t.Log("orgs:", []int{org0Resp.ID, org1Resp.ID})
	err = containsOrgs2([]int{org0Resp.ID}, orgs.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}
	err = containsOrgs2([]int{org1Resp.ID}, orgs.Data.Orgs)
	if !assert.Error(t, err) {
		return
	}
}


func Test_010(t *testing.T){
	//1.登录超级管理员
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	//2.创建订单来源
	resp, err := client.CreateOrderSources(ctx, "网易自选平台", superToken)
	if !assert.NoError(t, err) {
		return
	}

	//3.查看订单来源列表
	//4.录入学员，查看订单来源列表
	oss, err := client.ListOrderSources(ctx, superToken)
	if !assert.NoError(t, err) {
		return
	}

	flag := false
	for i := range oss.Sources {
		if resp.ID == oss.Sources[i].ID {
			flag = true
			break
		}
	}
	assert.Equal(t, true, flag)
}

func Test_011(t *testing.T){
	//1.登录超级管理员
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)
	//2.创建课程父分类

	cidRes, err := client.CreateSubject(ctx, entity.CreateSubjectRequest{
		ParentId: 0,
		Name:     "服装",
	}, superToken)
	if !assert.NoError(t, err) {
		return
	}

	//3.查看课程列表
	subjectsRes, err := client.ListSubjects(ctx, 0, superToken)
	if !assert.NoError(t, err) {
		return
	}
	flag := false
	for i := range subjectsRes.Subjects {
		if subjectsRes.Subjects[i].ID == cidRes.ID {
			flag = true
		}
	}
	assert.Equal(t, true, flag)

	//4.创建课程子分类
	cidRes2, err := client.CreateSubject(ctx, entity.CreateSubjectRequest{
		ParentId: cidRes.ID,
		Name:     "裁剪",
	}, superToken)

	//5.查看课程列表
	//6.录入学员，查看课程列表
	subjectsRes, err = client.ListSubjects(ctx, cidRes.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}
	flag = false
	for i := range subjectsRes.Subjects {
		if subjectsRes.Subjects[i].ID == cidRes2.ID {
			flag = true
		}
	}
	assert.Equal(t, true, flag)
}