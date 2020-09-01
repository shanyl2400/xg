package da

import (
	"context"
	"sync"
	"time"
	"xg/db"
)

type IOrderSourceModel interface{
	CreateOrderSources(ctx context.Context, name string)(int, error)
	ListOrderSources(ctx context.Context) ([]*OrderSource, error)
}

type OrderSource struct {
	ID int	`gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	Name string	 `gorm:"type:varchar(255);NOT NULL;column:org_id"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type DBOrderSourceModel struct {
}

func (d *DBOrderSourceModel) CreateOrderSources(ctx context.Context, name string) (int, error){
	os := OrderSource{
		Name: name,
	}
	err := db.Get().Create(&os).Error
	if err != nil{
		return -1, err
	}
	return os.ID, nil
}
func (d *DBOrderSourceModel) ListOrderSources(ctx context.Context) ([]*OrderSource, error){
	result := make([]*OrderSource, 0)
	err := db.Get().Find(&result).Error
	if err != nil{
		return nil, err
	}
	return result, nil
}

var(
	_orderSourceModel IOrderSourceModel
	_orderSourceModelOnce sync.Once
)

func GetOrderSourceModel() IOrderSourceModel{
	_orderSourceModelOnce.Do(func() {
		_orderSourceModel = new(DBOrderSourceModel)
	})
	return _orderSourceModel
}