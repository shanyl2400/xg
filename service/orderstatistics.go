package service

import (
	"context"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

const (
	DayTimeDuration = time.Hour * 24
)

type OrderStatisticsService struct {
}

func (s *OrderStatisticsService) StatisticsGroupByOrgs(ctx context.Context,
	et entity.StatisticRecordCondition) ([]*entity.OrderStatisticGroupTableItem, error) {
	records, err := s.fetchStatisticsRecords(ctx, &et)
	if err != nil {
		return nil, err
	}

	//统计结果
	orgsRecordsMap := make(map[int]entity.OrderStatisticTableMonth)
	for i := range records {
		//若没有创建value，则创建
		_, ok := orgsRecordsMap[records[i].OrgId]
		if !ok {
			orgsRecordsMap[records[i].OrgId] = entity.OrderStatisticTableMonth{}
		}
		//过滤时间
		recordDate := time.Date(records[i].Year, time.Month(records[i].Month), records[i].Date, 0, 0, 0, 0, time.Local)
		//在end_at之后或在start_at之前,排除
		if recordDate.After(*et.EndAt) || recordDate.Before(*et.StartAt) {
			continue
		}

		//统计records
		orgsRecordsMap[records[i].OrgId] = s.handleRecord(ctx, orgsRecordsMap[records[i].OrgId], records[i])
	}

	_, orgInfo, err := GetOrgService().ListOrgs(ctx, da.SearchOrgsCondition{
		IsSubOrg: false,
	})
	if err != nil {
		log.Error.Printf("ListOrderService failed, err: %v\n", err)
		return nil, err
	}
	orgMap := make(map[int]string)
	for i := range orgInfo {
		orgMap[orgInfo[i].ID] = orgInfo[i].Name
	}

	ret := make([]*entity.OrderStatisticGroupTableItem, 0)
	for orgID, orgName := range orgMap {
		item, ok := orgsRecordsMap[orgID]
		if !ok {
			item = entity.OrderStatisticTableMonth{
				Students:      0,
				Orders:        0,
				InvalidOrders: 0,
				SignedOrder:   0,
				Performance:   0,
			}
		}
		temp := entity.OrderStatisticTableItem{
			OrderStatisticTableMonth: item,
			Succeed:                  0,
		}
		temp.CalculateSucceed()
		ret = append(ret, &entity.OrderStatisticGroupTableItem{
			GroupID:                 orgID,
			GroupName:               orgName,
			OrderStatisticTableItem: temp,
		})
	}

	return ret, nil
}

func (s *OrderStatisticsService) StatisticsGroupByOrderSources(ctx context.Context,
	et entity.StatisticRecordCondition) ([]*entity.OrderStatisticGroupTableItem, error) {
	records, err := s.fetchStatisticsRecords(ctx, &et)
	if err != nil {
		return nil, err
	}

	//统计结果
	orderSourceRecordsMap := make(map[int]entity.OrderStatisticTableMonth)
	for i := range records {
		//若没有创建value，则创建
		_, ok := orderSourceRecordsMap[records[i].OrderSource]
		if !ok {
			orderSourceRecordsMap[records[i].OrderSource] = entity.OrderStatisticTableMonth{}
		}
		//过滤时间
		recordDate := time.Date(records[i].Year, time.Month(records[i].Month), records[i].Date, 0, 0, 0, 0, time.Local)
		//在end_at之后或在start_at之前,排除
		if recordDate.After(*et.EndAt) || recordDate.Before(*et.StartAt) {
			continue
		}

		//统计records
		orderSourceRecordsMap[records[i].OrderSource] = s.handleRecord(ctx, orderSourceRecordsMap[records[i].OrderSource], records[i])
	}

	orderSources, err := GetOrderSourceService().ListOrderService(ctx)
	if err != nil {
		log.Error.Printf("ListOrderService failed, err: %v\n", err)
		return nil, err
	}
	orderSourceMap := make(map[int]string)
	for i := range orderSources {
		orderSourceMap[orderSources[i].ID] = orderSources[i].Name
	}

	ret := make([]*entity.OrderStatisticGroupTableItem, 0)
	for orderSourceID, orderSourceName := range orderSourceMap {
		item, ok := orderSourceRecordsMap[orderSourceID]
		if !ok {
			item = entity.OrderStatisticTableMonth{
				Students:      0,
				Orders:        0,
				InvalidOrders: 0,
				SignedOrder:   0,
				Performance:   0,
			}
		}
		temp := entity.OrderStatisticTableItem{
			OrderStatisticTableMonth: item,
			Succeed:                  0,
		}
		temp.CalculateSucceed()
		ret = append(ret, &entity.OrderStatisticGroupTableItem{
			GroupID:                 orderSourceID,
			GroupName:               orderSourceName,
			OrderStatisticTableItem: temp,
		})
	}

	return ret, nil
}

//名单数，无效人数，报名人数，成交业绩，成功率
func (s *OrderStatisticsService) StatisticsTable(ctx context.Context, et entity.OrderStatisticRecordEntity) (*entity.OrderStatisticTable, error) {
	records, err := s.searchLast3MonthRecords(ctx, et)
	if err != nil {
		return nil, err
	}
	ret := entity.NewOrderStatisticTable()
	now := time.Now()
	for i := range records {
		recordDate := time.Date(records[i].Year, time.Month(records[i].Month), records[i].Date, 0, 0, 0, 0, time.Local)
		//当日数据
		if records[i].Year == now.Year() &&
			records[i].Month == int(now.Month()) &&
			records[i].Date == now.Day() {
			ret.DayData.OrderStatisticTableMonth = s.handleRecord(ctx, ret.DayData.OrderStatisticTableMonth, records[i])
		}

		//本周数据
		rYear, rWeek := recordDate.ISOWeek()
		nYear, nWeek := now.ISOWeek()
		if rYear == nYear && rWeek == nWeek {
			ret.WeekDayData.OrderStatisticTableMonth = s.handleRecord(ctx, ret.WeekDayData.OrderStatisticTableMonth, records[i])
		}

		//当月数据
		if records[i].Year == now.Year() &&
			records[i].Month == int(now.Month()) {
			ret.MonthDayData.OrderStatisticTableMonth = s.handleRecord(ctx, ret.MonthDayData.OrderStatisticTableMonth, records[i])
		}

		//三个月数据
		timeDiff := now.Sub(recordDate)
		if timeDiff < DayTimeDuration*90 {
			ret.ThreeMonthDayData.OrderStatisticTableMonth = s.handleRecord(ctx, ret.ThreeMonthDayData.OrderStatisticTableMonth, records[i])
		}

		//本年度
		if records[i].Year == now.Year() {
			month := records[i].Month - 1
			if ret.Data[month] == nil {
				ret.Data[month] = new(entity.OrderStatisticTableMonth)
			}
			data := s.handleRecord(ctx, *ret.Data[month], records[i])
			ret.Data[month] = &data
		}
	}
	//计算成功率
	ret.CalculateSucceed()
	return ret, nil
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
	payCondition := da.SearchPayRecordCondition{
		// Mode:       entity.OrderPayModePay,
		StatusList: []int{entity.OrderPayStatusChecked},
		PageSize:   1000000,
	}
	_, payRecords, err := da.GetOrderModel().SearchPayRecord(ctx, payCondition)
	if err != nil {
		log.Warning.Printf("Search pay records failed, condition: %#v, err: %v\n", payCondition, err)
		return nil, err
	}
	performanceTotal := float64(0)
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
		Page: 1000000,
	}
	total, orders, err := da.GetOrderModel().SearchOrder(ctx, orderCondition)
	if err != nil {
		log.Warning.Printf("Search orders failed, condition: %#v, err: %v\n", orderCondition, err)
		return nil, err
	}
	successTotal := 0
	for i := range orders {
		if orders[i].Status == entity.OrderStatusSigned ||
			orders[i].Status == entity.OrderStatusDeposit {
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
			value := float64(0)
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

func (s *OrderStatisticsService) AddConsiderOrder(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity) error {
	return s.addOrder(ctx, tx, osr, entity.OrderStatisticKeyConsiderOrder, 1)
}

func (s *OrderStatisticsService) AddPerformance(ctx context.Context, tx *gorm.DB, osr entity.OrderStatisticRecordEntity, performance float64) error {
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
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderStatisticsService) fetchStatisticsRecords(ctx context.Context,
	et *entity.StatisticRecordCondition) ([]*da.OrderStatisticsRecord, error) {
	//检查时间参数是否有误，若没有设置，则返回当年的数据
	if et.StartAt == nil || et.EndAt == nil {
		curYear := time.Now().Year()
		startAt := time.Date(curYear, 1, 1, 0, 0, 0, 0, time.Local)

		now := time.Now()
		et.StartAt = &startAt
		et.EndAt = &now
	}
	//查看时间是否有误
	if et.StartAt.After(*et.EndAt) {
		return nil, ErrInvalidSearchTime
	}

	endYear := et.EndAt.Year()
	startYear := et.StartAt.Year()
	yearArr := make([]int, 0)
	//查询不同年份的数据再筛除
	for i := startYear; i <= endYear; i++ {
		yearArr = append(yearArr, i)
	}
	//yearArr = utils.SliceDeduplicationInt(yearArr)

	//构建查询条件
	condition := s.entityToCondition(ctx, et.OrderStatisticRecordEntity)
	condition.Year = yearArr
	records, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
	if err != nil {
		log.Error.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	return records, nil
}
func (s *OrderStatisticsService) addValue(ctx context.Context, tx *gorm.DB, id entity.OrderStatisticRecordId, value float64, addCount bool) error {
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
		Key:         id.Key,
		Value:       value,
		Year:        now.Year(),
		Month:       int(now.Month()),
		Date:        now.Day(),
		Count:       count,
		Author:      id.Author,
		OrgId:       id.OrgId,
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
	}, float64(count), true)
}

func (s *OrderStatisticsService) searchLast3MonthRecords(ctx context.Context, et entity.OrderStatisticRecordEntity) ([]*da.OrderStatisticsRecord, error) {
	condition := s.entityToCondition(ctx, et)
	now := time.Now()

	records := make([]*da.OrderStatisticsRecord, 0)
	switch int(now.Month()) {
	case 1:
		condition.Year = []int{now.Year()}
		condition.Month = []int{1}
		records0, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
		if err != nil {
			log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}
		condition.Year = []int{now.Year() - 1}
		condition.Month = []int{11, 12}
		records1, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
		if err != nil {
			log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}
		records = append(records0, records1...)
	case 2:
		condition.Year = []int{now.Year()}
		condition.Month = []int{1, 2}
		records0, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
		if err != nil {
			log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}

		condition.Year = []int{now.Year() - 1}
		condition.Month = []int{12}
		records1, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
		if err != nil {
			log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}

		records = append(records0, records1...)
	default:
		month := int(now.Month())
		condition.Month = []int{month, month - 1, month - 2}
		condition.Year = []int{now.Year()}
		records0, err := da.GetOrderStatisticsRecordModel().SearchOrderStatisticsRecord(ctx, db.Get(), condition)
		if err != nil {
			log.Warning.Printf("SearchStatisticsRecord failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}
		records = records0
	}
	return records, nil
}

func (s *OrderStatisticsService) handleRecord(ctx context.Context, item entity.OrderStatisticTableMonth, record *da.OrderStatisticsRecord) entity.OrderStatisticTableMonth {
	switch record.Key {
	case entity.OrderStatisticKeyStudent:
		item.Students = item.Students + int(record.Value)
	case entity.OrderStatisticKeyOrder:
		item.Performance = item.Performance + record.Value
	case entity.OrderStatisticKeyNewOrder:
		item.Orders = item.Orders + int(record.Value)
	case entity.OrderStatisticKeySignupOrder:
		item.SignedOrder = item.SignedOrder + int(record.Value)
	case entity.OrderStatisticKeyInvalidOrder:
		item.InvalidOrders = item.InvalidOrders + int(record.Value)
	}
	return item
}

func (s *OrderStatisticsService) idToCondition(id entity.OrderStatisticRecordId) da.SearchOrderStatisticsRecordCondition {
	now := time.Now()
	condition := da.SearchOrderStatisticsRecordCondition{
		Key:   id.Key,
		Year:  []int{now.Year()},
		Month: []int{int(now.Month())},
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

func (s *OrderStatisticsService) entityToCondition(ctx context.Context, id entity.OrderStatisticRecordEntity) da.SearchOrderStatisticsRecordCondition {
	condition := da.SearchOrderStatisticsRecordCondition{}

	if id.Author > 0 {
		condition.Author = []int{id.Author}
	}
	if id.OrgId > 0 {
		subOrgs, err := GetOrgService().GetSubOrgs(ctx, id.OrgId)
		if err != nil {
			log.Warning.Println("Can't get sub orgs, error:", err)
			condition.OrgId = []int{id.OrgId}
		} else {
			ids := make([]int, len(subOrgs))
			for i := range ids {
				ids[i] = subOrgs[i].ID
			}
			ids = append(ids, id.OrgId)
			condition.OrgId = ids
		}

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
