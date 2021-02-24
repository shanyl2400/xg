package service

import (
	"context"
	"sync"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type IOrderSourceService interface{
	ListOrderService(ctx context.Context)([]*entity.OrderSource, error)
	CreateOrderService(ctx context.Context, name string)(int, error)
	DeleteOrderService(ctx context.Context, id int) error
}

type OrderSourceService struct {
}

func (o *OrderSourceService) ListOrderService(ctx context.Context)([]*entity.OrderSource, error){
	log.Info.Printf("list order services\n")
	os, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil{
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

func (o *OrderSourceService) CreateOrderService(ctx context.Context, name string)(int, error){
	log.Info.Printf("CreateOrderService, req: %#v\n", name)
	return da.GetOrderSourceModel().CreateOrderSources(ctx, name)
}

func (o *OrderSourceService)DeleteOrderService(ctx context.Context, id int) error {
	//删除订单来源
	//将学员名单订单来源更新为其他
	//将订单的订单来源更新为其他
	//将统计该订单来源的数据更新为其他

	return nil
}

var(
	_orderSourceService *OrderSourceService
	_orderSourceServiceOnce sync.Once
)

func GetOrderSourceService() *OrderSourceService{
	_orderSourceServiceOnce.Do(func() {
		if _orderSourceService == nil{
			_orderSourceService = new(OrderSourceService)
		}
	})
	return _orderSourceService
}
