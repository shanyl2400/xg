package da

import (
	"context"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
	"time"
)

type IStatisticsModel interface {
	CreateStatisticsRecord(ctx context.Context, tx *gorm.DB, c *StatisticsRecord) (int ,error)
	UpdateStatisticsRecord(ctx context.Context, tx *gorm.DB, rid int, value int) error

	SearchStatisticsRecord(ctx context.Context, tx *gorm.DB, c SearchStatisticsRecordCondition)([]*StatisticsRecord, error)
}

type StatisticsRecord struct {
	ID       int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`

	Key     string `gorm:"type:varchar(128);NOT NULL;column:key"`
	Value int    `gorm:"type:int;NOT NULL;column:value"`

	Year int `gorm:"type:int;NOT NULL;column:year"`
	Month int `gorm:"type:int;NOT NULL;column:month"`

	Author int `gorm:"type:int;NOT NULL;column:author"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SearchStatisticsRecordCondition struct {
	Key string `json:"key"`
	Year int `json:"year"`
	Month int `json:"month"`
	Author int `json:"author"`
}

func (s SearchStatisticsRecordCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if s.Key != "" {
		wheres = append(wheres, "`key` = ?")
		values = append(values, s.Key)
	}
	if s.Year > 0 {
		wheres = append(wheres, "`year` = ?")
		values = append(values, s.Year)
	}
	if s.Month > 0 {
		wheres = append(wheres, "`month` = ?")
		values = append(values, s.Month)
	}
	if s.Author > 0 {
		wheres = append(wheres, "`author` = ?")
		values = append(values, s.Month)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type DBStatisticsModel struct {

}

func (s *DBStatisticsModel) CreateStatisticsRecord(ctx context.Context, tx *gorm.DB, c *StatisticsRecord) (int ,error){
	now := time.Now()
	c.CreatedAt = &now
	c.UpdatedAt = &now
	err := tx.Create(&c).Error
	if err != nil {
		return -1, err
	}
	return c.ID, nil
}
func (s *DBStatisticsModel) UpdateStatisticsRecord(ctx context.Context, tx *gorm.DB, rid int, value int) error{
	record := new(StatisticsRecord)
	now := time.Now()
	record.UpdatedAt = &now
	record.Value = value
	whereRecord := &StatisticsRecord{ID: rid}
	err := tx.Model(StatisticsRecord{}).Where(whereRecord).Updates(&record).Error
	if err != nil {
		return err
	}
	return nil

}

func (s *DBStatisticsModel) SearchStatisticsRecord(ctx context.Context, tx *gorm.DB, c SearchStatisticsRecordCondition)([]*StatisticsRecord, error){
	records := make([]*StatisticsRecord, 0)
	where, values := c.GetConditions()
	err := tx.Where(where, values...).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

var (
	_statisticsRecordModel     *DBStatisticsModel
	_statisticsRecordModelOnce sync.Once
)

func GetStatisticsRecordModel() IStatisticsModel {
	_statisticsRecordModelOnce.Do(func() {
		_statisticsRecordModel = new(DBStatisticsModel)
	})
	return _statisticsRecordModel
}
