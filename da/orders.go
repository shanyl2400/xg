package da

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"xg/db"

	"github.com/jinzhu/gorm"
)

type IOrderModel interface {
	CreateOrder(ctx context.Context, o Order) (int, error)
	AddOrderPayRecord(ctx context.Context, o *OrderPayRecord) (int, error)
	AddRemarkRecord(ctx context.Context, o *OrderRemarkRecord) (int, error)
	AddOrderPayRecordTx(ctx context.Context, tx *gorm.DB, o *OrderPayRecord) (int, error)
	UpdateOrderStatusTx(ctx context.Context, db *gorm.DB, id, status int) error
	UpdateOrderPayRecordTx(ctx context.Context, tx *gorm.DB, id, status int) error

	GetOrderById(ctx context.Context, id int) (*OrderInfo, error)
	SearchOrder(ctx context.Context, s SearchOrderCondition) (int, []*Order, error)

	SearchPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, []*OrderPayRecord, error)
}

type OrderInfo struct {
	Order       *Order
	PaymentInfo []*OrderPayRecord
	RemarkInfo  []*OrderRemarkRecord
}

type Order struct {
	ID             int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	StudentID      int    `gorm:"type:int;NOT NULL;column:student_id"`
	ToOrgID        int    `gorm:"type:int;NOT NULL;column:org_id"`
	IntentSubjects string `gorm:"type:varchar(255);NOT NULL;column:intent_subjects"`
	PublisherID    int    `gorm:"type:int;NOT NULL;column:publisher_id"`

	Status int `gorm:"type:int;NOT NULL;column:status"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SearchOrdersCondition struct {
	IDList          []int
	StudentIDList   []int
	ToOrgIDList     []int
	PublisherIDList []int
}

type OrderPayRecord struct {
	ID      int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	OrderID int    `gorm:"type:int;NOT NULL;column:order_id"`
	Mode    int    `gorm:"type:int;NOT NULL;column:mode"`
	Title   string `gorm:"type:varchar(128);NOT NULL;column:title"`
	Amount  int    `gorm:"type:int;NOT NULL;column:amount"`

	Status int `gorm:"type:int;NOT NULL;column:status"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type OrderRemarkRecord struct {
	ID      int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	OrderID int    `gorm:"type:int;NOT NULL;column:order_id"`
	Author  int    `gorm:"type:int;NOT NULL;column:author"`
	Mode    int    `gorm:"type:int;NOT NULL;column:mode"`
	Content string `gorm:"type:text;NOT NULL;column:content"`
	Status  int    `gorm:"type:int;NOT NULL;column:status"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type SearchPayRecordCondition struct {
	PayRecordIDList []int
	OrderIDList     []int
	AuthorIDList    []int
	Mode            int
	StatusList      []int

	OrderBy string

	PageSize int
	Page     int
}

func (s SearchPayRecordCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.PayRecordIDList) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.OrderIDList)
	}
	if len(s.AuthorIDList) > 0 {
		wheres = append(wheres, "author IN (?)")
		values = append(values, s.AuthorIDList)
	}
	if len(s.OrderIDList) > 0 {
		wheres = append(wheres, "order_id IN (?)")
		values = append(values, s.OrderIDList)
	}
	if s.Mode > 0 {
		wheres = append(wheres, "mode = ?")
		values = append(values, s.Mode)
	}
	if len(s.StatusList) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Mode)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type SearchOrderCondition struct {
	OrderIDList    []int
	StudentIDList  []int
	ToOrgIDList    []int
	IntentSubjects string
	PublisherID    int

	Status int

	OrderBy string

	PageSize int
	Page     int
}

func (s SearchOrderCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.OrderIDList) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.OrderIDList)
	}
	if len(s.StudentIDList) > 0 {
		wheres = append(wheres, "student_id IN (?)")
		values = append(values, s.StudentIDList)
	}
	if len(s.ToOrgIDList) > 0 {
		wheres = append(wheres, "org_id IN (?)")
		values = append(values, s.ToOrgIDList)
	}
	if s.IntentSubjects != "" {
		wheres = append(wheres, "JSON_CONTAINS(intent_subjects, ?)")
		values = append(values, s.IntentSubjects)
	}
	if s.PublisherID > 0 {
		wheres = append(wheres, "publisher_id = ?")
		values = append(values, s.PublisherID)
	}
	if s.Status > 0 {
		wheres = append(wheres, "status = ?")
		values = append(values, s.Status)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type DBOrderModel struct{}

func (d *DBOrderModel) CreateOrder(ctx context.Context, o Order) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now
	err := db.Get().Create(&o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}

func (d *DBOrderModel) AddOrderPayRecord(ctx context.Context, o *OrderPayRecord) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now

	err := db.Get().Create(o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}
func (d *DBOrderModel) AddOrderPayRecordTx(ctx context.Context, tx *gorm.DB, o *OrderPayRecord) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now

	err := tx.Create(o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}

func (d *DBOrderModel) AddRemarkRecord(ctx context.Context, o *OrderRemarkRecord) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now
	err := db.Get().Create(o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}

func (d *DBOrderModel) UpdateOrderStatusTx(ctx context.Context, tx *gorm.DB, id, status int) error {
	now := time.Now()
	err := tx.Where(&Order{ID: id}).Updates(Order{Status: status, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderModel) UpdateOrderPayRecordTx(ctx context.Context, tx *gorm.DB, id, status int) error {
	now := time.Now()
	err := tx.Where(&OrderPayRecord{ID: id}).Updates(OrderPayRecord{Status: status, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderModel) GetOrderById(ctx context.Context, id int) (*OrderInfo, error) {
	orderInfo := new(OrderInfo)
	order := new(Order)
	err := db.Get().Where(&Order{ID: id}).First(&order).Error
	if err != nil {
		return nil, err
	}
	orderInfo.Order = order

	remarks := make([]*OrderRemarkRecord, 0)
	err = db.Get().Where(OrderRemarkRecord{OrderID: id}).Find(&remarks).Error
	if err != nil {
		fmt.Println("Get order remark record failed, err:", err)
	} else {
		orderInfo.RemarkInfo = remarks
	}

	payRecords := make([]*OrderPayRecord, 0)
	err = db.Get().Where(OrderRemarkRecord{OrderID: id}).Find(&payRecords).Error
	if err != nil {
		fmt.Println("Get order pay record failed, err:", err)
	} else {
		orderInfo.PaymentInfo = payRecords
	}

	return orderInfo, nil
}

func (d *DBOrderModel) SearchPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, []*OrderPayRecord, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(Order{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	records := make([]*OrderPayRecord, 0)
	tx := db.Get().Where(where, values...)
	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err = tx.Find(&records).Error
	if err != nil {
		return 0, nil, err
	}
	return total, records, nil
}

func (d *DBOrderModel) SearchOrder(ctx context.Context, s SearchOrderCondition) (int, []*Order, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(Order{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	orders := make([]*Order, 0)
	tx := db.Get().Where(where, values...)
	if s.PageSize > 0 {
		offset, limit := parsePage(s.Page, s.PageSize)
		tx = tx.Offset(offset).Limit(limit)
	}
	if s.OrderBy != "" {
		tx = tx.Order(s.OrderBy)
	}
	err = tx.Find(&orders).Error
	if err != nil {
		return 0, nil, err
	}
	return total, orders, nil
}

var (
	_orderModel     *DBOrderModel
	_orderModelOnce sync.Once
)

func GetOrderModel() IOrderModel {
	_orderModelOnce.Do(func() {
		_orderModel = new(DBOrderModel)
	})
	return _orderModel
}
