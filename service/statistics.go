package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
)

type StatisticsService struct {
}

func (s *StatisticsService) Summary(ctx context.Context) (*entity.SummaryInfo, error){
	orgsCount ,err := da.GetOrgModel().CountOrgs(ctx, da.SearchOrgsCondition{
		ParentIDs: []int{0},
	})
	if err != nil{
		return nil, err
	}
	studentsCount, err := da.GetStudentModel().CountStudents(ctx)
	if err != nil{
		return nil, err
	}
	_, payRecords, err := da.GetOrderModel().SearchPayRecord(ctx, da.SearchPayRecordCondition{
		Mode:            entity.OrderPayModePay,
		StatusList:      []int{entity.OrderPayStatusChecked},
		PageSize: 			1000000,
	})
	if err != nil{
		return nil, err
	}
	performanceTotal := 0
	for i := range payRecords {
		performanceTotal = performanceTotal + payRecords[i].Amount
	}
	
	total, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		Status:         []int{entity.OrderStatusSigned, entity.OrderStatusRevoked, entity.OrderStatusCreated},
		Page:           1000000,
	})
	if err != nil{
		return nil, err
	}
	successTotal := 0
	failedTotal := 0
	for i := range orders{
		if orders[i].Status == entity.OrderStatusSigned{
			successTotal = successTotal + 1
		}else{
			failedTotal = failedTotal + 1
		}
	}

	return &entity.SummaryInfo{
		OrgsTotal:        orgsCount,
		StudentsTotal:    studentsCount,
		PerformanceTotal: performanceTotal,
		SuccessRate:      successTotal * 10000 / total,
	}, nil
}

func (s *StatisticsService) SearchYearRecords(ctx context.Context, key string) ([]*entity.StatisticRecord, error){
	year := time.Now().Year()
	records, err := da.GetStatisticsRecordModel().SearchStatisticsRecord(ctx, db.Get(), da.SearchStatisticsRecordCondition{
		Key:    key,
		Year:   year,
	})
	if err != nil{
		return nil, err
	}
	monthRecord := make(map[int]*da.StatisticsRecord)
	for i := range records {
		monthRecord[records[i].Month] = records[i]
	}

	ret := make([]*entity.StatisticRecord, 12)
	for i := 1; i <=12; i ++ {
		record, ok := monthRecord[i]
		if ok {
			ret[i - 1] = &entity.StatisticRecord{
				Key:  record.Key,
				Value: record.Value,
				Year:  record.Year,
				Month: record.Month,
			}
		}else{
			ret[i - 1] = &entity.StatisticRecord{
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
func (s *StatisticsService) AddPerformance(ctx context.Context, tx *gorm.DB, performance int) error {
	return s.addValue(ctx, tx, entity.PerformanceStatisticsKey, performance)
}

func (s *StatisticsService) addValue(ctx context.Context, tx *gorm.DB, key string, value int) error {
	now := time.Now()
	records, err := da.GetStatisticsRecordModel().SearchStatisticsRecord(ctx, tx, da.SearchStatisticsRecordCondition{
		Key:       key,
		Year:      now.Year(),
		Month:     int(now.Month()),
		Author:    0,
	})
	if err != nil{
		return err
	}
	if len(records) > 0 {
		record := records[0]
		record.Value = record.Value + value
		err = da.GetStatisticsRecordModel().UpdateStatisticsRecord(ctx, tx, record.ID, record.Value)
		if err != nil{
			return err
		}
		return nil
	}

	_, err = da.GetStatisticsRecordModel().CreateStatisticsRecord(ctx, tx, &da.StatisticsRecord{
		Key:       key,
		Value:     value,
		Year:      now.Year(),
		Month:     int(now.Month()),
		Author:    0,
	})
	if err != nil{
		return err
	}
	return nil
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
