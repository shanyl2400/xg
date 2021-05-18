package da

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

type IOrderModel interface {
	CreateOrder(ctx context.Context, tx *gorm.DB, o Order) (int, error)
	AddOrderPayRecord(ctx context.Context, o *OrderPayRecord) (int, error)
	AddRemarkRecord(ctx context.Context, o *OrderRemarkRecord) (int, error)
	AddRemarkRecordTx(ctx context.Context, tx *gorm.DB, o *OrderRemarkRecord) (int, error)
	AddOrderPayRecordTx(ctx context.Context, tx *gorm.DB, o *OrderPayRecord) (int, error)
	UpdateOrderStatusTx(ctx context.Context, db *gorm.DB, id int, status int) error
	UpdateOrderUpdateAtTx(ctx context.Context, tx *gorm.DB, id int) error
	UpdateOrderPayRecordTx(ctx context.Context, tx *gorm.DB, id, status int) error
	UpdateOrderPayRecordForParentTx(ctx context.Context, tx *gorm.DB, id, status, parentID int) error
	UpdateOrderRemarkRecordTx(ctx context.Context, tx *gorm.DB, ids []int, status int) error

	ReplaceOrderSource(ctx context.Context, tx *gorm.DB, oldOrderSource, newOrderSource int) error
	UpdateOrderPayRecordPriceTx(ctx context.Context, tx *gorm.DB, id int, price float64) error

	GetOrderById(ctx context.Context, id int) (*OrderInfo, error)
	SearchOrder(ctx context.Context, s SearchOrderCondition) (int, []*Order, error)
	CountOrder(ctx context.Context, s SearchOrderCondition) (int, error)

	CountPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, error)
	GetPayRecordById(ctx context.Context, id int) (*OrderPayRecord, error)
	SearchPayRecord(ctx context.Context, s SearchPayRecordCondition) (int, []*OrderPayRecord, error)

	SearchRemarkRecord(ctx context.Context, s SearchRemarkRecordCondition) (int, []*OrderRemarkRecord, error)

	StatisticOrdersPayments(ctx context.Context, groupby string, limit int, s SearchOrderCondition, s0 SearchPayRecordCondition) ([]*entity.GroupbyStatisticEntity, error)
	StatisticOrders(ctx context.Context, groupby string, limit int, s SearchOrderCondition) ([]*entity.GroupbyStatisticEntity, error)
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

	ParentOrgID int `gorm:"type:int;NOT NULL;column:parent_org_id;index"`

	Address  string `gorm:"type:varchar(255);NULL;column:address"`
	AuthorID int    `gorm:"type:int;NOT NULL;DEFAULT 1;column:author_id;index"`

	OrderSource int `gorm:"type:int;NOT NULL;column:order_source;index"`

	Status int `gorm:"type:int;NOT NULL;column:status;index"`

	UpdatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:updated_at"`
	CreatedAt *time.Time `gorm:"type:datetime;NOT NULL;column:created_at"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at"`
}

type OrderPayRecord struct {
	ID      int     `gorm:"PRIMARY_KEY;AUTO_INCREMENT;column:id"`
	OrderID int     `gorm:"type:int;NOT NULL;column:order_id"`
	Mode    int     `gorm:"type:int;NOT NULL;column:mode"`
	Title   string  `gorm:"type:varchar(128);NOT NULL;column:title"`
	Amount  float64 `gorm:"type:DECIMAL(11,2);NOT NULL;column:amount"`
	Content string  `gorm:"type:text;NOT NULL;column:content"`

	RealPrice float64 `gorm:"type:DECIMAL(11,2);NOT NULL;column:real_price"`

	Status   int `gorm:"type:int;NOT NULL;column:status"`
	ParentID int `gorm:"type:int;NOT NULL;column:parent_id"`

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

	InfoType  int        `gorm:"type:int;NOT NULL;column:info_type"`
	Info      string     `gorm:"type:text;NULL;column:info"`
	RevisitAt *time.Time `gorm:"type:datetime;NULL;column:revisit_at"`

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
	CreateStartAt   *time.Time
	CreateEndAt     *time.Time

	OrderBy string

	PageSize int
	Page     int

	Prefix string
}

func (s SearchPayRecordCondition) prefix(column string) string {
	if s.Prefix == "" {
		return "" + column
	}
	return s.Prefix + "." + column
}

func (s SearchPayRecordCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.PayRecordIDList) > 0 {
		wheres = append(wheres, s.prefix("id")+" IN (?)")
		values = append(values, s.PayRecordIDList)
	}
	if len(s.AuthorIDList) > 0 {
		wheres = append(wheres, s.prefix("author")+" IN (?)")
		values = append(values, s.AuthorIDList)
	}
	if len(s.OrderIDList) > 0 {
		wheres = append(wheres, s.prefix("order_id")+" IN (?)")
		values = append(values, s.OrderIDList)
	}
	if s.CreateStartAt != nil && s.CreateEndAt != nil {
		wheres = append(wheres, s.prefix("created_at")+" BETWEEN ? AND ?")
		values = append(values, s.CreateStartAt, s.CreateEndAt)
	}

	if s.Mode > 0 {
		wheres = append(wheres, s.prefix("mode")+" = ?")
		values = append(values, s.Mode)
	}
	if len(s.StatusList) > 0 {
		wheres = append(wheres, s.prefix("status")+" IN (?)")
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
	PublisherIDList []int
	AuthorIDList    []int

	RelatedUserIDList []int

	StudentKeywords string

	Keywords string
	Address  string

	CreateStartAt *time.Time
	CreateEndAt   *time.Time
	UpdateStartAt *time.Time
	UpdateEndAt   *time.Time

	Status []int

	OrderBy string

	PageSize int
	Page     int

	Prefix string
}

func (s SearchOrderCondition) prefix(column string) string {
	if s.Prefix == "" {
		return "" + column
	}
	return s.Prefix + "." + column
}

func (s SearchOrderCondition) GetConditions() (string, []interface{}) {
	wheres := make([]string, 0)
	values := make([]interface{}, 0)

	if len(s.OrderIDList) > 0 {
		wheres = append(wheres, s.prefix("id")+" IN (?)")
		values = append(values, s.OrderIDList)
	}
	if len(s.StudentIDList) > 0 {
		wheres = append(wheres, s.prefix("student_id")+" IN (?)")
		values = append(values, s.StudentIDList)
	}
	if len(s.ToOrgIDList) > 0 {
		wheres = append(wheres, s.prefix("org_id")+" IN (?)")
		values = append(values, s.ToOrgIDList)
	}
	if s.IntentSubjects != "" {
		wheres = append(wheres, s.prefix("intent_subjects")+" LIKE ?")
		values = append(values, "%"+s.IntentSubjects+"%")
	}
	if len(s.RelatedUserIDList) > 0 {
		wheres = append(wheres, fmt.Sprintf("(%v IN (?) or %v IN (?))", s.prefix("publisher_id"), s.prefix("author_id")))
		values = append(values, s.RelatedUserIDList, s.RelatedUserIDList)
	}
	if len(s.PublisherIDList) > 0 {
		wheres = append(wheres, s.prefix("publisher_id")+" IN (?)")
		values = append(values, s.PublisherIDList)
	}
	if len(s.AuthorIDList) > 0 {
		wheres = append(wheres, s.prefix("author_id")+" IN (?)")
		values = append(values, s.AuthorIDList)
	}
	if s.Address != "" {
		wheres = append(wheres, s.prefix("address")+" LIKE ?")
		values = append(values, s.Address+"%")
	}
	if len(s.Status) > 0 {
		wheres = append(wheres, s.prefix("status")+" IN (?)")
		values = append(values, s.Status)
	}
	if len(s.OrderSourceList) > 0 {
		wheres = append(wheres, s.prefix("order_source")+" IN (?)")
		values = append(values, s.OrderSourceList)
	}

	if s.Keywords != "" {
		condition1 := "(intent_subjects LIKE ?)"

		studentsTable := "students"
		orderTable := "orders"
		sql := fmt.Sprintf(`select id from %v where (%v.name like ? or %v.telephone like ?) and %v.id = %v.student_id and deleted_at IS NULL`,
			studentsTable,
			studentsTable,
			studentsTable,
			studentsTable,
			orderTable)
		condition2 := fmt.Sprintf("(exists (%v))", sql)

		condition := "(" + condition1 + " or " + condition2 + ")"
		wheres = append(wheres, condition)

		values = append(values, "%"+s.Keywords+"%")
		values = append(values, s.Keywords+"%", s.Keywords+"%")
	}
	//search order by students info
	if s.StudentKeywords != "" {
		studentsTable := "students"
		orderTable := "orders"
		sql := fmt.Sprintf(`select id from %v where (%v.name like ? or %v.telephone like ?) and %v.id = %v.student_id and deleted_at IS NULL`,
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
		wheres = append(wheres, s.prefix("created_at")+" BETWEEN ? AND ?")
		values = append(values, s.CreateStartAt, s.CreateEndAt)
	}

	if s.UpdateStartAt != nil && s.UpdateEndAt != nil {
		wheres = append(wheres, s.prefix("updated_at")+" BETWEEN ? AND ?")
		values = append(values, s.UpdateStartAt, s.UpdateEndAt)
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

	if o.Mode == entity.OrderPayModePayback {
		o.RealPrice = -o.Amount
	} else {
		o.RealPrice = o.Amount
	}
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
	if o.Mode == entity.OrderPayModePayback {
		o.RealPrice = -o.Amount
	} else {
		o.RealPrice = o.Amount
	}
	err := tx.Create(o).Error
	if err != nil {
		return -1, err
	}
	return o.ID, nil
}

func (d *DBOrderModel) ReplaceOrderSource(ctx context.Context, tx *gorm.DB, oldOrderSource, newOrderSource int) error {
	err := tx.Model(&Order{}).Where(" order_source = ?", oldOrderSource).Updates(Order{OrderSource: newOrderSource}).Error
	if err != nil {
		return err
	}
	return nil
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

func (d *DBOrderModel) UpdateOrderUpdateAtTx(ctx context.Context, tx *gorm.DB, id int) error {
	now := time.Now()
	err := tx.Model(Order{}).Where(&Order{ID: id}).Updates(Order{UpdatedAt: &now}).Error
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

func (d *DBOrderModel) UpdateOrderPayRecordForParentTx(ctx context.Context, tx *gorm.DB, id, status, parentID int) error {
	now := time.Now()
	err := tx.Model(OrderPayRecord{}).Where(&OrderPayRecord{ID: id}).Updates(OrderPayRecord{Status: status, ParentID: parentID, UpdatedAt: &now}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DBOrderModel) UpdateOrderPayRecordPriceTx(ctx context.Context, tx *gorm.DB, id int, price float64) error {
	now := time.Now()

	err := tx.Model(OrderPayRecord{}).Where(&OrderPayRecord{ID: id}).Updates(OrderPayRecord{Amount: price, UpdatedAt: &now}).Error
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

func (d *DBOrderModel) StatisticOrdersWithStatus(ctx context.Context, groupby string, limit int, s SearchOrderCondition) ([]*entity.GroupbyStatisticEntity, error) {
	where, values := s.GetConditions()
	tx := db.Get().Table("orders").Select(fmt.Sprintf("%v as id, status, count(*) as cnt", groupby)).Where(where, values...).Group(groupby + ",status")
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	tx = tx.Order("cnt desc")
	entities := make([]*entity.GroupbyStatisticEntity, 0)
	err := tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *DBOrderModel) StatisticOrders(ctx context.Context, groupby string, limit int, s SearchOrderCondition) ([]*entity.GroupbyStatisticEntity, error) {
	where, values := s.GetConditions()
	tx := db.Get().Table("orders").Select(fmt.Sprintf("%v as id, status, count(*) as cnt", groupby)).Where(where, values...).Group(groupby + ",status")
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	tx = tx.Order("cnt desc")
	entities := make([]*entity.GroupbyStatisticEntity, 0)
	err := tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (d *DBOrderModel) StatisticOrdersPayments(ctx context.Context, groupby string, limit int, s SearchOrderCondition, s0 SearchPayRecordCondition) ([]*entity.GroupbyStatisticEntity, error) {
	//查询所有相关orders
	s.Prefix = "o"
	s0.Prefix = "p"
	where1, values1 := s.GetConditions()
	where2, values2 := s0.GetConditions()
	rawSQL := fmt.Sprintf(`SELECT o.%v as id, o.status as status, count(*) as cnt, sum(real_price) as amount
	FROM orders as o left join order_pay_records as p on o.id=p.order_id `, groupby)

	wheres := make([]string, 0)
	if where1 != "" {
		wheres = append(wheres, "("+where1+")")
	}
	if where2 != "" {
		wheres = append(wheres, "("+where2+")")
	}
	where := strings.Join(wheres, "AND")

	if where != "" {
		rawSQL = rawSQL + " WHERE " + where
	}

	rawSQL = rawSQL + fmt.Sprintf(" group by o.%v, o.status", groupby)
	values := append(values1, values2...)
	tx := db.Get().Raw(rawSQL, values...)

	fmt.Println(rawSQL)
	fmt.Println(values1)
	// if s.OrderBy != "" {
	// 	tx = tx.Order("o." + s.OrderBy)
	// }
	if limit > 0 {
		tx = tx.Limit(limit)
	}
	entities := make([]*entity.GroupbyStatisticEntity, 0)
	err := tx.Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
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
