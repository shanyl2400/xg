package service

import (
	"context"
	"sync"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type RoleService struct {
}

func (r *RoleService) CreateRole(ctx context.Context, name string, authList []int) (int, error) {
	id, err := da.GetRoleModel().CreateRole(ctx, name)
	if err != nil {
		log.Warning.Printf("Get role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
		return -1, err
	}
	err = da.GetRoleModel().SetRoleAuth(ctx, id, authList)
	if err != nil {
		log.Warning.Printf("Set role failed, name: %#v, authList: %#v, err: %v\n", name, authList, err)
		return -1, err
	}
	return id, nil
}

func (r *RoleService) ListRole(ctx context.Context) ([]*entity.Role, error) {
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
	err := da.GetRoleModel().SetRoleAuth(ctx, id, ids)
	if err != nil{
		log.Warning.Printf("Set role auth failed, roleId: %#v, authIds: %#v, err: %v\n", id, ids, err)
		return err
	}
	return nil
}

func (r *RoleService) GetRoleAuth(ctx context.Context, id int) ([]*entity.Auth, error) {
	res, err := da.GetRoleModel().ListRoleAuth(ctx, id)
	if err != nil {
		log.Warning.Printf("Set role auth failed, roleId: %#v, err: %v\n", id, err)
		return nil, err
	}

	return res.Auth, nil
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
