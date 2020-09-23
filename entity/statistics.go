package entity

const (
	StudentStatisticsKey     = "students"
	PerformanceStatisticsKey = "performance"

	StudentAuthorStatisticsKey     = "stu_author"
	OrgPerformanceStatisticsKey = "org"
	AuthorPerformanceStatisticsKey = "author"
	PublisherPerformanceStatisticsKey = "publisher"

	OrderSourcePerformanceStatisticsKey = "order_source"
)

const(
	OrderStatisticKeyStudent      = "student"
	OrderStatisticKeyOrder        = "order"
	OrderStatisticKeyInvalidOrder = "invalid_order"
	OrderStatisticKeyNewOrder     = "new_order"
	OrderStatisticKeySignupOrder  = "signup_order"
)


type OrderStatisticRecordId struct {
	Key string `json:"key"`
	Author int `json:"author"`
	OrgId int `json:"org_id"`
	PublisherId int `json:"publisher_id"`
	OrderSource int `json:"order_source"`
}

type OrderStatisticRecordEntity struct {
	Author int `json:"author"`
	OrgId int `json:"org_id"`
	PublisherId int `json:"publisher_id"`
	OrderSource int `json:"order_source"`
}

type SummaryInfo struct {
	OrgsTotal        int `json:"orgs_total"`
	StudentsTotal    int `json:"students_total"`
	PerformanceTotal int `json:"performance_total"`
	SuccessRate      int `json:"success_rate"`
}
type StatisticRecord struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Count int `json:"count"`

	Year  int `json:"year"`
	Month int `json:"month"`
}

type TotalStatisticRecord struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Count int `json:"count"`
}

type StatisticGraph struct {
	StudentsGraph     []*StatisticRecord `json:"students_graph"`
	PerformancesGraph []*StatisticRecord `json:"performances_graph"`
}

type PerformancesGraph struct {
	PerformancesGraph []*StatisticRecord `json:"performances_graph"`
}

type AuthorPerformancesGraph struct {
	AuthorPerformancesGraph []*StatisticRecord `json:"author_performances_graph"`
	PublisherPerformancesGraph []*StatisticRecord `json:"publisher_performances_graph"`
}

type OrderPerformanceInfo struct {
	OrgId int `json:"org_id"`
	AuthorId int `json:"author_id"`
	PublisherId int `json:"publisher_id"`
	OrderSourceId int `json:"order_source_id"`
}