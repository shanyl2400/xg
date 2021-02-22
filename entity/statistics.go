package entity

import "time"

const (
	StudentStatisticsKey     = "students"
	PerformanceStatisticsKey = "performance"

	StudentAuthorStatisticsKey        = "stu_author"
	OrgPerformanceStatisticsKey       = "org"
	AuthorPerformanceStatisticsKey    = "author"
	PublisherPerformanceStatisticsKey = "publisher"

	OrderSourcePerformanceStatisticsKey = "order_source"
)

const (
	OrderStatisticKeyStudent       = "student"
	OrderStatisticKeyOrder         = "order"
	OrderStatisticKeyInvalidOrder  = "invalid_order"
	OrderStatisticKeyConsiderOrder = "consider_order"
	OrderStatisticKeyNewOrder      = "new_order"
	OrderStatisticKeySignupOrder   = "signup_order"
)

type OrderStatisticRecordId struct {
	Key         string `json:"key"`
	Author      int    `json:"author"`
	OrgId       int    `json:"org_id"`
	PublisherId int    `json:"publisher_id"`
	OrderSource int    `json:"order_source"`
}

type OrderStatisticRecordEntity struct {
	Author      int `json:"author"`
	OrgId       int `json:"org_id"`
	PublisherId int `json:"publisher_id"`
	OrderSource int `json:"order_source"`
}

type StatisticRecordCondition struct {
	OrderStatisticRecordEntity
	StartAt *time.Time `json:"start_at"`
	EndAt   *time.Time `json:"end_at"`

	GroupBy string `json:"group_by"`
}

type SummaryInfo struct {
	OrgsTotal        int     `json:"orgs_total"`
	StudentsTotal    int     `json:"students_total"`
	PerformanceTotal float64 `json:"performance_total"`
	SuccessRate      int     `json:"success_rate"`
}
type StatisticRecord struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
	Count int     `json:"count"`

	Year  int `json:"year"`
	Month int `json:"month"`
}

type TotalStatisticRecord struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
	Count int     `json:"count"`
}

type StatisticGraph struct {
	StudentsGraph     []*StatisticRecord `json:"students_graph"`
	PerformancesGraph []*StatisticRecord `json:"performances_graph"`
}

type PerformancesGraph struct {
	PerformancesGraph []*StatisticRecord `json:"performances_graph"`
}

type AuthorPerformancesGraph struct {
	AuthorPerformancesGraph    []*StatisticRecord `json:"author_performances_graph"`
	PublisherPerformancesGraph []*StatisticRecord `json:"publisher_performances_graph"`
}

type OrderPerformanceInfo struct {
	OrgId         int `json:"org_id"`
	AuthorId      int `json:"author_id"`
	PublisherId   int `json:"publisher_id"`
	OrderSourceId int `json:"order_source_id"`
}
type OrderStatisticDate struct {
	Year  int
	Month int
	Date  int
}

//名单数，无效人数，报名人数，成交业绩，成功率
type OrderStatisticTable struct {
	Data              []*OrderStatisticTableMonth `json:"data"`
	DayData           OrderStatisticTableItem     `json:"day_data"`
	WeekDayData       OrderStatisticTableItem     `json:"week_day_data"`
	MonthDayData      OrderStatisticTableItem     `json:"month_day_data"`
	ThreeMonthDayData OrderStatisticTableItem     `json:"three_month_day_data"`
}

func (o *OrderStatisticTable) CalculateSucceed() {
	o.DayData.CalculateSucceed()
	o.WeekDayData.CalculateSucceed()
	o.MonthDayData.CalculateSucceed()
	o.ThreeMonthDayData.CalculateSucceed()
}

func NewOrderStatisticTable() *OrderStatisticTable {
	tb := new(OrderStatisticTable)

	//初始化12个月
	for i := 0; i < 12; i++ {
		tb.Data = append(tb.Data, new(OrderStatisticTableMonth))
	}
	return tb
}

type OrderStatisticTableMonth struct {
	Students      int     `json:"students"`
	Orders        int     `json:"orders"`
	InvalidOrders int     `json:"invalid_orders"`
	SignedOrder   int     `json:"signed_order"`
	Performance   float64 `json:"performance"`
}

type OrderStatisticTableItem struct {
	OrderStatisticTableMonth
	Succeed int `json:"succeed"`
}

type OrderStatisticGroupTableItem struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	OrderStatisticTableItem
}

func (o *OrderStatisticTableItem) CalculateSucceed() {
	if o.Orders == 0 {
		o.Succeed = 0
		return
	}
	o.Succeed = (o.SignedOrder * 10000) / o.Orders
}
