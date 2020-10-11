package da

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"xg/db"
	"xg/entity"

	"github.com/jinzhu/gorm"
)

type IOrgModel interface {
	CreateOrg(ctx context.Context, tx *gorm.DB, org Org) (int, error)
	GetOrgById(ctx context.Context, tx *gorm.DB, id int) (*Org, error)
	ListOrgs(ctx context.Context) ([]*Org, error)
	UpdateOrg(ctx context.Context, tx *gorm.DB, id int, org Org) error
	CountOrgs(ctx context.Context, s SearchOrgsCondition) (int, error)

	ListOrgsByIDs(ctx context.Context, ids []int) ([]*Org, error)
	DeleteOrgById(ctx context.Context, tx *gorm.DB, ids []int) error

	GetOrgsByParentId(ctx context.Context, parentId int) ([]*Org, error)
	SearchOrgs(ctx context.Context, s SearchOrgsCondition) (int, []*Org, error)

	SearchOrgsWithDistance(ctx context.Context, s SearchOrgsCondition, l *entity.Coordinate) (int, []*OrgWithDistance, error)
}
type Org struct {
	ID         int     `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name       string  `gorm:"type:varchar(128);NOT NULL;column:name"`
	Subjects   string  `gorm:"type:varchar(255);NOT NULL;column:subjects"`
	Address    string  `gorm:"type:varchar(255);NOT NULL; column:address"`
	AddressExt string  `gorm:"type:varchar(255);NOT NULL; column:address_ext"`
	ParentID   int     `gorm:"type:int;NOT NULL;column:parent_id;index"`
	Telephone  string  `gorm:"type:varchar(64);NOT NULL; column:telephone"`
	Longitude  float64 `gorm:"type:double(9,6);NOT NULL; column:longitude; default:0"`
	Latitude   float64 `gorm:"type:double(9,6);NOT NULL; column:latitude; default:0"`

	Status        int    `gorm:"type:int;NOT NULL;column:status;index"`
	SupportRoleID string `gorm:"type:varchar(255);NOT NULL;column:support_role_ids"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type OrgWithDistance struct {
	Org
	Distance float64 `json:"distance"`
}

type DBOrgModel struct{}

func (d *DBOrgModel) CreateOrg(ctx context.Context, tx *gorm.DB, org Org) (int, error) {
	now := time.Now()
	org.CreatedAt = &now
	org.UpdatedAt = &now
	err := tx.Create(&org).Error
	if err != nil {
		return -1, err
	}
	return org.ID, nil
}

func (d *DBOrgModel) UpdateOrg(ctx context.Context, tx *gorm.DB, id int, org Org) error {
	now := time.Now()
	err := db.Get().Model(Org{}).Where(&Org{ID: id}).Updates(Org{
		Status:     org.Status,
		Subjects:   org.Subjects,
		Address:    org.Address,
		AddressExt: org.AddressExt,
		Longitude:  org.Longitude,
		Latitude:   org.Latitude,
		Telephone:  org.Telephone,
		UpdatedAt:  &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrgModel) GetOrgById(ctx context.Context, tx *gorm.DB, id int) (*Org, error) {
	org := new(Org)
	err := tx.Where(&Org{ID: id}).First(&org).Error
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

func (d *DBOrgModel) CountOrgs(ctx context.Context, s SearchOrgsCondition) (int, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(Org{}).Where(where, values...).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DBOrgModel) ListOrgsByIDs(ctx context.Context, ids []int) ([]*Org, error) {
	orgList := make([]*Org, 0)
	err := db.Get().Where("id IN (?)", ids).Find(&orgList).Error
	if err != nil {
		return nil, err
	}
	return orgList, nil
}

func (d *DBOrgModel) DeleteOrgById(ctx context.Context, tx *gorm.DB, ids []int) error {
	err := tx.Delete(&Org{}, ids).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *DBOrgModel) GetOrgsByParentId(ctx context.Context, parentId int) ([]*Org, error) {
	condition := SearchOrgsCondition{
		ParentIDs: []int{parentId},
	}
	where, values := condition.GetConditions()
	result := make([]*Org, 0)
	err := db.Get().Where(where, values...).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBOrgModel) SearchOrgs(ctx context.Context, s SearchOrgsCondition) (int, []*Org, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(Org{}).Where(where, values...).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	result := make([]*Org, 0)
	tx := db.Get().Where(where, values...)

	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err = tx.Find(&result).Error

	if err != nil {
		return 0, nil, err
	}
	return count, result, nil
}

func (d *DBOrgModel) SearchOrgsWithDistance(ctx context.Context, s SearchOrgsCondition, l *entity.Coordinate) (int, []*OrgWithDistance, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(Org{}).Where(where, values...).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	result := make([]*OrgWithDistance, 0)
	tx := db.Get().Where(where, values...)

	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	} else {
		tx = tx.Order("distance")
	}
	err = tx.Model(Org{}).Select([]string{
		"*",
		fmt.Sprintf("(st_distance(point(longitude,latitude),point(%v,%v))*111195/1000 ) as distance", l.Longitude, l.Latitude),
	}).Scan(&result).Error

	if err != nil {
		return 0, nil, err
	}
	return count, result, nil
}

type SearchOrgsCondition struct {
	Subjects  []string
	Address   string
	Status    []int
	StudentID int

	ParentIDs []int
	IsSubOrg  bool

	OrderBy  string
	Page     int
	PageSize int
}

func (s SearchOrgsCondition) GetConditions() (string, []interface{}) {
	var wheres []string
	var values []interface{}

	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}

	if len(s.Subjects) != 0 {
		partsWhere := make([]string, 0)
		for i := range s.Subjects {
			partsWhere = append(partsWhere, "subjects LIKE ?")
			values = append(values, "%"+s.Subjects[i]+"%")
		}
		where := "(" + strings.Join(partsWhere, " or ") + ")"
		wheres = append(wheres, where)
	}
	if s.Address != "" {
		wheres = append(wheres, "address LIKE ?")
		values = append(values, "%"+s.Address+"%")
	}

	if len(s.ParentIDs) > 0 {
		wheres = append(wheres, "parent_id IN (?)")
		values = append(values, s.ParentIDs)
	}
	if s.IsSubOrg {
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
