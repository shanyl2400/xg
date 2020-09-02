package entity

const(
	StudentStatisticsKey = "students"
	PerformanceStatisticsKey = "performance"
)

type SummaryInfo struct {
	OrgsTotal int `json:"orgs_total"`
	StudentsTotal int `json:"students_total"`
	PerformanceTotal int `json:"performance_total"`
	SuccessRate int `json:"success_rate"`
}
type StatisticRecord struct {
	Key     string `json:"key"`
	Value int `json:"value"`

	Year int `json:"year"`
	Month int `json:"month"`
}

type StatisticGraph struct {
	StudentsGraph     []*StatisticRecord `json:"students_graph"`
	PerformancesGraph []*StatisticRecord `json:"performances_graph"`
}