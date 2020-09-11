package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"xg/da"
	"xg/entity"
)

func Test_005(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateDispatchUser(t)

	//1.登录派单员
	t.Log("Log in super user")
	res, err := client.Login(ctx, "dispatch0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token
	//2.根据兴趣和地区查询学员
	//3.进入派单页面，查看机构兴趣和地区是否匹配
	//见Test_SearchSubOrgs

	//4.派单
	//查询学员
	studentsRes, err := client.SearchStudents(ctx, &entity.SearchStudentRequest{
		PageSize:      4,
		Page:          1,
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, studentsRes.Result.Total, 1)

	size := 0
	for i, stu := range studentsRes.Result.Students {
		subOrgs, err := client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
			Subjects:  stu.IntentSubject,
			Address:   stu.Address,
			IsSubOrg:  true,
		}, token)
		if !assert.NoError(t, err) {
			return
		}
		if subOrgs.Data.Total > 0 {
			t.Log("Found sub orgs in ", i, ", subOrgId:", subOrgs.Data.Orgs[0].ID, "studentId:", stu.ID)
			_, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
				StudentID:      stu.ID,
				ToOrgID:        subOrgs.Data.Orgs[0].ID,
				IntentSubjects: subOrgs.Data.Orgs[0].Subjects,
			}, token)
			if !assert.NoError(t, err) {
				return
			}
			size ++
		}
	}
	//5.查看派单列表
	ordersResp, err := client.SearchOrderWithAuthor(ctx, &entity.SearchOrderCondition{PageSize:10}, token)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, ordersResp.Data.Total, size)
	t.Log("Total:", ordersResp.Data.Total)
	t.Log("Page:", len(ordersResp.Data.Orders))

	for _, order := range ordersResp.Data.Orders {
		//创建机构账号
		superToken := getSuperToken(t)
		_, err = client.CreateUser(ctx, &entity.CreateUserRequest{
			Name:   "org_" + strconv.Itoa(order.ToOrgID),
			OrgId:  order.ToOrgID,
			RoleId: 7,
		}, superToken)
		if err != nil{
			t.Log("create user failed, ", err)
		}

		//登录机构账号
		orgLoginResp, err := client.Login(ctx,  "org_" + strconv.Itoa(order.ToOrgID),"123456")
		if !assert.NoError(t, err) {
			return
		}
		orgToken := orgLoginResp.Data.Token

		//机构查询
		ordersResp0, err := client.SearchOrderWithOrgId(ctx, &entity.SearchOrderCondition{PageSize: 10000}, orgToken)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.NoError(t, containsOrders([]int{order.ID}, ordersResp0.Data.Orders)) {
			return
		}
	}

	t.Log("Done")
}
