package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
	"xg/utils"
)

var (
	ErrPaymentsNotExists    = errors.New("parts of payments are not exists")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
)

type ISettlementService interface {
	CreateSettlement(ctx context.Context, record entity.CreateSettlementRequest, operator *entity.JWTUser) error
	SearchSettlements(ctx context.Context, s da.SearchSettlementsCondition, operator *entity.JWTUser) (int, []*entity.SettlementData, error)
}

type SettlementService struct {
}

func (s *SettlementService) CreateSettlement(ctx context.Context, record entity.CreateSettlementRequest, operator *entity.JWTUser) error {
	allPaymentsStr := append(record.FailedOrders, record.SuccessOrders...)

	allPayment, err := utils.StringsToInts(allPaymentsStr)
	if err != nil {
		log.Error.Printf("Can't parse orders, allOrdersStr: %#v, err: %v\n", allPaymentsStr, err)
		return err
	}

	total, payments, err := da.GetOrderModel().SearchPayRecord(ctx, da.SearchPayRecordCondition{
		PayRecordIDList: allPayment,
	})
	if err != nil {
		log.Error.Printf("Can't search orders, allOrders: %#v, err: %v\n", allPayment, err)
		return err
	}
	if total != len(allPayment) {
		log.Error.Printf("Parts of orders are not exists, allOrders: %#v,total: %v err: %v\n", allPayment, total, ErrPaymentsNotExists)
		return ErrPaymentsNotExists
	}
	for i := range payments {
		if payments[i].Status != entity.OrderPayStatusChecked {
			log.Error.Printf("Invalid payment status, status: %#v,data: %#v err: %v\n", payments[i].Status, payments[i], ErrInvalidPaymentStatus)
			return ErrInvalidPaymentStatus
		}
	}

	err = da.GetSettlementModel().CreateSettlement(ctx, db.Get(), da.SettlementRecord{
		StartAt:       record.StartAt,
		EndAt:         record.EndAt,
		SuccessOrders: strings.Join(record.SuccessOrders, ","),
		FailedOrders:  strings.Join(record.FailedOrders, ","),
		Amount:        record.Amount,
		Status:        record.Status,
		Invoice:       record.Invoice,
		AuthorID:      operator.UserId,
	})
	if err != nil {
		log.Error.Printf("Can't create settlement, record: %#v, err: %v\n", record, err)
		return err
	}
	return nil
}
func (ss *SettlementService) SearchSettlements(ctx context.Context, sc da.SearchSettlementsCondition, operator *entity.JWTUser) (int, []*entity.SettlementData, error) {
	total, res, err := da.GetSettlementModel().SearchSettlements(ctx, sc)
	if err != nil {
		log.Error.Printf("Can't search settlement, condition: %#v, err: %v\n", sc, err)
		return 0, nil, err
	}
	authorIds := make([]int, len(res))
	for i := range res {
		authorIds[i] = res[i].AuthorID
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: utils.UniqueInts(authorIds),
	})
	if err != nil {
		log.Warning.Printf("Get User failed, ids: %#v, req: %#v, err: %v\n", authorIds, ss, err)
		return 0, nil, err
	}

	authorNameMaps := make(map[int]string)
	for i := range users {
		authorNameMaps[users[i].ID] = users[i].Name
	}

	ret := make([]*entity.SettlementData, len(res))
	for i := range res {
		ret[i] = &entity.SettlementData{
			ID:            res[i].ID,
			StartAt:       res[i].StartAt,
			EndAt:         res[i].EndAt,
			SuccessOrders: strings.Split(res[i].SuccessOrders, ","),
			FailedOrders:  strings.Split(res[i].FailedOrders, ","),
			Amount:        res[i].Amount,
			Status:        res[i].Status,
			Invoice:       res[i].Invoice,
			AuthorID:      res[i].AuthorID,
			AuthorName:    authorNameMaps[res[i].AuthorID],
			UpdatedAt:     res[i].UpdatedAt,
			CreatedAt:     res[i].CreatedAt,
		}
	}
	return total, ret, nil
}

var (
	_settlementService     *SettlementService
	_settlementServiceOnce sync.Once
)

func GetSettlementService() *SettlementService {
	_settlementServiceOnce.Do(func() {
		if _settlementService == nil {
			_settlementService = new(SettlementService)
		}
	})
	return _settlementService
}
