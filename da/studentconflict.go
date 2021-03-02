package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type StudentConflict struct {
	ID        int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Telephone string `gorm:"type:varchar(20);NOT NULL;column:telephone;index"`
	Status    int    `gorm:"type:int;NOT NULL;column:status;index"`
	Total     int    `gorm:"type:int;NOT NULL;column:total"`

	AuthorID int `gorm:"type:int;NOT NULL;column:author_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type IStudentConflictModel interface {
	CreateStudentConflict(ctx context.Context, tx *gorm.DB, record StudentConflict) (int, error)
	UpdateStudentConflict(ctx context.Context, tx *gorm.DB, id int, record StudentConflict) error
	GetStudentConflictByID(ctx context.Context, id int) (*StudentConflict, error)
	SearchStudentConflicts(ctx context.Context, s SearchStudentConflictCondition) (int, []*StudentConflict, error)
}
type DBStudentConflictModel struct {
}

func (d *DBStudentConflictModel) CreateStudentConflict(ctx context.Context, tx *gorm.DB, record StudentConflict) (int, error) {
	now := time.Now()
	record.CreatedAt = &now
	record.UpdatedAt = &now
	err := tx.Create(&record).Error
	if err != nil {
		return -1, err
	}
	return record.ID, nil
}
func (d *DBStudentConflictModel) UpdateStudentConflict(ctx context.Context, tx *gorm.DB, id int, record StudentConflict) error {
	now := time.Now()
	record.UpdatedAt = &now
	err := tx.Model(StudentConflict{ID: id}).Updates(&record).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *DBStudentConflictModel) GetStudentConflictByID(ctx context.Context, id int) (*StudentConflict, error) {
	record := new(StudentConflict)
	err := db.Get().Where(&StudentConflict{ID: id}).First(&record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (d *DBStudentConflictModel) SearchStudentConflicts(ctx context.Context, s SearchStudentConflictCondition) (int, []*StudentConflict, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(StudentConflict{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	records := make([]*StudentConflict, 0)
	tx := db.Get().Model(StudentConflict{}).Where(where, values...)
	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err = tx.Find(&records).Error
	if err != nil {
		return 0, nil, err
	}
	return total, records, nil
}

type SearchStudentConflictCondition struct {
	IDList    []int  `json:"student_id_list"`
	Telephone string `json:"telephone"`
	Status    []int  `json:"status"`

	AuthorIDList []int `json:"author_id_list"`

	OrderBy  string `json:"order_by"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

func (s SearchStudentConflictCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.IDList) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.IDList)
	}
	if s.Telephone != "" {
		wheres = append(wheres, "telephone LIKE ?")
		values = append(values, s.Telephone)
	}
	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}

	if len(s.AuthorIDList) > 0 {
		wheres = append(wheres, "author_id IN (?)")
		values = append(values, s.AuthorIDList)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

var (
	_studentConflictModel     *DBStudentConflictModel
	_studentConflictModelOnce sync.Once
)

func GetStudentConflictModel() IStudentConflictModel {
	_studentConflictModelOnce.Do(func() {
		_studentConflictModel = new(DBStudentConflictModel)
	})
	return _studentConflictModel
}
