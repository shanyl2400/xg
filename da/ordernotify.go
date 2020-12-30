package da

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
	"time"
	"xg/db"
)

type OrderNotifies struct {
	ID       int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	OrderID  int    `gorm:"type:int;NOT NULL;column:order_id;index"`
	Classify int    `gorm:"type:int;NOT NULL;column:classify;index"`
	Content  string `gorm:"type:text;NULL;column:content"`
	Author   int    `gorm:"type:int;NULL;column:author"`

	Status int `gorm:"type:int;NOT NULL;column:status;index"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type IOrderNotifiesModel interface {
	CreateOrderNotify(ctx context.Context, tx *gorm.DB, o OrderNotifies) (int, error)
	UpdateOrderNotifyStatus(ctx context.Context, tx *gorm.DB, id int, status int) error
	GetOrderNotifiesByOrderID(ctx context.Context, orderID int) ([]*OrderNotifies, error)
	GetOrderNotifyByID(ctx context.Context, id int) (*OrderNotifies, error)
	SearchOrderNotifies(ctx context.Context, s OrderNotifiesCondition) (int, []*OrderNotifies, error)
	QueryOrderNotifies(ctx context.Context, s OrderNotifiesCondition) ([]*OrderNotifies, error)
}
type DBOrderNotifiesModel struct {
}

func (d *DBOrderNotifiesModel) CreateOrderNotify(ctx context.Context, tx *gorm.DB, o OrderNotifies) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now
	err := db.Get().Create(&o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}
func (d *DBOrderNotifiesModel) UpdateOrderNotifyStatus(ctx context.Context, tx *gorm.DB, id int, status int) error {
	now := time.Now()
	err := tx.Model(OrderNotifies{}).Where(&OrderNotifies{ID: id}).Updates(OrderNotifies{Status: status, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}
func (d *DBOrderNotifiesModel) GetOrderNotifiesByOrderID(ctx context.Context, orderID int) ([]*OrderNotifies, error) {
	return d.QueryOrderNotifies(ctx, OrderNotifiesCondition{
		OrderIDs: []int{orderID},
	})
}
func (d *DBOrderNotifiesModel) GetOrderNotifyByID(ctx context.Context, id int) (*OrderNotifies, error) {
	order := new(OrderNotifies)
	err := db.Get().Where(&OrderNotifies{ID: id}).First(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (d *DBOrderNotifiesModel) QueryOrderNotifies(ctx context.Context, s OrderNotifiesCondition) ([]*OrderNotifies, error) {
	where, values := s.GetConditions()
	//获取学生名单
	records := make([]*OrderNotifies, 0)
	tx := db.Get().Where(where, values...)
	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err := tx.Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
func (d *DBOrderNotifiesModel) SearchOrderNotifies(ctx context.Context, s OrderNotifiesCondition) (int, []*OrderNotifies, error) {
	where, values := s.GetConditions()
	//获取数量
	var total int
	err := db.Get().Model(OrderNotifies{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	records, err := d.QueryOrderNotifies(ctx, s)
	if err != nil {
		return 0, nil, err
	}
	return total, records, nil
}

type OrderNotifiesCondition struct {
	OrderIDs         []int
	Status           []int
	Classifies       []int
	OrderAuthorID    int
	OrderPublisherID int

	PageSize int
	Page     int
	OrderBy  string
}

func (s OrderNotifiesCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.OrderIDs) > 0 {
		wheres = append(wheres, "order_id IN (?)")
		values = append(values, s.OrderIDs)
	}

	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}
	if len(s.Classifies) > 0 {
		wheres = append(wheres, "classify IN (?)")
		values = append(values, s.Classifies)
	}
	if s.OrderAuthorID > 0 {
		orderTable := "orders"
		orderNotifyTable := "order_notifies"
		sql := fmt.Sprintf(`select id from %v where (%v.author_id = ?) and %v.id = %v.order_id and deleted_at IS NULL`,
			orderTable,
			orderTable,
			orderTable,
			orderNotifyTable)
		condition := fmt.Sprintf("exists (%v)", sql)
		wheres = append(wheres, condition)
		values = append(values, s.OrderAuthorID)
	}

	if s.OrderPublisherID > 0 {
		orderTable := "orders"
		orderNotifyTable := "order_notifies"
		sql := fmt.Sprintf(`select id from %v where (%v.publisher_id = ?) and %v.id = %v.order_id and deleted_at IS NULL`,
			orderTable,
			orderTable,
			orderTable,
			orderNotifyTable)
		condition := fmt.Sprintf("exists (%v)", sql)
		wheres = append(wheres, condition)
		values = append(values, s.OrderPublisherID)
	}

	where := strings.Join(wheres, " AND ")

	return where, values
}

var (
	_orderNotifiesModel     *DBOrderNotifiesModel
	_orderNotifiesModelOnce sync.Once
)

func GetOrderNotifiesModel() IOrderNotifiesModel {
	_orderNotifiesModelOnce.Do(func() {
		_orderNotifiesModel = new(DBOrderNotifiesModel)
	})
	return _orderNotifiesModel
}
