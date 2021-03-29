package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type SettlementRecord struct {
	ID            int       `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	StartAt       time.Time `gorm:"type:datetime;NOT NULL;column:start_at;index"`
	EndAt         time.Time `gorm:"type:datetime;NOT NULL;column:end_at;index"`
	SuccessOrders string    `gorm:"type:varchar(1024);NOT NULL;column:success_orders"`
	FailedOrders  string    `gorm:"type:varchar(1024);NOT NULL;column:failed_orders"`
	Amount        float64   `gorm:"type:DECIMAL(11,2);NOT NULL;column:amount"`
	Status        int       `gorm:"type:int;NOT NULL;column:status"`
	Invoice       int       `gorm:"type:text;NULL;column:invoice"`
	AuthorID      int       `gorm:"type:int;NOT NULL; column:author_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type ISettlementModel interface {
	CreateSettlement(ctx context.Context, tx *gorm.DB, record SettlementRecord) error
	SearchSettlements(ctx context.Context, s SearchSettlementsCondition) (int, []*SettlementRecord, error)
}

type SearchSettlementsCondition struct {
	IDs     []int
	StartAt *time.Time
	EndAt   *time.Time

	OrderBy  string
	Page     int
	PageSize int
}

func (s SearchSettlementsCondition) GetConditions() (string, []interface{}) {
	var wheres []string
	var values []interface{}

	if len(s.IDs) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.IDs)
	}

	if s.StartAt != nil {
		wheres = append(wheres, "start_at > ?")
		values = append(values, s.StartAt)
	}
	if s.EndAt != nil {
		wheres = append(wheres, "end_at > ?")
		values = append(values, s.EndAt)
	}
	//wheres = append(wheres, "deleted_at IS NULL")
	where := strings.Join(wheres, " and ")
	return where, values
}

type SettlementModel struct{}

func (d *SettlementModel) CreateSettlement(ctx context.Context, tx *gorm.DB, record SettlementRecord) error {
	now := time.Now()
	record.CreatedAt = &now
	record.UpdatedAt = &now
	err := tx.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *SettlementModel) SearchSettlements(ctx context.Context, s SearchSettlementsCondition) (int, []*SettlementRecord, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(SettlementRecord{}).Where(where, values...).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	result := make([]*SettlementRecord, 0)
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

var (
	_settlementModel     ISettlementModel
	_settlementModelOnce sync.Once
)

func GetSettlementModel() ISettlementModel {
	_settlementModelOnce.Do(func() {
		_settlementModel = new(SettlementModel)
	})
	return _settlementModel
}
