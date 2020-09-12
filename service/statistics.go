package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

type IStatisticsService interface{
	Summary(ctx context.Context) (*entity.SummaryInfo, error)
	SearchYearRecords(ctx context.Context, key string) ([]*entity.StatisticRecord, error)
}
type StatisticsService struct {
}

func (s *StatisticsService) Summary(ctx context.Context) (*entity.SummaryInfo, error) {
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
		Status: []int{entity.OrderStatusSigned, entity.OrderStatusRevoked, entity.OrderStatusCreated},
		Page:   1000000,
	}
	total, orders, err := da.GetOrderModel().SearchOrder(ctx, orderCondition)
	if err != nil {
		log.Warning.Printf("Search orders failed, condition: %#v, err: %v\n", orderCondition, err)
		return nil, err
	}
	successTotal := 0
	failedTotal := 0
	for i := range orders {
		if orders[i].Status == entity.OrderStatusSigned {
			successTotal = successTotal + 1
		} else {
			failedTotal = failedTotal + 1
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

func (s *StatisticsService) SearchYearRecords(ctx context.Context, key string) ([]*entity.StatisticRecord, error) {
	year := time.Now().Year()
	condition := da.SearchStatisticsRecordCondition{
		Key:  key,
		Year: year,
	}
	records, err := da.GetStatisticsRecordModel().SearchStatisticsRecord(ctx, db.Get(), condition)
	if err != nil {
		log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	monthRecord := make(map[int]*da.StatisticsRecord)
	for i := range records {
		monthRecord[records[i].Month] = records[i]
	}

	ret := make([]*entity.StatisticRecord, 12)
	for i := 1; i <= 12; i++ {
		record, ok := monthRecord[i]
		if ok {
			ret[i-1] = &entity.StatisticRecord{
				Key:   record.Key,
				Value: record.Value,
				Year:  record.Year,
				Month: record.Month,
			}
		} else {
			ret[i-1] = &entity.StatisticRecord{
				Key:   key,
				Value: 0,
				Year:  year,
				Month: i,
			}
		}
	}
	return ret, nil
}

func (s *StatisticsService) AddStudent(ctx context.Context, tx *gorm.DB, count int) error {
	return s.addValue(ctx, tx, entity.StudentStatisticsKey, count)
}
func (s *StatisticsService) AddPerformance(ctx context.Context, tx *gorm.DB, info entity.OrderPerformanceInfo, performance int) error {
	err := s.addValue(ctx, tx, StatisticKeyId(entity.OrgPerformanceStatisticsKey, info.OrgId), performance)
	if err != nil{
		return err
	}
	err = s.addValue(ctx, tx, StatisticKeyId(entity.AuthorPerformanceStatisticsKey, info.AuthorId), performance)
	if err != nil{
		return err
	}
	err = s.addValue(ctx, tx, StatisticKeyId(entity.PublisherPerformanceStatisticsKey, info.PublisherId), performance)
	if err != nil{
		return err
	}

	err = s.addValue(ctx, tx, entity.PerformanceStatisticsKey, performance)
	if err != nil{
		return err
	}
	return nil
}

func (s *StatisticsService) addValue(ctx context.Context, tx *gorm.DB, key string, value int) error {
	now := time.Now()
	condition := da.SearchStatisticsRecordCondition{
		Key:    key,
		Year:   now.Year(),
		Month:  int(now.Month()),
		Author: 0,
	}
	records, err := da.GetStatisticsRecordModel().SearchStatisticsRecord(ctx, tx, condition)
	if err != nil {
		log.Warning.Printf("add graph value failed, condition: %#v, err: %v\n", condition, err)
		return err
	}
	if len(records) > 0 {
		record := records[0]
		record.Value = record.Value + value
		err = da.GetStatisticsRecordModel().UpdateStatisticsRecord(ctx, tx, record.ID, record.Value)
		if err != nil {
			log.Warning.Printf("UpdateStatisticsRecord failed, record: %#v, err: %v\n", record, err)
			return err
		}
		return nil
	}

	record := &da.StatisticsRecord{
		Key:    key,
		Value:  value,
		Year:   now.Year(),
		Month:  int(now.Month()),
		Author: 0,
	}
	_, err = da.GetStatisticsRecordModel().CreateStatisticsRecord(ctx, tx, record)
	if err != nil {
		log.Warning.Printf("CreateStatisticsRecord failed, record: %#v, err: %v\n", record, err)
		return err
	}
	return nil
}

func StatisticKeyId(prefix string, id int) string{
	return fmt.Sprintf("%v-%v", prefix, id)
}

var (
	_statisticsService     *StatisticsService
	_statisticsServiceOnce sync.Once
)

func GetStatisticsService() *StatisticsService {
	_statisticsServiceOnce.Do(func() {
		if _statisticsService == nil {
			_statisticsService = new(StatisticsService)
		}
	})
	return _statisticsService
}
