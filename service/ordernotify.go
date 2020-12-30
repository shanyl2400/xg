package service

import (
	"context"
	"github.com/jinzhu/gorm"
	"sync"
	"xg/da"
	"xg/entity"
	"xg/log"
)

type IOrderNotifyService interface{
	NotifyOrderSignup(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error
	NotifyOrderDeposit(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error
	NotifyOrderRevoke(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error
	NotifyOrderInvalid(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error

	MarkNotifyRead(ctx context.Context, tx *gorm.DB, id int, operator *entity.JWTUser) error
	SearchPublisherNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser)(int, []*entity.OrderNotify, error)
	SearchAuthorNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser)(int, []*entity.OrderNotify, error)
	SearchNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser)(int, []*entity.OrderNotify, error)
}

type OrderNotifyService struct {

}

func (o *OrderNotifyService) NotifyOrderSignup(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error {
	data := da.OrderNotifies{
		ID:        0,
		OrderID:   orderID,
		Classify:  entity.OrderNotifyClassifySignup,
		Content:   content,
		Author:    operator.UserId,
		Status:    entity.OrderNotifyStatusUnread,
	}
	_, err := da.GetOrderNotifiesModel().CreateOrderNotify(ctx, tx, data)
	if err != nil{
		log.Error.Printf("Can't create order notify, data: %#v, err:%v", data, err)
		return err
	}
	return nil
}

func (o *OrderNotifyService) NotifyOrderDeposit(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error {
	data := da.OrderNotifies{
		ID:        0,
		OrderID:   orderID,
		Classify:  entity.OrderNotifyClassifyDeposit,
		Content:   content,
		Author:    operator.UserId,
		Status:    entity.OrderNotifyStatusUnread,
	}
	_, err := da.GetOrderNotifiesModel().CreateOrderNotify(ctx, tx, data)
	if err != nil{
		log.Error.Printf("Can't create order notify, data: %#v, err:%v", data, err)
		return err
	}
	return nil
}

func (o *OrderNotifyService) NotifyOrderRevoke(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error {
	data := da.OrderNotifies{
		ID:        0,
		OrderID:   orderID,
		Classify:  entity.OrderNotifyClassifyRevoke,
		Content:   content,
		Author:    operator.UserId,
		Status:    entity.OrderNotifyStatusUnread,
	}
	_, err := da.GetOrderNotifiesModel().CreateOrderNotify(ctx, tx, data)
	if err != nil{
		log.Error.Printf("Can't create order notify, data: %#v, err:%v", data, err)
		return err
	}
	return nil
}

func (o *OrderNotifyService) NotifyOrderInvalid(ctx context.Context, tx *gorm.DB, orderID int, content string, operator *entity.JWTUser) error {
	data := da.OrderNotifies{
		ID:        0,
		OrderID:   orderID,
		Classify:  entity.OrderNotifyClassifyInvalid,
		Content:   content,
		Author:    operator.UserId,
		Status:    entity.OrderNotifyStatusUnread,
	}
	_, err := da.GetOrderNotifiesModel().CreateOrderNotify(ctx, tx, data)
	if err != nil{
		log.Error.Printf("Can't create order notify, data: %#v, err:%v", data, err)
		return err
	}
	return nil
}

func (o *OrderNotifyService) MarkNotifyRead(ctx context.Context, tx *gorm.DB, id int, operator *entity.JWTUser) error {
	err := da.GetOrderNotifiesModel().UpdateOrderNotifyStatus(ctx, tx, id, entity.OrderNotifyStatusRead)
	if err != nil{
		log.Error.Printf("Can't create order notify, id: %#v, err:%v", id, err)
		return err
	}
	return nil
}
func (o *OrderNotifyService) SearchNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser)(int, []*entity.OrderNotify, error){
	total, notifies, err := da.GetOrderNotifiesModel().SearchOrderNotifies(ctx, condition)
	if err != nil{
		log.Error.Printf("search order notify failed, condition: %#v, err:%v", condition, err)
		return 0, nil, err
	}
	orderIDs := make([]int, len(notifies))
	for i := range notifies {
		orderIDs[i] = notifies[i].OrderID
	}
	res, err := GetOrderService().SearchOrders(ctx, &entity.SearchOrderCondition{
		IDs: orderIDs,
	}, operator)

	orderObj := make(map[int]*entity.OrderInfoDetails)
	for i := range res.Orders {
		orderObj[res.Orders[i].ID] = res.Orders[i]
	}

	ret := make([]*entity.OrderNotify, len(notifies))
	for i := range notifies {
		ret[i] = &entity.OrderNotify{
			ID:        notifies[i].ID,
			OrderID:   notifies[i].OrderID,
			Classify:  notifies[i].Classify,
			Content:   notifies[i].Content,
			Author:    notifies[i].Author,
			OrderInfo: orderObj[notifies[i].OrderID],
			Status:    notifies[i].Status,
			UpdatedAt: notifies[i].UpdatedAt,
			CreatedAt: notifies[i].CreatedAt,
		}
	}
	return total, ret, nil
}

func (o *OrderNotifyService) SearchAuthorNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser)(int, []*entity.OrderNotify, error){
	condition.OrderAuthorID = operator.UserId
	return o.SearchNotifies(ctx, condition, operator)
}
func (o *OrderNotifyService) SearchPublisherNotifies(ctx context.Context, condition da.OrderNotifiesCondition, operator *entity.JWTUser) (int, []*entity.OrderNotify, error) {
	condition.OrderPublisherID = operator.UserId
	return o.SearchNotifies(ctx, condition, operator)
}

var (
	_orderNotifyService     *OrderNotifyService
	_orderNotifyServiceOnce sync.Once
)

func GetOrderNotifyService() *OrderNotifyService {
	_orderNotifyServiceOnce.Do(func() {
		if _orderNotifyService == nil {
			_orderNotifyService = new(OrderNotifyService)
		}
	})
	return _orderNotifyService
}
