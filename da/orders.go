package da

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"xg/db"
	"xg/log"

	"github.com/jinzhu/gorm"
)

type IOrderModel interface {
	CreateOrder(ctx context.Context, tx *gorm.DB, o Order) (int, error)
	AddOrderPayRecord(ctx context.Context, o *OrderPayRecord) (int, error)
	AddRemarkRecord(ctx context.Context, o *OrderRemarkRecord) (int, error)
	AddRemarkRecordTx(ctx context.Context, tx *gorm.DB, o *OrderRemarkRecord) (int, error)
	AddOrderPayRecordTx(ctx context.Context, tx *gorm.DB, o *OrderPayRecord) (int, error)
	UpdateOrderStatusTx(ctx context.Context, db *gorm.DB, id, status int) error
	UpdateOrderPayRecordTx(ctx context.Context, tx *gorm.DB, id, status int) error
	UpdateOrderRemarkRecordTx(ctx context.Context, tx *gorm.DB, ids []int, status int) error

	GetOrderById(ctx context.Context, id int) (*OrderInfo, error)
	SearchOrder(ctx context.Context, s SearchOrderCondition) (int, []*Order, error)
	CountOrder(ctx context.Context, s SearchOrderCondition) (int, error)

	CountPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, error)
	GetPayRecordById(ctx context.Context, id int) (*OrderPayRecord, error)
	SearchPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, []*OrderPayRecord, error)

	SearchRemarkRecord(ctx context.Context, s SearchRemarkRecordCondition) (int, []*OrderRemarkRecord, error)
}

type OrderInfo struct {
	Order       *Order
	PaymentInfo []*OrderPayRecord
	RemarkInfo  []*OrderRemarkRecord
}

type Order struct {
	ID             int    `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	StudentID      int    `gorm:"type:int;NOT NULL;column:student_id;index"`
	ToOrgID        int    `gorm:"type:int;NOT NULL;column:org_id;index"`
	IntentSubjects string `gorm:"type:varchar(255);NOT NULL;column:intent_subjects"`
	PublisherID    int    `gorm:"type:int;NOT NULL;column:publisher_id;index"`

	OrderSource int `gorm:"type:int;NOT NULL;column:order_source;index"`

	Status int `gorm:"type:int;NOT NULL;column:status;index"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
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
	Status  int    `gorm:"type:int;NOT NULL;DEFAULT 1;column:status"`

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
		values = append(values, s.PayRecordIDList)
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
		values = append(values, s.StatusList)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type SearchRemarkRecordCondition struct {
	RemarkRecordIDList []int
	OrderIDList        []int
	AuthorIDList       []int
	Mode               int
	StatusList         []int

	OrderBy string

	PageSize int
	Page     int
}

func (s SearchRemarkRecordCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.RemarkRecordIDList) > 0 {
		wheres = append(wheres, "id IN (?)")
		values = append(values, s.RemarkRecordIDList)
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
		values = append(values, s.StatusList)
	}

	where := strings.Join(wheres, " and ")

	return where, values
}

type SearchOrderCondition struct {
	OrderIDList     []int
	StudentIDList   []int
	ToOrgIDList     []int
	OrderSourceList []int
	IntentSubjects  string
	PublisherID     []int

	StudentKeywords string

	CreateStartAt *time.Time
	CreateEndAt   *time.Time

	Status []int

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
		wheres = append(wheres, "intent_subjects LIKE ?")
		values = append(values, "%"+s.IntentSubjects+"%")
	}
	if len(s.PublisherID) > 0 {
		wheres = append(wheres, "publisher_id IN (?)")
		values = append(values, s.PublisherID)
	}
	if len(s.Status) > 0 {
		wheres = append(wheres, "status IN (?)")
		values = append(values, s.Status)
	}
	if len(s.OrderSourceList) > 0 {
		wheres = append(wheres, "order_source IN (?)")
		values = append(values, s.OrderSourceList)
	}
	//search order by students info
	if s.StudentKeywords != "" {
		studentsTable := "students"
		orderTable := "orders"
		sql := fmt.Sprintf(`select id from %v where (%v.name like ? or %v.telephone like ?) and %v.id = %v.student_id and deleted_at = NULL`,
			studentsTable,
			studentsTable,
			studentsTable,
			studentsTable,
			orderTable)
		condition := fmt.Sprintf("exists (%v)", sql)
		wheres = append(wheres, condition)
		values = append(values, s.StudentKeywords+"%", s.StudentKeywords+"%")
	}

	if s.CreateStartAt != nil && s.CreateEndAt != nil {
		wheres = append(wheres, "created_at BETWEEN ? AND ?")
		values = append(values, s.CreateStartAt, s.CreateEndAt)
	}

	where := strings.Join(wheres, " AND ")

	return where, values
}

type DBOrderModel struct{}

func (d *DBOrderModel) CreateOrder(ctx context.Context, tx *gorm.DB, o Order) (int, error) {
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

func (d *DBOrderModel) AddRemarkRecordTx(ctx context.Context, tx *gorm.DB, o *OrderRemarkRecord) (int, error) {
	now := time.Now()
	o.CreatedAt = &now
	o.UpdatedAt = &now
	err := tx.Create(o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}

func (d *DBOrderModel) UpdateOrderStatusTx(ctx context.Context, tx *gorm.DB, id, status int) error {
	now := time.Now()
	err := tx.Model(Order{}).Where(&Order{ID: id}).Updates(Order{Status: status, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderModel) UpdateOrderPayRecordTx(ctx context.Context, tx *gorm.DB, id, status int) error {
	now := time.Now()
	err := tx.Model(OrderPayRecord{}).Where(&OrderPayRecord{ID: id}).Updates(OrderPayRecord{Status: status, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderModel) UpdateOrderRemarkRecordTx(ctx context.Context, tx *gorm.DB, ids []int, status int) error {
	now := time.Now()
	err := tx.Model(OrderRemarkRecord{}).Where("id in (?)", ids).Updates(OrderRemarkRecord{Status: status, UpdatedAt: &now}).Error
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
		log.Error.Println("Get order remark record failed, err:", err)
	} else {
		orderInfo.RemarkInfo = remarks
	}

	payRecords := make([]*OrderPayRecord, 0)
	err = db.Get().Where(OrderPayRecord{OrderID: id}).Find(&payRecords).Error
	if err != nil {
		log.Error.Println("Get order pay record failed, err:", err)
	} else {
		orderInfo.PaymentInfo = payRecords
	}

	return orderInfo, nil
}

func (d *DBOrderModel) CountPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(Order{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (d *DBOrderModel) GetPayRecordById(ctx context.Context, id int) (*OrderPayRecord, error) {
	record := new(OrderPayRecord)
	err := db.Get().Where(&OrderPayRecord{ID: id}).First(&record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}
func (d *DBOrderModel) SearchPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, []*OrderPayRecord, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(OrderPayRecord{}).Where(where, values...).Count(&total).Error
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

func (d *DBOrderModel) SearchRemarkRecord(ctx context.Context, s SearchRemarkRecordCondition) (int, []*OrderRemarkRecord, error) {
	where, values := s.GetConditions()

	//获取数量
	var total int
	err := db.Get().Model(OrderRemarkRecord{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, nil, err
	}

	//获取学生名单
	records := make([]*OrderRemarkRecord, 0)
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

func (d *DBOrderModel) CountOrder(ctx context.Context, s SearchOrderCondition) (int, error) {
	where, values := s.GetConditions()
	//获取数量
	var total int
	err := db.Get().Model(Order{}).Where(where, values...).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
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
