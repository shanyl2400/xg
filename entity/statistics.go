package entity

const (
	StudentStatisticsKey     = "students"
	PerformanceStatisticsKey = "performance"

	OrgPerformanceStatisticsKey = "org"
	AuthorPerformanceStatisticsKey = "author"
	PublisherPerformanceStatisticsKey = "publisher"
)

type SummaryInfo struct {
	OrgsTotal        int `json:"orgs_total"`
	StudentsTotal    int `json:"students_total"`
	PerformanceTotal int `json:"performance_total"`
	SuccessRate      int `json:"success_rate"`
}
type StatisticRecord struct {
	Key   string `json:"key"`
	Value int    `json:"value"`

	Year  int `json:"year"`
	Month int `json:"month"`
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
}