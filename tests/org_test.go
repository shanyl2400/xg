package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xg/entity"
)

func Test_008(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	createOrgRes, err := client.CreateOrg(ctx, orgs[0], superToken)
	if !assert.NoError(t, err) {
		return
	}

	orgName := "org" + RandString(3)
	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   orgName,
		OrgId:  createOrgRes.ID,
		RoleId: 7,
	}, superToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.ApprovePendingOrgById(ctx, createOrgRes.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}

	//1.登录派单员
	tryCreateDispatchUser(t)
	dispatchRes, err := client.Login(ctx, "dispatch0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	dispatchToken := dispatchRes.Data.Token

	//2.派单员给机构A派3名学生（A,B,C）
	stuRes0, err := client.CreateStudent(ctx, students[0], superToken)
	if !assert.NoError(t, err) {
		return
	}
	stuRes1, err := client.CreateStudent(ctx, students[1], superToken)
	if !assert.NoError(t, err) {
		return
	}
	stuRes2, err := client.CreateStudent(ctx, students[2], superToken)
	if !assert.NoError(t, err) {
		return
	}

	order1Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes0.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[0].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}
	order2Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes1.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[1].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}
	order3Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes2.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[2].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}


	//3.登录机构A账号
	orgARes, err := client.Login(ctx, orgName, "123456")
	if !assert.NoError(t, err) {
		return
	}
	orgAToken := orgARes.Data.Token
	//4.查看派单情况
	ordersRes, err := client.SearchOrderWithOrgId(ctx, &entity.SearchOrderCondition{
		PageSize:       10,
		Page:           0,
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, ordersRes.Data.Total, 3)

	//5.学生A报名，学生B取消，学生C报名
	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  1000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.RevokeOrder(ctx, order2Res.ID, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order3Res.ID,
		Amount:  2000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	//6.学生A缴费，学生C退费
	err = client.PayOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  500,
		Title:   "加一门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.PaybackOrder(ctx, &entity.OrderPayRequest{
		OrderID: order2Res.ID,
		Amount:  1000,
		Title:   "取消一门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	//7.登录审核员，审核订单（通过A，拒绝C）
	tryCreateCheckUser(t)
	//3.登录机构A账号
	checkRes, err := client.Login(ctx, "check0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	checkToken := checkRes.Data.Token

	recordRes, err := client.SearchOrderPayRecords(ctx, &SearchPayRecordCondition{
		PageSize:        10,
		Page:            1,
	}, checkToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, recordRes.Data.Total, 2)

	err = client.AcceptPayment(ctx, recordRes.Data.Records[0].ID, checkToken)
	if !assert.NoError(t, err) {
		return
	}
	err = client.RejectPayment(ctx, recordRes.Data.Records[1].ID, checkToken)
	if !assert.NoError(t, err) {
		return
	}

	//8.登录机构账号，查看审核情况

	summaryResponse, err := client.Summary(ctx, superToken)
	if !assert.NoError(t, err) {
		return
	}

	assert.GreaterOrEqual(t, summaryResponse.Summary.PerformanceTotal, 0)
	t.Log("Summary", summaryResponse.Summary.PerformanceTotal)
	t.Log("Done")
}

func Test_009(t *testing.T) {

	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	createOrgRes, err := client.CreateOrg(ctx, orgs[0], superToken)
	if !assert.NoError(t, err) {
		return
	}

	orgName := "org" + RandString(3)
	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   orgName,
		OrgId:  createOrgRes.ID,
		RoleId: 7,
	}, superToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.ApprovePendingOrgById(ctx, createOrgRes.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}

	//1.登录派单员
	tryCreateDispatchUser(t)
	dispatchRes, err := client.Login(ctx, "dispatch0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	dispatchToken := dispatchRes.Data.Token

	//2.给机构A派送3个订单
	stuRes0, err := client.CreateStudent(ctx, students[0], superToken)
	if !assert.NoError(t, err) {
		return
	}
	stuRes1, err := client.CreateStudent(ctx, students[1], superToken)
	if !assert.NoError(t, err) {
		return
	}
	stuRes2, err := client.CreateStudent(ctx, students[2], superToken)
	if !assert.NoError(t, err) {
		return
	}

	order1Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes0.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[0].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}
	order2Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes1.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[1].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}
	order3Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes2.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[2].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}


	//3.登录机构A
	orgARes, err := client.Login(ctx, orgName, "123456")
	if !assert.NoError(t, err) {
		return
	}
	orgAToken := orgARes.Data.Token
	//4.3个订单报名，并分别对订单A交费，订单B交费，订单C退费
	ordersRes, err := client.SearchOrderWithOrgId(ctx, &entity.SearchOrderCondition{
		PageSize:       10,
		Page:           0,
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, ordersRes.Data.Total, 3)

	//5.学生报名
	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  1000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order2Res.ID,
		Amount:  1000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order3Res.ID,
		Amount:  2000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	//学生A缴费，学生B缴费，学生C退费
	err = client.PayOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  500,
		Title:   "加1门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.PayOrder(ctx, &entity.OrderPayRequest{
		OrderID: order2Res.ID,
		Amount:  1500,
		Title:   "加2门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.PaybackOrder(ctx, &entity.OrderPayRequest{
		OrderID: order3Res.ID,
		Amount:  1000,
		Title:   "取消1门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	//5.登录审核员
	tryCreateCheckUser(t)
	checkRes, err := client.Login(ctx, "check0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	checkToken := checkRes.Data.Token

	//6.查看待审核订单
	recordRes, err := client.SearchOrderPayRecords(ctx, &SearchPayRecordCondition{
		PageSize:        10,
		Page:            1,
	}, checkToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, recordRes.Data.Total, 3)

	//7.审核订单A通过，订单B驳回，订单C通过
	err = client.AcceptPayment(ctx, recordRes.Data.Records[0].ID, checkToken)
	if !assert.NoError(t, err) {
		return
	}
	err = client.RejectPayment(ctx, recordRes.Data.Records[1].ID, checkToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.AcceptPayment(ctx, recordRes.Data.Records[2].ID, checkToken)
	if !assert.NoError(t, err) {
		return
	}

	//8.查看待审核状态
	//9.登录机构A，查看订单情况
	//10.登录派单员，查看业绩
	//11.登录录单员，查看业绩

	summaryResponse, err := client.Summary(ctx, superToken)
	if !assert.NoError(t, err) {
		return
	}

	assert.GreaterOrEqual(t, summaryResponse.Summary.PerformanceTotal, 0)
	t.Log("Summary", summaryResponse.Summary.PerformanceTotal)
	t.Log("Done")
}



func TestCheckPerformance(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	superToken := getSuperToken(t)

	summaryResponse0, err := client.Summary(ctx, superToken)
	if !assert.NoError(t, err) {
		return
	}

	assert.GreaterOrEqual(t, summaryResponse0.Summary.PerformanceTotal, 0)
	t.Log("Pre start Summary", summaryResponse0.Summary.PerformanceTotal)


	createOrgRes, err := client.CreateOrg(ctx, orgs[0], superToken)
	if !assert.NoError(t, err) {
		return
	}

	orgName := "org" + RandString(3)
	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   orgName,
		OrgId:  createOrgRes.ID,
		RoleId: 7,
	}, superToken)
	if !assert.NoError(t, err) {
		return
	}

	err = client.ApprovePendingOrgById(ctx, createOrgRes.ID, superToken)
	if !assert.NoError(t, err) {
		return
	}

	//1.登录派单员
	tryCreateDispatchUser(t)
	dispatchRes, err := client.Login(ctx, "dispatch0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	dispatchToken := dispatchRes.Data.Token

	//2.给机构A派送3个订单
	tryCreateEnterUser(t)
	enterRes, err := client.Login(ctx, "enter0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	enterToken := enterRes.Data.Token

	stuRes0, err := client.CreateStudent(ctx, students[0], enterToken)
	if !assert.NoError(t, err) {
		return
	}

	order1Res, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
		StudentID:      stuRes0.Result.ID,
		ToOrgID:        createOrgRes.ID,
		IntentSubjects: students[0].IntentSubject,
	}, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}

	//3.登录机构A
	orgARes, err := client.Login(ctx, orgName, "123456")
	if !assert.NoError(t, err) {
		return
	}
	orgAToken := orgARes.Data.Token
	//4.3个订单报名，并分别对订单A交费，订单B交费，订单C退费
	ordersRes, err := client.SearchOrderWithOrgId(ctx, &entity.SearchOrderCondition{
		PageSize:       10,
		Page:           0,
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, ordersRes.Data.Total, 1)

	//5.学生报名
	err = client.SignUpOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  1000,
		Title:   "学费",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}


	//学生A缴费，学生B缴费，学生C退费
	err = client.PayOrder(ctx, &entity.OrderPayRequest{
		OrderID: order1Res.ID,
		Amount:  500,
		Title:   "加1门课",
	}, orgAToken)
	if !assert.NoError(t, err) {
		return
	}

	//5.登录审核员
	tryCreateCheckUser(t)
	checkRes, err := client.Login(ctx, "check0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	checkToken := checkRes.Data.Token

	//6.查看待审核订单
	recordRes, err := client.SearchOrderPayRecords(ctx, &SearchPayRecordCondition{
		OrderIDList: 	 []int{order1Res.ID},
		PageSize:        10,
		Page:            1,
	}, checkToken)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, recordRes.Data.Total, 2)

	//7.审核订单A通过，订单B驳回，订单C通过
	for i := range recordRes.Data.Records{
		err = client.AcceptPayment(ctx, recordRes.Data.Records[i].ID, checkToken)
		if !assert.NoError(t, err) {
			return
		}
	}

	//8.计算总业绩
	summaryResponse1, err := client.Summary(ctx, superToken)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, summaryResponse0.Summary.PerformanceTotal + 1500, summaryResponse1.Summary.PerformanceTotal)
	t.Log("Pre start Summary", summaryResponse1.Summary.PerformanceTotal)

	//9.登录机构A，查看订单情况
	graphRes0, err := client.OrgGraph(ctx, orgAToken)
	if !assert.NoError(t, err) {
		return
	}
	for i := range graphRes0.Graph.PerformancesGraph {
		t.Logf("Org: %#v", graphRes0.Graph.PerformancesGraph[i])
		now := time.Now()
		if graphRes0.Graph.PerformancesGraph[i].Month == int(now.Month()) {
			assert.Equal(t, 1500, graphRes0.Graph.PerformancesGraph[i].Value)
		}
	}

	//10.登录派单员，查看业绩
	graphRes1, err := client.DispatchGraph(ctx, dispatchToken)
	if !assert.NoError(t, err) {
		return
	}
	for i := range graphRes1.Graph.PerformancesGraph {
		t.Logf("Dispatch: %#v", graphRes1.Graph.PerformancesGraph[i])
		now := time.Now()
		if graphRes1.Graph.PerformancesGraph[i].Month == int(now.Month()) {
			assert.Equal(t, 1500, graphRes1.Graph.PerformancesGraph[i].Value)
		}
	}

	//11.登录录单员，查看业绩
	graphRes2, err := client.EnterGraph(ctx, enterToken)
	if !assert.NoError(t, err) {
		return
	}
	for i := range graphRes2.Graph.AuthorPerformancesGraph {
		t.Logf("Dispatch: %#v", graphRes2.Graph.AuthorPerformancesGraph[i])
		now := time.Now()
		if graphRes2.Graph.AuthorPerformancesGraph[i].Month == int(now.Month()) {
			assert.Equal(t, 1500, graphRes2.Graph.AuthorPerformancesGraph[i].Value)
		}
	}
	for i := range graphRes2.Graph.PublisherPerformancesGraph {
		t.Logf("Dispatch: %#v", graphRes2.Graph.PublisherPerformancesGraph[i])
		assert.Equal(t, 0, graphRes2.Graph.PublisherPerformancesGraph[i].Value)
	}

	t.Log("Done")
}