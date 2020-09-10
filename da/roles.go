package da

import (
	"context"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
	"xg/db"
	"xg/entity"
)

type IRolesModel interface {
	CreateRole(ctx context.Context, tx *gorm.DB, name string) (int, error)
	CreateRoleWithID(ctx context.Context, id int, name string) (int, error)
	ListRoles(ctx context.Context) ([]*Role, error)
	GetRoleById(ctx context.Context, id int) (*Role, error)

	SetRoleAuth(ctx context.Context, tx *gorm.DB, id int, authIdList []int) error
	ListRoleAuth(ctx context.Context, id int) (*RoleAuthInfo, error)
}

type Role struct {
	ID   int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name string `gorm:"type:varchar(128);NOT NULL;column:name"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type RoleAuth struct {
	ID     int `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	RoleId int `gorm:"type:int;NOT NULL;column:role_id"`
	AuthId int `gorm:"type:int;NOT NULL;column:auth_id"`
}

type RoleAuthInfo struct {
	RoleID    int
	Auth 	[]*entity.Auth
}

type DBRoleModel struct{}

func (d *DBRoleModel) CreateRole(ctx context.Context, tx *gorm.DB, name string) (int, error) {
	now := time.Now()
	role := &Role{
		Name:      name,
		UpdatedAt: &now,
		CreatedAt: &now,
	}
	err := tx.Create(role).Error
	if err != nil {
		return -1, err
	}
	return role.ID, nil
}

func (d *DBRoleModel) CreateRoleWithID(ctx context.Context, id int, name string) (int, error) {
	now := time.Now()
	role := &Role{
		ID:        id,
		Name:      name,
		UpdatedAt: &now,
		CreatedAt: &now,
	}
	err := db.Get().Create(role).Error
	if err != nil {
		return -1, err
	}
	return role.ID, nil
}

func (d *DBRoleModel) GetRoleById(ctx context.Context, id int) (*Role, error) {
	role := new(Role)
	err := db.Get().Where(&Role{ID: id}).First(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}
func (d *DBRoleModel) ListRoles(ctx context.Context) ([]*Role, error) {
	result := make([]*Role, 0)
	err := db.Get().Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBRoleModel) SetRoleAuth(ctx context.Context, tx *gorm.DB, id int, authIdList []int) error {
	for i := range authIdList {
		ra := &RoleAuth{
			RoleId: id,
			AuthId: authIdList[i],
		}
		err := tx.Create(ra).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DBRoleModel) ListRoleAuth(ctx context.Context, id int) (*RoleAuthInfo, error) {
	ras := make([]*RoleAuth, 0)
	err := db.Get().Where(&RoleAuth{RoleId: id}).Find(&ras).Error
	if err != nil {
		return nil, err
	}
	ret := new(RoleAuthInfo)
	ret.RoleID = id

	authIds := make([]int, len(ras))
	for i := range ras {
		authIds[i] = ras[i].AuthId
	}
	auths, err := GetAuthModel().ListAuthByIDs(ctx, authIds)
	if err != nil {
		return nil, err
	}
	for i := range auths {
		ret.Auth = append(ret.Auth, &entity.Auth{
			ID:   auths[i].ID,
			Name: auths[i].Name,
		})
	}

	return ret, nil
}

var (
	_roleModel     *DBRoleModel
	_roleModelOnce sync.Once
)

func GetRoleModel() IRolesModel {
	_roleModelOnce.Do(func() {
		_roleModel = new(DBRoleModel)
	})
	return _roleModel
}
