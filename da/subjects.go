package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"
)

type ISubjectModel interface{
	CreateSubject(ctx context.Context, subject Subject) (int, error)

	GetSubjectById(ctx context.Context, id int) (*Subject, error)
	SearchSubject(ctx context.Context, s SearchSubjectCondition) ([]*Subject, error)
}

type Subject struct {
	ID int	`gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Level int `gorm:"type:int;NOT NULL;column:level"`
	ParentId int	`gorm:"type:int;NOT NULL;column:parent_id"`
	Name string	 `gorm:"type:varchar(128);NOT NULL;column:name"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SearchSubjectCondition struct {
	IDList []int
	Level int
	ParentId int
}

func (s SearchSubjectCondition) GetConditions()(string, []interface{}){
	wheres := make([]string, 0)
	values := make([]interface{}, 0)
	if len(s.IDList) > 0 {
		wheres = append(wheres, "id in (?)")
		values = append(values, s.IDList)
	}
	if s.Level > 0 {
		wheres = append(wheres, "level = ?")
		values = append(values, s.Level)
	}
	if s.ParentId > 0 {
		wheres = append(wheres, "parent_id = ?")
		values = append(values, s.ParentId)
	}
	wheres = append(wheres, "deleted_at IS NULL")
	where := strings.Join(wheres, " and ")

	return where, values
}

type DBSubjectModel struct {}

func (d *DBSubjectModel) CreateSubject(ctx context.Context, subject Subject) (int, error) {
	now := time.Now()
	subject.CreatedAt = &now
	subject.UpdatedAt = &now
	err := db.Get().Create(&subject).Error
	if err != nil{
		return -1, err
	}
	return subject.ID, nil
}

func (d *DBSubjectModel) GetSubjectById(ctx context.Context, id int) (*Subject, error) {
	subject := new(Subject)
	err := db.Get().Where(&Subject{ID: id}).First(&subject).Error
	if err != nil{
		return nil, err
	}
	return subject, nil
}

func (d *DBSubjectModel) SearchSubject(ctx context.Context, s SearchSubjectCondition) ([]*Subject, error) {
	where, values := s.GetConditions()

	subjects := make([]*Subject, 0)
	err := db.Get().Where(where, values...).Find(&subjects).Error
	if err != nil{
		return nil, err
	}
	return subjects, nil
}
var(
	_subjectModel *DBSubjectModel
	_subjectModelOnce sync.Once
)

func GetSubjectModel() ISubjectModel{
	_subjectModelOnce.Do(func() {
		_subjectModel = new(DBSubjectModel)
	})
	return _subjectModel
}