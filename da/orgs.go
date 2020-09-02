package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"
)

type IOrgModel interface {
	CreateOrg(ctx context.Context, org Org) (int, error)
	GetOrgById(ctx context.Context, id int) (*Org, error)
	ListOrgs(ctx context.Context) ([]*Org, error)
	UpdateOrg(ctx context.Context, id int, org Org) error
	CountOrgs(ctx context.Context) (int, error)

	SearchOrgs(ctx context.Context, s SearchOrgsCondition) (int, []*Org, error)
}

type Org struct {
	ID       int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name     string `gorm:"type:varchar(128);NOT NULL;column:name"`
	Subjects string `gorm:"type:varchar(255);NOT NULL;column:subjects"`
	Address  string `gorm:"type:varchar(255);NOT NULL; column:address"`
	ParentID int    `gorm:"type:int;NOT NULL;column:parent_id"`

	Status int `gorm:"type:int;NOT NULL;column:status"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type DBOrgModel struct{}

func (d *DBOrgModel) CreateOrg(ctx context.Context, org Org) (int, error) {
	now := time.Now()
	org.CreatedAt = &now
	org.UpdatedAt = &now
	err := db.Get().Create(&org).Error
	if err != nil {
		return -1, err
	}
	return org.ID, nil
}

func (d *DBOrgModel) UpdateOrg(ctx context.Context, id int, org Org) error {
	now := time.Now()
	err := db.Get().Where(&Org{ID: id}).Updates(Org{Status: org.Status, Subjects: org.Subjects, Address: org.Address, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrgModel) GetOrgById(ctx context.Context, id int) (*Org, error) {
	org := new(Org)
	err := db.Get().Where(&Org{ID: id}).First(&org).Error
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (d *DBOrgModel) ListOrgs(ctx context.Context) ([]*Org, error) {
	result := make([]*Org, 0)
	err := db.Get().Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBOrgModel) CountOrgs(ctx context.Context) (int, error) {
	count := 0
	err := db.Get().Model(Org{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DBOrgModel) SearchOrgs(ctx context.Context, s SearchOrgsCondition) (int, []*Org, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(Org{}).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	result := make([]*Org, 0)
	err = db.Get().Where(where, values...).Find(&result).Error
	if err != nil {
		return 0, nil, err
	}
	return count, result, nil
}

type SearchOrgsCondition struct {
	Subjects string
	Address  string
	Status   []int

	ParentIDs []int
	IsSubOrg bool
}

func (s SearchOrgsCondition) GetConditions() (string, []interface{}) {
	var wheres []string
	var values []interface{}

	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}

	if s.Subjects != "" {
		wheres = append(wheres, "subjects LIKE ?")
		values = append(values, "%"+s.Subjects+"%")
	}
	if s.Address != "" {
		wheres = append(wheres, "address LIKE ?")
		values = append(values, "%"+s.Address+"%")
	}

	if len(s.ParentIDs) > 0 {
		wheres = append(wheres, "parent_id IN (?)")
		values = append(values, s.ParentIDs)
	}
	if s.IsSubOrg{
		wheres = append(wheres, "parent_id != 0")
	}

	//wheres = append(wheres, "deleted_at IS NULL")
	where := strings.Join(wheres, " and ")

	return where, values
}

var (
	_orgModel     *DBOrgModel
	_orgModelOnce sync.Once
)

func GetOrgModel() IOrgModel {
	_orgModelOnce.Do(func() {
		_orgModel = new(DBOrgModel)
	})
	return _orgModel
}
