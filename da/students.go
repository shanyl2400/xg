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

type IStudentsModel interface {
	CreateStudent(ctx context.Context, tx *gorm.DB, student Student) (int, error)
	UpdateStudent(ctx context.Context, tx *gorm.DB, id int, student Student) error
	GetStudentById(ctx context.Context, id int) (*Student, error)
	SearchStudents(ctx context.Context, s SearchStudentCondition) (int, []*Student, error)
	CountStudents(ctx context.Context) (int, error)

	StatisticStudentsWithStatus(ctx context.Context, groupby string, limit int, s SearchStudentCondition) ([]*entity.GroupbyStatisticEntity, error)
	StatisticStudents(ctx context.Context, groupby string, limit int, s SearchStudentCondition) ([]*entity.GroupbyStatisticEntity, error)

	ReplaceStudentOrderSource(ctx context.Context, tx *gorm.DB, oldOrderSource, newOrderSource int) error
}

type Student struct {
	ID            int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name          string `gorm:"type:varchar(128);NOT NULL;column:name;index"`
	Gender        bool   `gorm:"type:int;NOT NULL;column:gender"`
	Telephone     string `gorm:"type:varchar(20);NOT NULL;column:telephone;index"`
	Address       string `gorm:"type:varchar(128);NOT NULL;column:address"`
	AddressExt    string `gorm:"type:varchar(256);NOT NULL;column:address_ext"`
	Email         string `gorm:"type:varchar(128);NOT NULL;column:email"`
	IntentSubject string `gorm:"type:varchar(255);NOT NULL;column:intent_subject"`
	Status        int    `gorm:"type:int;NOT NULL;column:status;index"`
	Note          string `gorm:"type:text;NOT NULL;column:note"`
	OrderSourceID int    `gorm:"type:int;NOT NULL;column:order_source_id"`

	OrderSourceExt string `gorm:"type:varchar(256);NULL;column:order_source_ext"`

	OrderCount int `gorm:"type:int;NOT NULL;column:order_count"`

	Longitude float64 `gorm:"type:double(9,6);NOT NULL; column:longitude; default:0"`
	Latitude  float64 `gorm:"type:double(9,6);NOT NULL; column:latitude; default:0"`

	AuthorID int `gorm:"type:int;NOT NULL;column:author_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type StudentNote struct {
	ID   int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Note string `gorm:"type:text;NOT NULL;column:note"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type DBStudentsModel struct{}

func (d *DBStudentsModel) CreateStudent(ctx context.Context, tx *gorm.DB, student Student) (int, error) {
	now := time.Now()
	student.CreatedAt = &now
	student.UpdatedAt = &now
	err := tx.Create(&student).Error
	if err != nil {
		return -1, err
	}
	return student.ID, nil
}
func (d *DBStudentsModel) UpdateStudent(ctx context.Context, tx *gorm.DB, id int, student Student) error {
	now := time.Now()
	student.UpdatedAt = &now
	//student.ID = id
	err := tx.Model(Student{ID: id}).Updates(&student).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBStudentsModel) ReplaceStudentOrderSource(ctx context.Context, tx *gorm.DB, oldOrderSource, newOrderSource int) error {
	err := tx.Model(&Student{}).Where(" order_source_id = ?", oldOrderSource).Updates(Student{OrderSourceID: newOrderSource}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBStudentsModel) GetStudentById(ctx context.Context, id int) (*Student, error) {
	student := new(Student)
	err := db.Get().Where(&Student{ID: id}).First(&student).Error
	if err != nil {
		return nil, err
	}
	return student, nil
}
func (d *DBStudentsModel) StatisticStudentsWithStatus(ctx context.Context, groupby string, limit int, s SearchStudentCondition) ([]*entity.GroupbyStatisticEntity, error) {
	where, values := s.GetConditions()
	tx := db.Get().Table("students").Select(fmt.Sprintf("%v as id, status, count(*) as cnt", groupby)).Where(where, values...).Group(groupby + ",status")
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	tx = tx.Order("cnt desc")
	entities := make([]*entity.GroupbyStatisticEntity, 0)
	err := tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *DBStudentsModel) StatisticStudents(ctx context.Context, groupby string, limit int, s SearchStudentCondition) ([]*entity.GroupbyStatisticEntity, error) {
	where, values := s.GetConditions()
	tx := db.Get().Table("students").Select(fmt.Sprintf("%v as id, count(*) as cnt", groupby)).Where(where, values...).Group(groupby)
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	tx = tx.Order("cnt desc")
	entities := make([]*entity.GroupbyStatisticEntity, 0)
	err := tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
func (d *DBStudentsModel) SearchStudents(ctx context.Context, s SearchStudentCondition) (int, []*Student, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(Student{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	students := make([]*Student, 0)
	tx := db.Get().Where(where, values...)
	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err = tx.Find(&students).Error
	if err != nil {
		return 0, nil, err
	}
	return total, students, nil
}

func (d *DBStudentsModel) CountStudents(ctx context.Context) (int, error) {
	count := 0
	err := db.Get().Model(Student{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func parsePage(page, pageSize int) (int, int) {
	return (page - 1) * pageSize, pageSize
}

type SearchStudentCondition struct {
	StudentIDList   []int  `json:"student_id_list"`
	Name            string `json:"name"`
	Telephone       string `json:"telephone"`
	Address         string `json:"address"`
	IntentString    string `json:"intent_string"`
	Status          []int  `json:"status"`
	NoDispatchOrder bool   `json:"no_dispatch_order"`
	Keywords        string `json:"keywords"`

	AuthorIDList []int `json:"author_id_list"`

	OrderSourceIDs []int      `json:"order_source_ids"`
	CreatedStartAt *time.Time `json:"created_start_at"`
	CreatedEndAt   *time.Time `json:"created_end_at"`

	OrderBy  string `json:"order_by"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}

func (s SearchStudentCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.StudentIDList) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.StudentIDList)
	}
	if s.Name != "" {
		wheres = append(wheres, "name LIKE ?")
		values = append(values, "%"+s.Name+"%")
	}
	if s.Keywords != "" {
		wheres = append(wheres, "(name LIKE ? OR telephone LIKE ?)")
		values = append(values, s.Keywords+"%")
		values = append(values, s.Keywords+"%")
	}
	if s.Telephone != "" {
		wheres = append(wheres, "telephone LIKE ?")
		values = append(values, s.Telephone)
	}
	if s.Address != "" {
		wheres = append(wheres, "address LIKE ?")
		//values = append(values, "%"+s.Address+"%")
		values = append(values, s.Address+"%")
	}
	if len(s.OrderSourceIDs) > 0 {
		wheres = append(wheres, "order_source_id IN (?)")
		values = append(values, s.OrderSourceIDs)
	}

	if s.IntentString != "" {
		wheres = append(wheres, "intent_subject LIKE ?")
		values = append(values, "%"+s.IntentString+"%")
	}
	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}
	if s.NoDispatchOrder {
		wheres = append(wheres, "order_count = 0")
	}

	if len(s.AuthorIDList) > 0 {
		wheres = append(wheres, "author_id IN (?)")
		values = append(values, s.AuthorIDList)
	}
	if s.CreatedStartAt != nil {
		wheres = append(wheres, "created_at >= ?")
		values = append(values, s.CreatedStartAt)
	}
	if s.CreatedEndAt != nil {
		wheres = append(wheres, "created_at <= ?")
		values = append(values, s.CreatedEndAt)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

var (
	_studentModel     *DBStudentsModel
	_studentModelOnce sync.Once
)

func GetStudentModel() IStudentsModel {
	_studentModelOnce.Do(func() {
		_studentModel = new(DBStudentsModel)
	})
	return _studentModel
}
