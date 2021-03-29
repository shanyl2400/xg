package da

import (
	"context"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type IOrderSourceModel interface {
	CreateOrderSources(ctx context.Context, name string) (int, error)
	CreateOrderSourcesWithID(ctx context.Context, id int, name string) error
	ListOrderSources(ctx context.Context) ([]*OrderSource, error)
	GetOrderSourceById(ctx context.Context, orderSourceId int) (*OrderSource, error)
	UpdateOrderSourceByID(ctx context.Context, tx *gorm.DB, id int, name string) error

	DeleteOrderSourceByID(ctx context.Context, tx *gorm.DB, id int) error
}

type OrderSource struct {
	ID   int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name string `gorm:"type:varchar(255);NOT NULL;column:org_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type DBOrderSourceModel struct {
}

func (d *DBOrderSourceModel) CreateOrderSourcesWithID(ctx context.Context, id int, name string) error {
	os := OrderSource{
		ID:   id,
		Name: name,
	}
	err := db.Get().Create(&os).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *DBOrderSourceModel) CreateOrderSources(ctx context.Context, name string) (int, error) {
	os := OrderSource{
		Name: name,
	}
	err := db.Get().Create(&os).Error
	if err != nil {
		return -1, err
	}
	return os.ID, nil
}

func (d *DBOrderSourceModel) UpdateOrderSourceByID(ctx context.Context, tx *gorm.DB, id int, name string) error {
	now := time.Now()
	err := tx.Model(OrderSource{}).Where(&OrderSource{ID: id}).Updates(OrderSource{Name: name, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderSourceModel) DeleteOrderSourceByID(ctx context.Context, tx *gorm.DB, id int) error {
	os, err := d.GetOrderSourceById(ctx, id)
	if err != nil {
		return err
	}
	now := time.Now()
	os.DeletedAt = &now

	err = tx.Model(OrderSource{}).Where(&OrderSource{ID: id}).Updates(os).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *DBOrderSourceModel) ListOrderSources(ctx context.Context) ([]*OrderSource, error) {
	result := make([]*OrderSource, 0)
	err := db.Get().Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DBOrderSourceModel) GetOrderSourceById(ctx context.Context, orderSourceId int) (*OrderSource, error) {
	orderSource := new(OrderSource)
	err := db.Get().Where(&OrderSource{ID: orderSourceId}).First(&orderSource).Error
	if err != nil {
		return nil, err
	}
	return orderSource, nil
}

var (
	_orderSourceModel     IOrderSourceModel
	_orderSourceModelOnce sync.Once
)

func GetOrderSourceModel() IOrderSourceModel {
	_orderSourceModelOnce.Do(func() {
		_orderSourceModel = new(DBOrderSourceModel)
	})
	return _orderSourceModel
}
