package da

import (
	"context"
	"xg/crypto"
	"xg/db"
	"xg/entity"
)

func AutoMigrate(){
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

func InitData(flag bool){
	if !flag {
		return
	}
	//创建超级管理员角色
	roleId, err := GetRoleModel().CreateRole(context.Background(), "超级管理员")
	if err !=nil {
		panic(err)
	}
	//主机构
	orgId, err := GetOrgModel().CreateOrg(context.Background(), Org{
		Name:      "学果网",
		Subjects:  "",
		Status:    entity.OrgStatusCertified,
	})
	if err !=nil {
		panic(err)
	}

	//创建超级管理员用户
	password := crypto.Hash("123456")
	_, err = GetUsersModel().CreateUser(context.Background(), User{
		ID:        1,
		Name:      "Admin",
		Password:  password,
		OrgId:     orgId,
		RoleId:    roleId,
	})
	if err !=nil {
		panic(err)
	}

	_, err = GetOrderSourceModel().CreateOrderSources(context.Background(), "百度平台")
	if err !=nil {
		panic(err)
	}

	initSubject()
	initAuth(roleId)
}

func initAuth(roleId int){
	enterStudentId, err := GetAuthModel().CreateAuth(context.Background(), "学员录入")
	if err != nil{
		panic(err)
	}

	err = GetRoleModel().SetRoleAuth(context.Background(), roleId, []int{enterStudentId})
	if err != nil{
		panic(err)
	}
}

func initSubject(){
	//添加课程
	designSubjectId, err := GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:     1,
		ParentId:  0,
		Name:      "设计",
	})
	if err !=nil {
		panic(err)
	}
	languageSubjectId, err := GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:     1,
		ParentId:  0,
		Name:      "外语",
	})
	if err !=nil {
		panic(err)
	}

	_, err = GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:     2,
		ParentId:  designSubjectId,
		Name:      "Photoshop",
	})
	if err !=nil {
		panic(err)
	}

	_, err = GetSubjectModel().CreateSubject(context.Background(), Subject{
		Level:     2,
		ParentId:  languageSubjectId,
		Name:      "英语",
	})
	if err !=nil {
		panic(err)
	}
}