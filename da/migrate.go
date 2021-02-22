package da

import (
	"context"
	"fmt"
	"xg/crypto"
	"xg/db"
	"xg/entity"
	"xg/log"
)

func AutoMigrate() {
	db.Get().AutoMigrate(OrderNotifies{})
	db.Get().AutoMigrate(Auth{})
	db.Get().AutoMigrate(StudentConflict{})
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
	db.Get().AutoMigrate(StatisticsRecord{})

	db.Get().AutoMigrate(OrderStatisticsRecord{})
	db.Get().Model(&OrderStatisticsRecord{}).AddIndex("idx_new_statistics_date", "year", "month", "day")
}

func InitData(flag bool) {
	if !flag {
		return
	}
	o, _ := GetOrgModel().GetOrgById(context.Background(), db.Get(), 1)
	if o != nil {
		return
	}

	//主机构
	orgId, err := GetOrgModel().CreateOrg(context.Background(), db.Get(), Org{
		Name:     "学果网",
		Subjects: "",
		Status:   entity.OrgStatusCertified,
		SupportRoleID: fmt.Sprintf("%v,%v,%v,%v,%v,%v", entity.RoleSuperAdmin,
			entity.RoleEnterWorker, entity.RoleDispatchWorker, entity.RoleUserManager,
			entity.RoleOrgManager, entity.RoleChecker),
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
	err := GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthEnterStudent, "录单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthDispatchSelfOrder, "自派单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthDispatchOrder, "全名单派单权")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthCheckOrder, "审核订单权限")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthListAllOrder, "查看所有订单")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthListOrgOrder, "机构订单权限")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageOrderSource, "订单来源管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageSubject, "课程分类管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageOrg, "机构管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthCheckOrg, "机构审核")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageUser, "用户管理")
	if err != nil {
		panic(err)
	}
	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageRole, "角色管理")
	if err != nil {
		panic(err)
	}

	err = GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageSelfOrg, "管理本机构信息")
	if err != nil {
		panic(err)
	}
}

func Migrate0(flag bool) {
	if !flag {
		return
	}
	err := GetAuthModel().CreateAuthWithID(context.Background(), entity.AuthManageSelfOrg, "管理本机构信息")
	if err != nil {
		log.Error.Println(err)
		return
	}
	outOrgId2, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleSeniorOutOrg, "高级机构账号")
	if err != nil {
		log.Error.Println(err)
		return
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), outOrgId2, []int{6, 13})
	if err != nil {
		log.Error.Println(err)
		return
	}
}

func initRole() {
	//创建超级管理员角色
	adminId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleSuperAdmin, "超级管理员")
	if err != nil {
		panic(err)
	}

	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), adminId, []int{1, 3, 4, 5, 7, 8, 9, 10, 11, 12})
	if err != nil {
		panic(err)
	}

	enterId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleEnterWorker, "录单员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), enterId, []int{1, 2})
	if err != nil {
		panic(err)
	}

	dispatchId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleDispatchWorker, "派单员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), dispatchId, []int{3})
	if err != nil {
		panic(err)
	}

	userId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleUserManager, "人员管理员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), userId, []int{11})
	if err != nil {
		panic(err)
	}

	orgId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleOrgManager, "机构管理员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), orgId, []int{9})
	if err != nil {
		panic(err)
	}

	checkId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleChecker, "审核员")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), checkId, []int{4, 5})
	if err != nil {
		panic(err)
	}

	outOrgId, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleOutOrg, "机构账号")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), outOrgId, []int{6})
	if err != nil {
		panic(err)
	}

	outOrgId2, err := GetRoleModel().CreateRoleWithID(context.Background(), entity.RoleSeniorOutOrg, "高级机构账号")
	if err != nil {
		panic(err)
	}
	err = GetRoleModel().SetRoleAuth(context.Background(), db.Get(), outOrgId2, []int{6, 13})
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
