package da

import (
	"context"
	"xg/crypto"
	"xg/db"
	"xg/entity"
)

func AutoMigrate() {
	db.Get().AutoMigrate(Auth{})
	db.Get().AutoMigrate(Order{})
	db.Get().AutoMigrate(OrderPayRecord{})
	db.Get().AutoMigrate(OrderRemarkRecord{})
	db.Get().AutoMigrate(OrderSource{})
	db.Get().AutoMigrate(Org{})
	db.Get().AutoMigrate(Role{})
	db.Get().AutoMigrate(RoleAuth{})
	db.Get().AutoMigrate(Student{})
	db.Get().AutoMigrate(StudentNote{})
	db.Get().AutoMigrate(Subject{})
	db.Get().AutoMigrate(User{})
}

func InitData(flag bool) {
	if !flag {
		return
	}
	o, _ := GetOrgModel().GetOrgById(context.Background(), 1)
	if o != nil {
		return
	}

	//主机构
	orgId, err := GetOrgModel().CreateOrg(context.Background(), Org{
		Name:     "学果网",
		Subjects: "",
		Status:   entity.OrgStatusCertified,
	})
	if err != nil {
		panic(err)
	}

	initAuth()
	initRole()

	password := crypto.Hash("123456")
	_, err = GetUsersModel().CreateUser(context.Background(), User{
		ID:       1,
		Name:     "Admin",
		Password: password,
		OrgId:    orgId,
		RoleId:   1,
	})
	if err != nil {
		panic(err)
	}

	_, err = GetOrderSourceModel().CreateOrderSources(context.Background(), "百度平台")
	if err != nil {
		panic(err)
	}

	initSubject()
}

func initAuth() {
	err := GetAuthModel().CreateAuthWithID(context.Background(), 1, "录单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 2, "自派单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 3, "全名单派单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 4, "审核订单权限")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 5, "查看所有订单")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 6, "机构订单权限")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 7, "订单来源管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 8, "课程分类管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 9, "机构管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 10, "机构审核")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 11, "用户管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), 12, "角色管理")
	if err != nil {
		panic(err)
	}
}

func initRole() {
	//创建超级管理员角色
	adminId, err := GetRoleModel().CreateRoleWithID(context.Background(), 1, "超级管理员")
	if err != nil {
		panic(err)
	}

	err = GetRoleModel().SetRoleAuth(context.Background(), adminId, []int{1, 3, 4, 5, 7, 8, 9, 10, 11, 12})
	if err != nil {
		panic(err)
	}

	enterId, err := GetRoleModel().CreateRoleWithID(context.Background(), 2, "录单员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), enterId, []int{1, 2})
	if err != nil {
		panic(err)
	}

	dispatchId, err := GetRoleModel().CreateRoleWithID(context.Background(), 3, "派单员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), dispatchId, []int{3, 5})
	if err != nil {
		panic(err)
	}

	userId, err := GetRoleModel().CreateRoleWithID(context.Background(), 4, "人员管理员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), userId, []int{11})
	if err != nil {
		panic(err)
	}

	orgId, err := GetRoleModel().CreateRoleWithID(context.Background(), 5, "机构管理员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), orgId, []int{9})
	if err != nil {
		panic(err)
	}

	checkId, err := GetRoleModel().CreateRoleWithID(context.Background(), 6, "审核员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), checkId, []int{4, 5})
	if err != nil {
		panic(err)
	}
}

func initSubject() {
	//添加课程
	designSubjectId, err := GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:    1,
		ParentId: 0,
		Name:     "设计",
	})
	if err != nil {
		panic(err)
	}
	languageSubjectId, err := GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:    1,
		ParentId: 0,
		Name:     "外语",
	})
	if err != nil {
		panic(err)
	}

	_, err = GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:    2,
		ParentId: designSubjectId,
		Name:     "Photoshop",
	})
	if err != nil {
		panic(err)
	}

	_, err = GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:    2,
		ParentId: languageSubjectId,
		Name:     "英语",
	})
	if err != nil {
		panic(err)
	}
}
