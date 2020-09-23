package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
)

type OrderStatisticsService struct {
}

//名单数，无效人数，报名人数，成交业绩，成功率
func (s *OrderStatisticsService) StatisticsTable(ctx context.Context, durationDay int) {

}

func (s *OrderStatisticsService) Summary(ctx context.Context) (*entity.SummaryInfo, error) {
	log.Info.Printf("Get summary\n")
	orgsCount, err := da.GetOrgModel().CountOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{0},
	})
	if err != nil {
		log.Warning.Printf("CountOrgs failed, err: %v\n", err)
		return nil, err
	}
	studentsCount, err := da.GetStudentModel().CountStudents(ctx)
	if err != nil {
		log.Warning.Printf("CountStudents failed, err: %v\n", err)
		return nil, err
	}
	payCondition :=  da.SearchPayRecordCondition{
		// Mode:       entity.OrderPayModePay,
		StatusList: []int{entity.OrderPayStatusChecked},
		PageSize:   1000000,
	}
	_, payRecords, err := da.GetOrderModel().SearchPayRecord(ctx, payCondition)
	if err != nil {
		log.Warning.Printf("Search pay records failed, condition: %#v, err: %v\n", payCondition, err)
		return nil, err
	}
	performanceTotal := 0
	for i := range payRecords {
		if payRecords[i].Mode == entity.OrderPayModePay {
			performanceTotal = performanceTotal + payRecords[i].Amount
		} else {
			performanceTotal = performanceTotal - payRecords[i].Amount
		}
	}

	orderCondition := da.SearchOrderCondition{
		Status: []int{entity.OrderStatusSigned,
			entity.OrderStatusDeposit,
			entity.OrderStatusCreated,
			entity.OrderStatusRevoked,
			entity.OrderStatusInvalid},
		Page:   1000000,
	}
	total, orders, err := da.GetOrderModel().SearchOrder(ctx, orderCondition)
	if err != nil {
		log.Warning.Printf("Search orders failed, condition: %#v, err: %v\n", orderCondition, err)
		return nil, err
	}
	successTotal := 0
	for i := range orders {
		if orders[i].Status == entity.OrderStatusSigned ||
			orders[i].Status == entity.OrderStatusDeposit{
			successTotal = successTotal + 1
		}
	}
	successRate := 0
	if total > 0 {
		successRate = successTotal * 10000 / total
	}
	return &entity.SummaryInfo{
		OrgsTotal:        orgsCount,
		StudentsTotal:    studentsCount,
		PerformanceTotal: performanceTotal,
		SuccessRate:      successRate,
	}, nil
}

func (s *OrderStatisticsService) SearchRecords(ctx context.Context, condition da.SearchOrderStatisticsRecordCondition) ([]*da.OrderStatisticsRecord, error) {
	log.Info.Printf("SearchYearRecords, condition: %#v\n", condition)
	records, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
	if err != nil {
		log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	return records, nil
}

func (s *OrderStatisticsService) SearchRecordsTotal(ctx context.Context, condition da.SearchOrderStatisticsRecordCondition) (*entity.TotalStatisticRecord, error) {
	log.Info.Printf("SearchYearRecords, condition: %#v\n", condition)
	records, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
	if err != nil {
		log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	ret := new(entity.TotalStatisticRecord)
	for i := range records {
		ret.Key = records[i].Key
		ret.Count = ret.Count + records[i].Count
		ret.Value = ret.Value + records[i].Value
	}
	return ret, nil
}


func (s *OrderStatisticsService) SearchRecordsMonth(ctx context.Context, condition da.SearchOrderStatisticsRecordCondition) ([]*entity.StatisticRecord, error) {
	log.Info.Printf("SearchYearRecords, condition: %#v\n", condition)
	if condition.Key == "" {
		log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v\n", condition)
		return nil, ErrInvalidStatisticKey
	}
	records, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
	if err != nil {
		log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	monthRecord := make(map[int][]*da.OrderStatisticsRecord)
	for i := range records {
		monthRecord[records[i].Month] = append(monthRecord[records[i].Month], records[i])
	}

	ret := make([]*entity.StatisticRecord, 12)
	year := time.Now().Year()
	for i := 1; i <= 12; i++ {
		_, ok := monthRecord[i]
		if ok {
			value := 0
			count := 0
			for j := range monthRecord[i] {
				value = value + monthRecord[i][j].Value
				count = count + monthRecord[i][j].Count
			}
			ret[i-1] = &entity.StatisticRecord{
				Key:   condition.Key,
				Year:  year,
				Month: i,
				Value: value,
				Count: count,
			}
		} else {
			ret[i-1] = &entity.StatisticRecord{
				Key:   condition.Key,
				Year:  year,
				Month: i,
				Value: 0,
				Count: 0,
			}
		}
	}
	return ret, nil
}

func (s *OrderStatisticsService) AddStudent(ctx context.Context, tx *gorm.DB, authorId, orderSourceId int) error {
	log.Info.Printf("AddStudent, authorId: %#v, orderSourceId: %#v\n", authorId, orderSourceId)
	return s.addValue(ctx, tx, entity.OrderStatisticRecordId{
		Key:         entity.OrderStatisticKeyStudent,
		Author:      authorId,
		OrderSource: orderSourceId,
	}, 1, true)
}

func (s *OrderStatisticsService) AddNewOrder(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity) error {
	return s.addOrder(ctx, tx, osr, entity.OrderStatisticKeyNewOrder, 1)
}

func (s *OrderStatisticsService) AddSignupOrder(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity) error {
	return s.addOrder(ctx, tx, osr, entity.OrderStatisticKeySignupOrder, 1)
}

func (s *OrderStatisticsService) AddInvalidOrder(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity) error {
	return s.addOrder(ctx, tx, osr, entity.OrderStatisticKeyInvalidOrder, 1)
}

func (s *OrderStatisticsService) AddPerformance(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity, performance int) error {
	log.Info.Printf("AddPerformance, value: %#v, performance: %#v\n", osr, performance)
	addCount := false
	//大于0表示成交，计算成交量
	if performance > 0 {
		addCount = true
	}

	err := s.addValue(ctx, tx, entity.OrderStatisticRecordId{
		Key:         entity.OrderStatisticKeyOrder,
		Author:      osr.Author,
		OrgId:       osr.OrgId,
		PublisherId: osr.PublisherId,
		OrderSource: osr.OrderSource,
	}, performance, addCount)
	if err != nil{
		return err
	}
	return nil
}

func (s *OrderStatisticsService) addValue(ctx context.Context, tx *gorm.DB, id entity.OrderStatisticRecordId, value int, addCount bool) error {
	condition := s.idToCondition(id)

	records, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, tx, condition)
	if err != nil {
		log.Warning.Printf("add graph value failed, condition: %#v, err: %v\n", condition, err)
		return err
	}
	if len(records) > 0 {
		record := records[0]
		record.Value = record.Value + value
		if addCount {
			record.Count = record.Count + 1
		}
		err = da.GetOrderStatisticsRecordModel().UpdateOrderStatisticsRecord(ctx, tx, record.ID, record.Value, record.Count)
		if err != nil {
			log.Warning.Printf("UpdateStatisticsRecord failed, record: %#v, err: %v\n", record, err)
			return err
		}
		return nil
	}

	count := 0
	if addCount {
		count = 1
	}
	now := time.Now()
	record := &da.OrderStatisticsRecord{
		Key:    id.Key,
		Value:  value,
		Year:   now.Year(),
		Month:  int(now.Month()),
		Date: now.Day(),
		Count: count,
		Author: id.Author,
		OrgId: id.OrgId,
		OrderSource: id.OrderSource,
		PublisherId: id.PublisherId,
	}
	_, err = da.GetOrderStatisticsRecordModel().CreateOrderStatisticsRecord(ctx, tx, record)
	if err != nil {
		log.Warning.Printf("CreateStatisticsRecord failed, record: %#v, err: %v\n", record, err)
		return err
	}
	return nil
}

func (s *OrderStatisticsService) addOrder(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity, key string, count int) error {
	log.Info.Printf("AddNewOrder, count: %#v, key: %#v\n", count, key)
	return s.addValue(ctx, tx, entity.OrderStatisticRecordId{
		Key:         key,
		Author:      osr.Author,
		OrgId:       osr.OrgId,
		PublisherId: osr.PublisherId,
		OrderSource: osr.OrderSource,
	}, count, true)
}
func (s *OrderStatisticsService) idToCondition(id entity.OrderStatisticRecordId) da.SearchOrderStatisticsRecordCondition{
	now := time.Now()
	condition := da.SearchOrderStatisticsRecordCondition{
		Key:    id.Key,
		Year:   now.Year(),
		Month:  []int{int(now.Month())},
		Date:  []int{now.Day()},
	}
	if id.Author > 0 {
		condition.Author = []int{id.Author}
	}
	if id.OrgId > 0 {
		condition.OrgId = []int{id.OrgId}
	}
	if id.OrderSource > 0 {
		condition.OrderSource = []int{id.OrderSource}
	}
	if id.PublisherId > 0 {
		condition.PublisherId = []int{id.PublisherId}
	}
	return condition
}

var (
	_orderStatisticsService     *OrderStatisticsService
	_orderStatisticsServiceOnce sync.Once
)

func GetOrderStatisticsService() *OrderStatisticsService {
	_orderStatisticsServiceOnce.Do(func() {
		if _orderStatisticsService == nil {
			_orderStatisticsService = new(OrderStatisticsService)
		}
	})
	return _orderStatisticsService
}
