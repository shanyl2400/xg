package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"
)

type IUsersModel interface {
	CreateUser(ctx context.Context, user User) (int, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, user User) error

	GetUserById(ctx context.Context, id int) (*User, error)
	SearchUsers(ctx context.Context, s SearchUserCondition) ([]*User, error)
}

type User struct {
	ID       int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name     string `gorm:"type:varchar(128);NOT NULL;column:name"`
	Password string `gorm:"type:varchar(128);NOT NULL;column:password"`
	OrgId    int    `gorm:"type:int;NOT NULL;column:org_id"`
	RoleId   int    `gorm:"type:int;NOT NULL;column:role_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SearchUserCondition struct {
	Name       string
	Password   string
	IDList     []int
	OrgIdList  []int
	RoleIdList []int
}

func (s SearchUserCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)
	if len(s.IDList) > 0 {
		wheres = append(wheres, "id in (?)")
		values = append(values, s.IDList)
	}
	if len(s.OrgIdList) > 0 {
		wheres = append(wheres, "org_id in (?)")
		values = append(values, s.OrgIdList)
	}
	if len(s.RoleIdList) > 0 {
		wheres = append(wheres, "role_id in (?)")
		values = append(values, s.RoleIdList)
	}
	if s.Name != "" {
		wheres = append(wheres, "name = ?")
		values = append(values, s.Name)
	}
	if s.Password != "" {
		wheres = append(wheres, "password = ?")
		values = append(values, s.Password)
	}

	//wheres = append(wheres, "deleted_at IS NULL")
	where := strings.Join(wheres, " and ")

	return where, values
}

type DBUsersModel struct{}

func (d *DBUsersModel) CreateUser(ctx context.Context, user User) (int, error) {
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	err := db.Get().Create(&user).Error
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (d *DBUsersModel) DeleteUser(ctx context.Context, id int) error {
	user := User{ID: id}
	err := db.Get().Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBUsersModel) UpdateUser(ctx context.Context, user User) error {
	err := db.Get().Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBUsersModel) GetUserById(ctx context.Context, id int) (*User, error) {
	user := new(User)
	err := db.Get().Where(&User{ID: id}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *DBUsersModel) SearchUsers(ctx context.Context, s SearchUserCondition) ([]*User, error) {
	where, values := s.GetConditions()

	users := make([]*User, 0)
	err := db.Get().Where(where, values...).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

var (
	_usersModel     *DBUsersModel
	_usersModelOnce sync.Once
)

func GetUsersModel() IUsersModel {
	_usersModelOnce.Do(func() {
		_usersModel = new(DBUsersModel)
	})
	return _usersModel
}
