package service

import (
	"context"
	"sync"
	"xg/da"
	"xg/entity"
)

type RoleService struct {
}

func (r *RoleService) CreateRole(ctx context.Context, name string, authList []int) (int, error) {
	id, err := da.GetRoleModel().CreateRole(ctx, name)
	if err != nil {
		return -1, err
	}
	err = da.GetRoleModel().SetRoleAuth(ctx, id, authList)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *RoleService) ListRole(ctx context.Context) ([]*entity.Role, error) {
	roles, err := da.GetRoleModel().ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*entity.Role, len(roles))
	for i := range roles {
		authInfo, err := da.GetRoleModel().ListRoleAuth(ctx, roles[i].ID)
		if err != nil {
			return nil, err
		}
		authList := make([]*entity.Auth, len(authInfo.AuthIDs))
		for i := range authInfo.AuthIDs {
			authList[i] = &entity.Auth{
				ID:   authInfo.AuthIDs[i],
				Name: authInfo.AuthNames[i],
			}
		}

		res[i] = &entity.Role{
			ID:       roles[i].ID,
			Name:     roles[i].Name,
			AuthList: authList,
		}
	}
	return res, nil
}
func (r *RoleService) SetRoleAuth(ctx context.Context, id int, ids []int) error {
	return da.GetRoleModel().SetRoleAuth(ctx, id, ids)
}

func (r *RoleService) GetRoleAuth(ctx context.Context, id int) ([]string, error) {
	res, err := da.GetRoleModel().ListRoleAuth(ctx, id)
	if err != nil {
		return nil, err
	}

	return res.AuthNames, nil
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
