package service

import (
	"context"
	"errors"
	"sync"
	"xg/crypto"
	"xg/da"
	"xg/entity"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidPublisherID   = errors.New("invalid publisher id")
	ErrInvalidToOrgID       = errors.New("invalid to org id")
	ErrInvalidStudentID     = errors.New("invalid student id")
	ErrNoAuthorizeToOperate = errors.New("no auth to operate")
	ErrNoAuthToOperateOrder = errors.New("no auth to operate order")
	ErrNoNeedToOperate      = errors.New("no need to operate")
)

type UserService struct {
}

func (u *UserService) Login(ctx context.Context, name, password string) (*entity.UserLoginResponse, error) {
	users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, ErrUserNotFound
	}
	user := users[0]
	if crypto.Hash(password) != user.Password {
		return nil, ErrInvalidPassword
	}

	token, err := crypto.GenerateToken(user.ID, user.OrgId, user.RoleId)
	if err != nil {
		return nil, err
	}
	userInfo, err := u.fillUserInfo(ctx, user)
	if err != nil {
		return nil, err
	}

	return &entity.UserLoginResponse{
		UserDetailsInfo: entity.UserDetailsInfo{
			UserId:   userInfo.UserId,
			RoleId:   userInfo.RoleId,
			OrgId:    userInfo.OrgId,
			RoleName: userInfo.RoleName,
			OrgName:  userInfo.OrgName,
			Auths:    userInfo.Auths,
		},
		Token: token,
	}, nil
}

func (u *UserService) fillUserInfo(ctx context.Context, user *da.User) (*entity.UserDetailsInfo, error) {
	//获取角色和权限
	roleInfo, err := da.GetRoleModel().GetRoleById(ctx, user.RoleId)
	if err != nil {
		return nil, err
	}
	auth, err := da.GetRoleModel().ListRoleAuth(ctx, user.RoleId)
	if err != nil {
		return nil, err
	}

	orgInfo, err := da.GetOrgModel().GetOrgById(ctx, user.OrgId)
	if err != nil {
		return nil, err
	}
	return &entity.UserDetailsInfo{
		UserId:   user.ID,
		RoleId:   user.RoleId,
		OrgId:    user.OrgId,
		RoleName: roleInfo.Name,
		OrgName:  orgInfo.Name,
		Auths:    auth.AuthNames,
	}, nil
}

func (u *UserService) UpdatePassword(ctx context.Context, newPassword string, operator *entity.JWTUser) error {

	user, err := da.GetUsersModel().GetUserById(ctx, operator.UserId)
	if err != nil {
		return err
	}
	user.Password = crypto.Hash(newPassword)
	err = da.GetUsersModel().UpdateUser(ctx, *user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ListUserAuthority(ctx context.Context, operator *entity.JWTUser) ([]*entity.Auth, error) {
	authList, err := da.GetRoleModel().ListRoleAuth(ctx, operator.RoleId)
	if err != nil {
		return nil, err
	}
	authObjList, err := da.GetAuthModel().ListAuthByIDs(ctx, authList.AuthIDs)
	if err != nil {
		return nil, err
	}
	res := make([]*entity.Auth, len(authObjList))
	for i := range authObjList {
		res[i] = (*entity.Auth)(authObjList[i])
	}

	return res, nil
}

func (u *UserService) checkUserEntity(ctx context.Context, req *entity.CreateUserRequest) error {
	_, err := da.GetOrgModel().GetOrgById(ctx, req.OrgId)
	if err != nil {
		return err
	}
	_, err = da.GetRoleModel().GetRoleById(ctx, req.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (int, error) {
	//check orgId & roleId
	err := u.checkUserEntity(ctx, req)
	if err != nil {
		return -1, err
	}
	return da.GetUsersModel().CreateUser(ctx, da.User{
		Name:     req.Name,
		Password: crypto.Hash("123456"),
		OrgId:    req.OrgId,
		RoleId:   req.RoleId,
	})
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
