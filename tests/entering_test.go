package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
	"xg/conf"
	"xg/da"
	"xg/entity"
)

func Test_001(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "abcabc")
	if !assert.Error(t, err) {
		return
	}
	//2.使用正确密码登录录入员
	res, err = client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token
	t.Log("Token:", token)
	//3.查看录入员权限
	auths, err := client.ListAuths(ctx, token)
	if !assert.NoError(t, err) {
		return
	}

	for i := range auths.AuthList {
		t.Logf("%v[%v]:%v", "auth", i, auths.AuthList[i])
	}

	//4.修改密码为错误密码1
	err = client.UpdatePassword(ctx, "123123", token)
	if !assert.NoError(t, err) {
		return
	}

	//5.使用正确密码登录
	res, err = client.Login(ctx, "Admin", "123456")
	if !assert.Error(t, err) {
		return
	}

	//6.使用错误密码1登录
	res, err = client.Login(ctx, "Admin", "123123")
	if !assert.NoError(t, err) {
		return
	}

	err = client.UpdatePassword(ctx, "123456", token)
	if !assert.NoError(t, err) {
		return
	}
	t.Log("Done")
}

func tryCreateEnterUser(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "enter0",
		OrgId:  1,
		RoleId: 2,
	}, superToken)
}

func tryCreateCheckUser(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "check0",
		OrgId:  1,
		RoleId: 6,
	}, superToken)
}


func tryCreateDispatchUser(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "dispatch0",
		OrgId:  1,
		RoleId: 3,
	}, superToken)
}


func tryCreateUserManagerUser(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "user0",
		OrgId:  1,
		RoleId: 4,
	}, superToken)
}


func tryCreateOrgManagerUser(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "morg0",
		OrgId:  1,
		RoleId: 5,
	}, superToken)
}

func tryCreateEnterUser2(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "enter0",
		OrgId:  1,
		RoleId: 2,
	}, superToken)

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "enter1",
		OrgId:  1,
		RoleId: 2,
	}, superToken)
}


func tryCreateUserAndOrg(t *testing.T){
	client := new(APIClient)
	ctx := context.Background()
	//1.使用错误密码1登录录入员
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res.Data.Token

	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "enter0",
		OrgId:  1,
		RoleId: 2,
	}, superToken)

	for i := range orgs {
		client.CreateOrg(ctx, orgs[i], superToken)
	}
}

func Test_002(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateEnterUser(t)
	//1.登录录入员
	t.Log("Log in enter user")
	res, err := client.Login(ctx, "enter0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token
	//2.录入3个学生
	t.Log("Create student1")
	students[0].IntentSubject = []string{subjects[0], subjects[1]}
	students[0].Address = addresses[0]
	res1, err := client.CreateStudent(ctx, students[0], token)
	if !assert.NoError(t, err) {
		return
	}

	t.Log("Create student2")
	students[1].IntentSubject = []string{subjects[1], subjects[2]}
	students[1].Address = addresses[0]
	res2, err := client.CreateStudent(ctx, students[1], token)
	if !assert.NoError(t, err) {
		return
	}

	t.Log("Create student3")
	students[2].IntentSubject = []string{subjects[0], subjects[2]}
	students[2].Address = addresses[1]
	res3, err := client.CreateStudent(ctx, students[2], token)
	if !assert.NoError(t, err) {
		return
	}

	//3.查询学生信息
	//检查学生1
	t.Log("Checking students")
	studentRes, err := client.GetStudentById(ctx, res1.Result.ID, token)
	if !assert.NoError(t, err) {
		return
	}

	stu1 := studentRes.Student
	assert.Equal(t, addresses[0], stu1.Address)
	assert.Equal(t, entity.StudentCreated, stu1.Status)
	assert.Equal(t, students[0].Name, stu1.Name)
	assert.Equal(t, students[0].Telephone, stu1.Telephone)


	//检查学生2
	studentRes2, err := client.GetStudentById(ctx, res2.Result.ID, token)
	if !assert.NoError(t, err) {
		return
	}
	stu2 := studentRes2.Student
	assert.Equal(t, addresses[0], stu2.Address)
	assert.Equal(t, entity.StudentCreated, stu2.Status)
	assert.Equal(t, students[1].Name, stu2.Name)
	assert.Equal(t, students[1].Telephone, stu2.Telephone)

	//检查学生3
	studentRes3, err := client.GetStudentById(ctx, res3.Result.ID, token)
	if !assert.NoError(t, err) {
		return
	}
	stu3 := studentRes3.Student
	assert.Equal(t, addresses[1], stu3.Address)
	assert.Equal(t, entity.StudentCreated, stu3.Status)
	assert.Equal(t, students[2].Name, stu3.Name)
	assert.Equal(t, students[2].Telephone, stu3.Telephone)

	//4.根据姓名查询学生
	t.Log("checking search students by name")
	result0, err := client.SearchPrivateStudents(ctx, &entity.SearchStudentRequest{
		Name:       students[0].Name,
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	//地址要包含1和2
	if !assert.NoError(t, containsStudents([]int{res1.Result.ID}, result0.Result.Students)) {
		return
	}

	//5.根据地址和兴趣查询学生
	t.Log("checking search students by address & subjects")
	result1, err := client.SearchPrivateStudents(ctx, &entity.SearchStudentRequest{
		Address:       addresses[0],
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	//地址要包含1和2
	if !assert.NoError(t, containsStudents([]int{res1.Result.ID, res2.Result.ID}, result1.Result.Students)) {
		return
	}

	result2, err := client.SearchPrivateStudents(ctx, &entity.SearchStudentRequest{
		IntentSubject:       subjects[0],
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	//课程包括1和3
	if !assert.NoError(t, containsStudents([]int{res1.Result.ID, res3.Result.ID}, result2.Result.Students)) {
		return
	}
	t.Log("Done")
}

func Test_003(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateEnterUser(t)
	//1.登录录入员
	t.Log("Log in enter user")
	res, err := client.Login(ctx, "enter0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token

	//2.录入3个学生
	sids := make([]int, 3)
	for i := 0; i < 3; i ++ {
		res, err := client.CreateStudent(ctx, students[i], token)
		if !assert.NoError(t, err) {
			return
		}
		t.Logf("Student: %#v", students[i])
		sids[i]= res.Result.ID
	}

	//3.为3个学生派单（每个学生分配4个机构）
	subOrgRes, err := client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		IsSubOrg:  true,
	}, token)
	if !assert.NoError(t, err) {
		return
	}
	assert.GreaterOrEqual(t, subOrgRes.Data.Total, 1)
	assert.GreaterOrEqual(t, len(subOrgRes.Data.Orgs), 1)
	for i := range sids{
		_, err := client.CreateOrder(ctx, &entity.CreateOrderRequest{
			StudentID:      sids[i],
			ToOrgID:        subOrgRes.Data.Orgs[i].ID,
			IntentSubjects: students[i].IntentSubject,
		}, token)
		if !assert.NoError(t, err) {
			return
		}
	}

	//4.查询学生派单情况
	for i := range sids {
		studentRes, err := client.GetStudentById(ctx, sids[i], token)
		assert.NoError(t, err)
		if assert.GreaterOrEqual(t, len(studentRes.Student.Orders), 1) {
			order := studentRes.Student.Orders[0]
			assert.Equal(t, subOrgRes.Data.Orgs[i].Name, order.OrgName)
			assert.Equal(t, subOrgRes.Data.Orgs[i].ID, order.ToOrgID)
			assert.Equal(t, len(students[i].IntentSubject), len(order.IntentSubject))
		}
	}

	//5.根据机构查询学员派单
	res0, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	superToken := res0.Data.Token

	if !assert.NotEqual(t, 0, subOrgRes.Data.Orgs[0].ParentID){
		return
	}
	client.CreateUser(ctx, &entity.CreateUserRequest{
		Name:   "org0",
		OrgId:  subOrgRes.Data.Orgs[0].ParentID,
		RoleId: 7,
	}, superToken)

	//登录机构账号
	logRes, err := client.Login(ctx,  "org0", "123456")
	if !assert.NoError(t, err) {
		return
	}

	orgToken := logRes.Data.Token
	//查询机构订单
	ordersRes, err := client.SearchOrderWithOrgId(ctx, &entity.SearchOrderCondition{}, orgToken)
	if !assert.NoError(t, err) {
		return
	}
	flag := false
	for _, order := range ordersRes.Data.Orders{
		if sids[0] == order.StudentID{
			flag = true
			t.Logf("Find student, orderId: %v, studentId: %v", order.ID, order.StudentID)
			break
		}
	}
	assert.Equal(t, true, flag)

	t.Log("Done")
}

func Test_004(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateEnterUser2(t)
	//1.登录录入员
	t.Log("Log in enter user")
	res, err := client.Login(ctx, "enter0", "123456")
	if !assert.NoError(t, err) {
		return
	}

	token1 := res.Data.Token
	//2.录入学生1，学生2
	cres0, err := client.CreateStudent(ctx, students[0], token1)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, cres0.Result.Status, entity.StudentCreated)

	cres1, err := client.CreateStudent(ctx, students[1], token1)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, cres1.Result.Status, entity.StudentCreated)

	//3.登录录入员B
	res1, err := client.Login(ctx, "enter0", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token2 := res1.Data.Token
	//4.录入学生2，学生3
	cres, err := client.CreateStudent(ctx, students[1], token2)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, cres.Result.Status, entity.StudentConflictFailed)

	cres, err = client.CreateStudent(ctx, students[2], token2)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, cres.Result.Status, entity.StudentCreated)


	cres, err = client.CreateStudent(ctx, students[1], token1)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, cres.Result.Status, entity.StudentConflictFailed)

	//5.修改学生1添加时间为30天前
	conf.Set(&conf.Config{DBConnectionString: "root:Badanamu123456@tcp(localhost:3306)/xg?parseTime=true&charset=utf8mb4"})
	createdAt := time.Now().Add(-time.Hour * 24 * 64)
	err = da.GetStudentModel().UpdateStudent(ctx, cres1.Result.ID, da.Student{CreatedAt: &createdAt})
	if !assert.NoError(t, err) {
		return
	}
	//6.登录录入员B，录入学生1
	cres, err = client.CreateStudent(ctx, students[1], token2)
	if !assert.NoError(t, err) {
		return
	}
	t.Log("StudentID:", cres1.Result.ID)
	assert.Equal(t, cres.Result.Status, entity.StudentConflictSuccess)
	t.Log("Done")
}

func Test_SearchOrders(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateEnterUser(t)

	//1.登录录入员
	t.Log("Log in super user")
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token
	ordersResp, err := client.SearchOrderWithAuthor(ctx, &entity.SearchOrderCondition{PageSize:10}, token)
	if !assert.NoError(t, err) {
		return
	}
	fmt.Printf("%#v", ordersResp.Data)
}

func Test_SearchSubOrgs(t *testing.T) {
	client := new(APIClient)
	ctx := context.Background()
	tryCreateEnterUser(t)

	//1.登录录入员
	t.Log("Log in super user")
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return
	}
	token := res.Data.Token

	t.Logf("Create org %#v", orgs[0])
	orgs[0].SubOrgs[0].Subjects = []string{subjects[0], subjects[1]}
	orgs[0].SubOrgs[0].Address = addresses[8]
	orgRes, err := client.CreateOrg(ctx, orgs[0], token)
	if !assert.NoError(t, err) {
		return
	}
	orgInfoRes, err := client.GetOrgById(ctx, orgRes.ID, token)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.GreaterOrEqual(t, len(orgInfoRes.Org.SubOrgs), 1) {
		return
	}
	t.Logf("Checking org")
	err = client.ApprovePendingOrgById(ctx, orgRes.ID, token)
	if !assert.NoError(t, err) {
		return
	}

	//检查包含Address和Subjects的情况
	t.Logf("Checking search with address & subjects")
	orgsRes, err := client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		Subjects:  []string{subjects[0], subjects[1]},
		Address:   addresses[8],
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}

	//检查只包含Address的情况
	t.Logf("Checking search with address")
	orgsRes, err = client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		Address:   addresses[8],
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}

	//检查只包含部分subjects的情况
	t.Logf("Checking search with parts of subjects")
	orgsRes, err = client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		Subjects:  []string{subjects[1]},
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}

	//检查包含错误address的情况
	t.Logf("Checking search with wrong address")
	orgsRes, err = client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		Subjects:  []string{subjects[1]},
		Address: addresses[0],
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.Error(t, err) {
		return
	}

	//检查包含部分subjects的情况
	t.Logf("Checking search with parts subjects")
	orgsRes, err = client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
		Subjects:  []string{subjects[1], subjects[6]},
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}

	//检查全查询
	t.Logf("Checking search with all")
	orgsRes, err = client.SearchSubOrgs(ctx, da.SearchOrgsCondition{
	}, token)
	if !assert.NoError(t, err) {
		return
	}

	err = containsOrgs([]int{orgInfoRes.Org.SubOrgs[0].ID},orgsRes.Data.Orgs)
	if !assert.NoError(t, err) {
		return
	}

}

func containsStudents(ids []int, students []*entity.StudentInfo) error {
	flags := make([]bool, len(ids))

	for i := range students {
		for j := range ids {
			if students[i].ID == ids[j] {
				flags[j] = true
			}
		}
	}

	for i := range flags {
		if !flags[i] {
			return errors.New("Can't find student, id:" + strconv.Itoa(ids[i]))
		}
	}
	return nil
}

func getSuperToken(t *testing.T) string{
	client := new(APIClient)
	ctx := context.Background()
	res, err := client.Login(ctx, "Admin", "123456")
	if !assert.NoError(t, err) {
		return ""
	}
	return res.Data.Token
}

func containsOrgs(ids []int, orgs []*entity.SubOrgWithDistance) error {
	flags := make([]bool, len(ids))

	for i := range orgs {
		for j := range ids {
			if orgs[i].ID == ids[j] {
				fmt.Println("Find org:", orgs[i].ID)
				flags[j] = true
			}
		}
	}

	for i := range flags {
		if !flags[i] {
			return errors.New("Can't find org, id:" + strconv.Itoa(ids[i]))
		}
	}
	return nil
}


func containsOrgs2(ids []int, orgs []*entity.Org) error {
	flags := make([]bool, len(ids))

	for i := range orgs {
		for j := range ids {
			if orgs[i].ID == ids[j] {
				fmt.Println("Find org:", orgs[i].ID)
				flags[j] = true
			}
		}
	}

	for i := range flags {
		if !flags[i] {
			return errors.New("Can't find org, id:" + strconv.Itoa(ids[i]))
		}
	}
	return nil
}

func containsOrders(ids []int, orgs []*entity.OrderInfoDetails) error {
	flags := make([]bool, len(ids))

	for i := range orgs {
		for j := range ids {
			if orgs[i].ID == ids[j] {
				fmt.Println("Find org:", orgs[i].ID)
				flags[j] = true
			}
		}
	}

	for i := range flags {
		if !flags[i] {
			return errors.New("Can't find org, id:" + strconv.Itoa(ids[i]))
		}
	}
	return nil
}
