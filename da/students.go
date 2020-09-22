package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type IStudentsModel interface {
	CreateStudent(ctx context.Context, tx *gorm.DB, student Student) (int, error)
	UpdateStudent(ctx context.Context, id int, student Student) error
	GetStudentById(ctx context.Context, id int) (*Student, error)
	SearchStudents(ctx context.Context, s SearchStudentCondition) (int, []*Student, error)
	CountStudents(ctx context.Context) (int, error)
}

type Student struct {
	ID            int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name          string `gorm:"type:varchar(128);NOT NULL;column:name"`
	Gender        bool   `gorm:"type:int;NOT NULL;column:gender"`
	Telephone     string `gorm:"type:varchar(20);NOT NULL;column:telephone"`
	Address       string `gorm:"type:varchar(128);NOT NULL;column:address"`
	AddressExt	  string `gorm:"type:varchar(256);NOT NULL;column:address_ext"`
	Email         string `gorm:"type:varchar(128);NOT NULL;column:email"`
	IntentSubject string `gorm:"type:varchar(255);NOT NULL;column:intent_subject"`
	Status        int    `gorm:"type:int;NOT NULL;column:status;index"`
	Note          string `gorm:"type:text;NOT NULL;column:note"`
	OrderSourceID int    `gorm:"type:int;NOT NULL;column:order_source_id"`

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
func (d *DBStudentsModel) UpdateStudent(ctx context.Context, id int, student Student) error {
	now := time.Now()
	student.UpdatedAt = &now
	//student.ID = id
	err := db.Get().Model(Student{ID: id}).Updates(&student).Error
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
	StudentIDList []int  `json:"student_id_list"`
	Name          string `json:"name"`
	Telephone     string `json:"telephone"`
	Address       string `json:"address"`
	IntentString  string `json:"intent_string"`
	Status        int    `json:"status"`

	AuthorIDList []int `json:"author_id_list"`

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
	if s.Telephone != "" {
		wheres = append(wheres, "telephone LIKE ?")
		values = append(values, s.Telephone)
	}
	if s.Address != "" {
		wheres = append(wheres, "address LIKE ?")
		values = append(values, "%"+s.Address+"%")
	}

	if s.IntentString != "" {
		wheres = append(wheres, "intent_subject LIKE ?")
		values = append(values, "%"+s.IntentString+"%")
	}
	if s.Status > 0 {
		wheres = append(wheres, "status=?")
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
	_studentModel     *DBStudentsModel
	_studentModelOnce sync.Once
)

func GetStudentModel() IStudentsModel {
	_studentModelOnce.Do(func() {
		_studentModel = new(DBStudentsModel)
	})
	return _studentModel
}
