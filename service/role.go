package service

import (
	"context"
	"errors"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

const (
	RoleModeRootRole = 1
	RoleModeOrgRole  = 2
)

var (
	ErrInvalidAuth     = errors.New("invalid auth")
	ErrInvalidAuthMode = errors.New("invalid auth mode")
)

type IRoleService interface {
	CreateRole(ctx context.Context, name string, authList []int) (int, error)
	ListRole(ctx context.Context) ([]*entity.Role, error)
	SetRoleAuth(ctx context.Context, id int, ids []int) error
	GetRoleAuth(ctx context.Context, id int) ([]*entity.Auth, error)
}

type RoleService struct {
}

func (r *RoleService) CreateRole(ctx context.Context, name string, authList []int) (int, error) {
	log.Info.Printf("CreateRole, name: %#v, authList: %#v\n", name, authList)
	id, err := db.GetTransResult(ctx, func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		id, err := da.GetRoleModel().CreateRole(ctx, tx, name)
		if err != nil {
			log.Warning.Printf("Get role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
			return -1, err
		}
		err = da.GetRoleModel().SetRoleAuth(ctx, tx, id, authList)
		if err != nil {
			log.Warning.Printf("Set role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
			return -1, err
		}
		return id, nil
	})
	if err != nil {
		return -1, err
	}
	return id.(int), nil
}

func (r *RoleService) CreateRoleForOrgs(ctx context.Context, name string, authList []int, roleMode int) (int, error) {
	log.Info.Printf("CreateRole, name: %#v, authList: %#v\n", name, authList)

	err := r.checkAuth(ctx, authList, roleMode)
	if err != nil {
		return -1, err
	}
	id, err := db.GetTransResult(ctx, func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		id, err := da.GetRoleModel().CreateRole(ctx, tx, name)
		if err != nil {
			log.Warning.Printf("Get role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
			return -1, err
		}
		err = da.GetRoleModel().SetRoleAuth(ctx, tx, id, authList)
		if err != nil {
			log.Warning.Printf("Set role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
			return -1, err
		}

		//Update orgs
		err = da.GetOrgModel().UpdateOrgsRole(ctx, tx, roleMode == RoleModeOrgRole, id)
		if err != nil {
			log.Warning.Printf("UpdateOrgsRole failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
			return -1, err
		}

		return id, nil
	})
	if err != nil {
		return -1, err
	}
	return id.(int), nil
}

func (r *RoleService) ListRole(ctx context.Context) ([]*entity.Role, error) {
	log.Info.Printf("ListRole\n")
	roles, err := da.GetRoleModel().ListRoles(ctx)
	if err != nil {
		log.Warning.Printf("List role failed, err: %v\n", err)
		return nil, err
	}
	res := make([]*entity.Role, len(roles))
	for i := range roles {
		authInfo, err := da.GetRoleModel().ListRoleAuth(ctx, roles[i].ID)
		if err != nil {
			log.Warning.Printf("List role auth failed, role: %#v, err: %v\n", roles[i], err)
			return nil, err
		}

		res[i] = &entity.Role{
			ID:       roles[i].ID,
			Name:     roles[i].Name,
			AuthList: authInfo.Auth,
		}
	}
	return res, nil
}
func (r *RoleService) SetRoleAuth(ctx context.Context, id int, ids []int) error {
	log.Info.Printf("SetRoleAuth, roleId: %#v, authList: %#v\n", id, ids)
	err := db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err := da.GetRoleModel().SetRoleAuth(ctx, tx, id, ids)
		if err != nil {
			log.Warning.Printf("Set role auth failed, roleId: %#v, authIds: %#v, err: %v\n", id, ids, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RoleService) GetRoleAuth(ctx context.Context, id int) ([]*entity.Auth, error) {
	log.Info.Printf("GetRoleAuth, roleId: %#v\n", id)
	res, err := da.GetRoleModel().ListRoleAuth(ctx, id)
	if err != nil {
		log.Warning.Printf("Set role auth failed, roleId: %#v, err: %v\n", id, err)
		return nil, err
	}

	return res.Auth, nil
}

func (r *RoleService) checkAuth(ctx context.Context, authList []int, roleMode int) error {
	auths, err := GetAuthService().ListAuths(ctx)
	if err != nil {
		log.Warning.Printf("Get Auths failed, authList: %#v, err: %v\n", authList, err)
		return err
	}
	authMap := make(map[int]*entity.Auth)
	for i := range auths {
		authMap[auths[i].ID] = auths[i]
	}

	for i := range authList {
		auth, ok := authMap[authList[i]]
		if !ok {
			log.Warning.Printf("invalid auth, aluth: %#v authList: %#v\n", authList[i], authList)
			return ErrInvalidAuth
		}
		if (auth.Mode == entity.AuthModeOrgAuth && roleMode != RoleModeOrgRole) ||
			(auth.Mode == entity.AuthModeRootAuth && roleMode != RoleModeRootRole) {
			log.Warning.Printf("invalid auth mode, aluth : %#v authList: %#v\n", auth, authList)
			return ErrInvalidAuthMode
		}

	}

	return nil
}

var (
	_roleService     *RoleService
	_roleServiceOnce sync.Once
)

func GetRoleService() *RoleService {
	_roleServiceOnce.Do(func() {
		if _roleService == nil {
			_roleService = new(RoleService)
		}
	})
	return _roleService
}
