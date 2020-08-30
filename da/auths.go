package da

import (
	"context"
	"sync"
	"xg/db"
)

type IAuthModel interface {
	CreateAuth(ctx context.Context, name string) (int, error)
	CreateAuthWithID(ctx context.Context, id int, name string) error
	ListAuth(ctx context.Context) ([]*Auth, error)
	ListAuthByIDs(ctx context.Context, ids []int) ([]*Auth, error)
}

type Auth struct {
	ID   int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name string `gorm:"type:varchar(128);NOT NULL;column:name"`
}

type DBAuthModel struct {
}

func (d *DBAuthModel) CreateAuth(ctx context.Context, name string) (int, error) {
	auth := Auth{
		Name: name,
	}
	err := db.Get().Create(&auth).Error
	if err != nil {
		return -1, err
	}
	return auth.ID, nil
}

func (d *DBAuthModel) CreateAuthWithID(ctx context.Context, id int, name string) error {
	auth := Auth{
		ID:   id,
		Name: name,
	}
	err := db.Get().Create(&auth).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBAuthModel) ListAuth(ctx context.Context) ([]*Auth, error) {
	result := make([]*Auth, 0)
	err := db.Get().Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (d *DBAuthModel) ListAuthByIDs(ctx context.Context, ids []int) ([]*Auth, error) {
	authList := make([]*Auth, 0)
	err := db.Get().Where("id IN (?)", ids).Find(&authList).Error
	if err != nil {
		return nil, err
	}
	return authList, nil
}

var (
	_authModel     *DBAuthModel
	_authModelOnce sync.Once
)

func GetAuthModel() IAuthModel {
	_authModelOnce.Do(func() {
		_authModel = new(DBAuthModel)
	})
	return _authModel
}
