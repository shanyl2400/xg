package service

import (
	"context"
	"errors"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoSuchOrderService = errors.New("no such order service")
)

type IOrderSourceService interface {
	ListOrderService(ctx context.Context) ([]*entity.OrderSource, error)
	CreateOrderService(ctx context.Context, name string) (int, error)
	DeleteOrderService(ctx context.Context, id int) error
	UpdateOrderSourceByID(ctx context.Context, req entity.UpdateOrderSourceRequest) error
}

type OrderSourceService struct {
}

func (o *OrderSourceService) ListOrderService(ctx context.Context) ([]*entity.OrderSource, error) {
	log.Info.Printf("list order services\n")
	os, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil {
		log.Warning.Printf("Get order source failed, err: %v\n", err)
		return nil, err
	}
	res := make([]*entity.OrderSource, len(os))
	for i := range os {
		res[i] = &entity.OrderSource{
			ID:   os[i].ID,
			Name: os[i].Name,
		}
	}
	return res, nil
}

func (o *OrderSourceService) CreateOrderService(ctx context.Context, name string) (int, error) {
	log.Info.Printf("CreateOrderService, req: %#v\n", name)
	return da.GetOrderSourceModel().CreateOrderSources(ctx, name)
}

func (o *OrderSourceService) UpdateOrderSourceByID(ctx context.Context, req entity.UpdateOrderSourceRequest) error {
	err := da.GetOrderSourceModel().UpdateOrderSourceByID(ctx, db.Get(), req.ID, req.Name)
	if err != nil {
		log.Error.Printf("Update Order service failed, req: %#v\n, err: %#v\n", req, err)
		return err
	}
	return nil
}

func (o *OrderSourceService) DeleteOrderService(ctx context.Context, id int) error {
	//删除订单来源
	os, err := da.GetOrderSourceModel().GetOrderSourceById(ctx, id)
	if err != nil {
		log.Error.Printf("Get Order service failed, id: %#v\n, err: %#v\n", id, err)
		return err
	}
	if os == nil {
		log.Error.Printf("No such order service, id: %#v\n", id)
		return ErrNoSuchOrderService
	}
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err := da.GetOrderSourceModel().DeleteOrderSourceByID(ctx, tx, id)
		if err != nil {
			log.Error.Printf("Delete Order failed, id: %#v\n, err: %#v\n", id, err)
			return err
		}
		//将学员名单订单来源更新为其他
		err = da.GetStudentModel().ReplaceStudentOrderSource(ctx, tx, id, entity.OtherOrderSource)
		if err != nil {
			log.Error.Printf("Replace Student Order source failed, id: %#v\n, err: %#v\n", id, err)
			return err
		}
		//将订单的订单来源更新为其他
		err = da.GetOrderModel().ReplaceOrderSource(ctx, tx, id, entity.OtherOrderSource)
		if err != nil {
			log.Error.Printf("Replace Order Order source failed, id: %#v\n, err: %#v\n", id, err)
			return err
		}

		//将统计该订单来源的数据更新为其他
		err = GetOrderStatisticsService().RemoveStatisticsByOrderSource(ctx, tx, id)
		if err != nil {
			log.Error.Printf("Update order statistics records failed, id: %#v\n, err: %#v\n", id, err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Warning.Printf("Delete Order service failed, id: %#v, err: %#v\n", id, err)
		return err
	}

	return nil
}

var (
	_orderSourceService     *OrderSourceService
	_orderSourceServiceOnce sync.Once
)

func GetOrderSourceService() *OrderSourceService {
	_orderSourceServiceOnce.Do(func() {
		if _orderSourceService == nil {
			_orderSourceService = new(OrderSourceService)
		}
	})
	return _orderSourceService
}
