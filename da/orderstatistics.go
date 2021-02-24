package da

import (
	"context"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
	"time"
)

type IOrderStatisticsModel interface {
	CreateOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, c *OrderStatisticsRecord) (int, error)
	UpdateOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, rid int, value float64, count int) error

	SearchOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, c SearchOrderStatisticsRecordCondition) ([]*OrderStatisticsRecord, error)
}

type OrderStatisticsRecord struct {
	ID int `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`

	Key   string `gorm:"type:varchar(128);NOT NULL;column:key;index"`

	Year  int `gorm:"type:int;NOT NULL;column:year"`
	Month int `gorm:"type:int;NOT NULL;column:month"`
	Date int `gorm:"type:int;NOT NULL;column:day"`

	Author int `gorm:"type:int;NOT NULL;column:author;index"`
	PublisherId int `gorm:"type:int;NOT NULL;column:publisher_id;index"`
	OrgId int `gorm:"type:int;NOT NULL;column:org_id;index"`
	OrderSource int `gorm:"type:int;NOT NULL;column:order_source;index"`

	Value float64    `gorm:"type:DECIMAL(11,2);NOT NULL;column:value"`
	Count int    `gorm:"type:int;NOT NULL;column:count"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}


type SearchOrderStatisticsRecordCondition struct {
	Key    string `json:"key"`
	Year   []int    `json:"year"`
	Month  []int    `json:"month"`
	Date  []int    `json:"day"`
	Author []int    `json:"author"`
	OrgId []int    `json:"org_id"`
	PublisherId []int `json:"publisher_id"`
	OrderSource []int    `json:"order_source"`
}

func (s SearchOrderStatisticsRecordCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if s.Key != "" {
		wheres = append(wheres, "`key` = ?")
		values = append(values, s.Key)
	}
	if len(s.Year) > 0 {
		wheres = append(wheres, "`year` IN (?)")
		values = append(values, s.Year)
	}
	if len(s.Month) > 0 {
		wheres = append(wheres, "`month` IN (?)")
		values = append(values, s.Month)
	}
	if len(s.Date) > 0 {
		wheres = append(wheres, "`day` IN (?)")
		values = append(values, s.Date)
	}
	if len(s.Author) > 0 {
		wheres = append(wheres, "`author` IN (?)")
		values = append(values, s.Author)
	}
	if len(s.OrgId) > 0 {
		wheres = append(wheres, "`org_id` IN (?)")
		values = append(values, s.OrgId)
	}
	if len(s.OrderSource) > 0 {
		wheres = append(wheres, "`order_source` IN (?)")
		values = append(values, s.OrderSource)
	}
	if len(s.PublisherId) > 0 {
		wheres = append(wheres, "`publisher_id` IN (?)")
		values = append(values, s.PublisherId)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type DBOrderStatisticsModel struct {
}

func (s *DBOrderStatisticsModel) CreateOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, c *OrderStatisticsRecord) (int, error) {
	now := time.Now()
	c.CreatedAt = &now
	c.UpdatedAt = &now
	err := tx.Create(&c).Error
	if err != nil {
		return -1, err
	}
	return c.ID, nil
}
func (s *DBOrderStatisticsModel) UpdateOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, rid int, value float64, count int) error {
	record := new(OrderStatisticsRecord)
	now := time.Now()
	record.UpdatedAt = &now
	record.Value = value
	record.Count = count
	whereRecord := &OrderStatisticsRecord{ID: rid}
	err := tx.Model(OrderStatisticsRecord{}).Where(whereRecord).Updates(&record).Error
	if err != nil {
		return err
	}
	return nil

}

func (s *DBOrderStatisticsModel) SearchOrderStatisticsRecord(ctx context.Context, tx *gorm.DB, c SearchOrderStatisticsRecordCondition) ([]*OrderStatisticsRecord, error) {
	records := make([]*OrderStatisticsRecord, 0)
	where, values := c.GetConditions()
	err := tx.Where(where, values...).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

var (
	_orderStatisticsRecordModel     *DBOrderStatisticsModel
	_orderStatisticsRecordModelOnce sync.Once
)

func GetOrderStatisticsRecordModel() IOrderStatisticsModel {
	_orderStatisticsRecordModelOnce.Do(func() {
		_orderStatisticsRecordModel = new(DBOrderStatisticsModel)
	})
	return _orderStatisticsRecordModel
}
