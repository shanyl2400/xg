package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"
	"xg/utils"

	"github.com/jinzhu/gorm"
)

var (
	ErrPaymentsNotExists       = errors.New("parts of payments are not exists")
	ErrInvalidPaymentStatus    = errors.New("invalid payment status")
	ErrInvalidSettlementStatus = errors.New("invalid settlement status")
)

type ISettlementService interface {
	CreateSettlement(ctx context.Context, record entity.CreateSettlementRequest, operator *entity.JWTUser) error
	SearchSettlements(ctx context.Context, s da.SearchSettlementsCondition, operator *entity.JWTUser) (int, []*entity.SettlementData, error)
	CreateCommissionSettlement(ctx context.Context, tx *gorm.DB, record entity.CreateCommissionSettlementRequest, operator *entity.JWTUser) error
	SearchCommissionSettlements(ctx context.Context, s da.SearchCommissionSettlementsCondition, operator *entity.JWTUser) (int, []*entity.CommissionSettlementData, error)
}

type SettlementService struct {
}

func (s *SettlementService) CreateSettlement(ctx context.Context, record entity.CreateSettlementRequest, operator *entity.JWTUser) error {
	allPayment := append(record.FailedOrders, record.SuccessOrders...)

	total, _, err := da.GetOrderModel().SearchPayRecord(ctx, da.SearchPayRecordCondition{
		PayRecordIDList: allPayment,
	})
	// total, _, err := da.GetSettlementModel().SearchCommissionSettlements(ctx, da.SearchCommissionSettlementsCondition{
	// 	IDs: allPayment,
	// })
	if err != nil {
		log.Error.Printf("Can't search orders, allOrders: %#v, err: %v\n", allPayment, err)
		return err
	}
	if total != len(allPayment) {
		log.Error.Printf("Parts of orders are not exists, allOrders: %#v,total: %v err: %v\n", allPayment, total, ErrPaymentsNotExists)
		return ErrPaymentsNotExists
	}
	if record.Status != entity.SettlementStatusSettled &&
		record.Status != entity.SettlementStatusUnsettled {
		log.Error.Printf("Invalid settlement status, record: %#v,\n", record)
		return ErrInvalidSettlementStatus
	}

	err = da.GetSettlementModel().CreateSettlement(ctx, db.Get(), da.SettlementRecord{
		StartAt:       time.Unix(int64(record.StartAt), 0),
		EndAt:         time.Unix(int64(record.EndAt), 0),
		SuccessOrders: strings.Join(utils.IntsToStrings(record.SuccessOrders), ","),
		FailedOrders:  strings.Join(utils.IntsToStrings(record.FailedOrders), ","),
		Amount:        record.Amount,
		Status:        record.Status,
		Commission:    record.Commission,
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
			SuccessOrders: utils.ParseInts(res[i].SuccessOrders),
			FailedOrders:  utils.ParseInts(res[i].FailedOrders),
			Amount:        res[i].Amount,
			Status:        res[i].Status,
			Commission:    res[i].Commission,
			Invoice:       res[i].Invoice,
			AuthorID:      res[i].AuthorID,
			AuthorName:    authorNameMaps[res[i].AuthorID],
			UpdatedAt:     res[i].UpdatedAt,
			CreatedAt:     res[i].CreatedAt,
		}
	}
	return total, ret, nil
}

func (ss *SettlementService) CreateCommissionSettlement(ctx context.Context, tx *gorm.DB, record entity.CreateCommissionSettlementRequest, operator *entity.JWTUser) error {
	payment, err := da.GetOrderModel().GetPayRecordById(ctx, record.PaymentID)
	if err != nil {
		log.Error.Printf("Can't find payment, record: %#v, err: %v\n", record, err)
		return err
	}
	if payment.Status != entity.OrderPayStatusChecked {
		log.Error.Printf("Invalid payment status, record: %#v, payment: %#v, err: %v\n", record, payment, ErrInvalidPaymentStatus)
		return ErrInvalidPaymentStatus
	}

	err = da.GetSettlementModel().CreateCommissionSettlement(ctx, tx, da.CommissionSettlementRecord{
		OrderID:        payment.OrderID,
		PaymentID:      record.PaymentID,
		Amount:         payment.Amount,
		Commission:     record.Commission,
		SettlementNote: record.SettlementNote,
		Status:         entity.CommissionSettlementStatusCreated,

		AuthorID: operator.UserId,
		Note:     record.Note,
	})
	if err != nil {
		log.Error.Printf("Can't create settlement, record: %#v, payment: %#v, err: %v\n", record, payment, err)
		return err
	}
	err = da.GetOrderModel().UpdateOrderPayRecordTx(ctx, tx, record.PaymentID, entity.OrderPayStatusSettled)
	if err != nil {
		log.Error.Printf("Can't update pay record, record: %#v, payment: %#v, err: %v\n", record, payment, err)
		return err
	}
	return nil
}
func (ss *SettlementService) SearchCommissionSettlements(ctx context.Context, s da.SearchCommissionSettlementsCondition, operator *entity.JWTUser) (int, []*entity.CommissionSettlementData, error) {
	total, res, err := da.GetSettlementModel().SearchCommissionSettlements(ctx, s)
	if err != nil {
		log.Error.Printf("Can't search settlement, condition: %#v, err: %v\n", s, err)
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
	ret := make([]*entity.CommissionSettlementData, len(res))
	for i := range res {
		ret[i] = &entity.CommissionSettlementData{
			ID:             res[i].ID,
			OrderID:        res[i].OrderID,
			PaymentID:      res[i].PaymentID,
			Amount:         res[i].Amount,
			Commission:     res[i].Commission,
			Status:         res[i].Status,
			SettlementNote: res[i].SettlementNote,
			AuthorID:       res[i].AuthorID,
			Note:           res[i].Note,
			AuthorName:     authorNameMaps[res[i].AuthorID],
			UpdatedAt:      res[i].UpdatedAt,
			CreatedAt:      res[i].CreatedAt,
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
