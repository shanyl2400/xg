package service

import (
	"context"
	"sync"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type OrderSourceService struct {

}

func (o *OrderSourceService) ListOrderService(ctx context.Context)([]*entity.OrderSource, error){
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
	return da.GetOrderSourceModel().CreateOrderSources(ctx, name)
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
