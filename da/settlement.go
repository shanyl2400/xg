package da

import (
	"context"
	"strings"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type CommissionSettlementRecord struct {
	ID             int     `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	OrderID        int     `gorm:"type:int;NOT NULL;column:order_id;index"`
	PaymentID      int     `gorm:"type:int;NOT NULL;column:payment_id;index"`
	Amount         float64 `gorm:"type:DECIMAL(11,2);NOT NULL;column:amount"`
	Commission     float64 `gorm:"type:DECIMAL(11,2);NOT NULL;column:commission"`
	SettlementNote string  `gorm:"type:text;NULL;column:settlement_note;COLLATION(utf8_general_ci)"`
	Status         int     `gorm:"type:int;NOT NULL;column:status"`
	AuthorID       int     `gorm:"type:int;NOT NULL; column:author_id"`
	Note           string  `gorm:"type:text;NULL;column:note;COLLATION(utf8_general_ci)"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SettlementRecord struct {
	ID            int       `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	StartAt       time.Time `gorm:"type:datetime;NOT NULL;column:start_at;index"`
	EndAt         time.Time `gorm:"type:datetime;NOT NULL;column:end_at;index"`
	SuccessOrders string    `gorm:"type:varchar(1024);NOT NULL;column:success_orders"`
	FailedOrders  string    `gorm:"type:varchar(1024);NOT NULL;column:failed_orders"`
	Amount        float64   `gorm:"type:DECIMAL(11,2);NOT NULL;column:amount"`
	Status        int       `gorm:"type:int;NOT NULL;column:status"`
	Invoice       string    `gorm:"type:text;NULL;column:invoice"`
	AuthorID      int       `gorm:"type:int;NOT NULL; column:author_id"`
	Commission    float64   `gorm:"type:DECIMAL(11,2);NOT NULL;column:commission"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type ISettlementModel interface {
	CreateSettlement(ctx context.Context, tx *gorm.DB, record SettlementRecord) error
	SearchSettlements(ctx context.Context, s SearchSettlementsCondition) (int, []*SettlementRecord, error)
	CreateCommissionSettlement(ctx context.Context, tx *gorm.DB, record CommissionSettlementRecord) error
	GetCommissionSetSettlementByID(ctx context.Context, tx *gorm.DB, id int) (*CommissionSettlementRecord, error)
	UpdateCommissionSettlement(ctx context.Context, tx *gorm.DB, id int, record CommissionSettlementRecord) error
	SearchCommissionSettlements(ctx context.Context, s SearchCommissionSettlementsCondition) (int, []*CommissionSettlementRecord, error)
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

type SearchCommissionSettlementsCondition struct {
	IDs        []int
	PaymentIDs []int
	OrderIDs   []int

	OrderBy  string
	Page     int
	PageSize int
}

func (s SearchCommissionSettlementsCondition) GetConditions() (string, []interface{}) {
	var wheres []string
	var values []interface{}

	if len(s.IDs) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.IDs)
	}
	if len(s.PaymentIDs) > 0 {
		wheres = append(wheres, "payment_id IN (?)")
		values = append(values, s.PaymentIDs)
	}

	if len(s.OrderIDs) > 0 {
		wheres = append(wheres, "order_id IN (?)")
		values = append(values, s.OrderIDs)
	}

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

func (d *SettlementModel) CreateCommissionSettlement(ctx context.Context, tx *gorm.DB, record CommissionSettlementRecord) error {
	now := time.Now()
	record.CreatedAt = &now
	record.UpdatedAt = &now
	err := tx.Create(&record).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *SettlementModel) GetCommissionSetSettlementByID(ctx context.Context, tx *gorm.DB, id int) (*CommissionSettlementRecord, error) {
	res := new(CommissionSettlementRecord)
	err := tx.Where(&CommissionSettlementRecord{ID: id}).First(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (d *SettlementModel) UpdateCommissionSettlement(ctx context.Context, tx *gorm.DB, id int, record CommissionSettlementRecord) error {
	now := time.Now()
	record.ID = id
	record.UpdatedAt = &now
	err := tx.Model(CommissionSettlementRecord{}).Save(record).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *SettlementModel) SearchCommissionSettlements(ctx context.Context, s SearchCommissionSettlementsCondition) (int, []*CommissionSettlementRecord, error) {
	where, values := s.GetConditions()
	count := 0
	err := db.Get().Model(CommissionSettlementRecord{}).Where(where, values...).Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	result := make([]*CommissionSettlementRecord, 0)
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
