package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/log"

	"github.com/jinzhu/gorm"
)

var (
	ErrInvalidOrderRemarkStatus = errors.New("invalid order remark status")
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, operator *entity.JWTUser) (int, error)

	SignUpOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error
	DepositOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error
	RevokeOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error
	InvalidOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error
	ConsiderOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error

	UpdateOrderStatus(ctx context.Context, req *entity.OrderUpdateStatusRequest, operator *entity.JWTUser) error

	PayOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error
	PaybackOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error

	ConfirmOrderPay(ctx context.Context, orderPayId int, status int, operator *entity.JWTUser) error
	UpdateOrderPayPrice(ctx context.Context, orderPayId int, price float64, operator *entity.JWTUser) error
	AddOrderRemark(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error
	MarkOrderRemark(ctx context.Context, remarkIDs []int, status int, operator *entity.JWTUser) error
	SearchOrderPayRecords(ctx context.Context, condition *entity.SearchPayRecordCondition, operator *entity.JWTUser) (*entity.PayRecordInfoList, error)
	SearchOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error)
	SearchOrderWithAuthor(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error)
	SearchOrderWithOrgId(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error)
	SearchOrderRemarks(ctx context.Context, condition *da.SearchRemarkRecordCondition, operator *entity.JWTUser) (*entity.OrderRemarkList, error)
	GetOrderById(ctx context.Context, orderId int, operator *entity.JWTUser) (*entity.OrderInfoWithRecords, error)

	ExportOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) ([]byte, error)
}

type OrderService struct {
	lock sync.Mutex
}
type orderEntity struct {
	StudentID int `json:"student_id"`
	ToOrgID   int `json:"to_org_id"`
}

func (o *OrderService) CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, operator *entity.JWTUser) (int, error) {
	//Check entity
	org, err := da.GetOrgModel().GetOrgById(ctx, db.Get(), req.ToOrgID)
	if err != nil {
		log.Warning.Println("Can't find org when check entity, orderEntity:", req)
		return -1, err
	}
	err = o.checkEntity(ctx, org, orderEntity{
		StudentID: req.StudentID,
		ToOrgID:   req.ToOrgID,
	})
	if err != nil {
		log.Warning.Printf("Create order failed, req: %#v, err: %v\n", req, err)
		return -1, err
	}
	student, err := da.GetStudentModel().GetStudentById(ctx, req.StudentID)
	if err != nil {
		log.Warning.Printf("Get Student failed, req: %#v, err: %v\n", req, err)
		return -1, err
	}

	//若学员状态非法，不能创建
	if student.Status != entity.StudentCreated && student.Status != entity.StudentConflictSuccess {
		log.Warning.Printf("Student is conflict, req: %#v", req)
		return -1, ErrStudentIsConflict
	}

	//TODO:检查重复订单？
	data := da.Order{
		StudentID:      req.StudentID,
		ToOrgID:        req.ToOrgID,
		Address:        org.Address,
		IntentSubjects: strings.Join(req.IntentSubjects, ","),
		PublisherID:    operator.UserId,
		AuthorID:       student.AuthorID,
		OrderSource:    student.OrderSourceID,
		Status:         entity.OrderStatusCreated,
	}

	id, err := db.GetTransResult(ctx, func(ctx context.Context, tx *gorm.DB) (interface{}, error) {
		log.Info.Printf("create order: %#v\n", data)
		id, err := da.GetOrderModel().CreateOrder(ctx, tx, data)
		if err != nil {
			log.Warning.Printf("Create order failed, data: %#v, err: %v\n", data, err)
			return -1, err
		}

		record := entity.OrderStatisticRecordEntity{
			Author:      student.AuthorID,
			OrgId:       req.ToOrgID,
			PublisherId: operator.UserId,
			OrderSource: student.OrderSourceID,
		}
		err = GetOrderStatisticsService().AddNewOrder(ctx, tx, record)
		if err != nil {
			log.Warning.Printf("Add statistics record failed, record: %#v, err: %v\n", record, err)
			return -1, err
		}

		err = GetStudentService().UpdateStudentOrderCount(ctx, tx, req.StudentID, 1)
		if err != nil {
			log.Warning.Printf("Update student order count failed, req: %#v, err: %v\n", req, err)
			return -1, err
		}

		return id, nil
	})
	if err != nil {
		return -1, err
	}

	return id.(int), nil
}

func (o *OrderService) UpdateOrderStatus(ctx context.Context, req *entity.OrderUpdateStatusRequest, operator *entity.JWTUser) error {
	switch req.Status {
	case entity.OrderStatusCreated:
		log.Warning.Printf("Invalid status, req: %#v\n", req)
		return ErrInvalidOrderStatus
	case entity.OrderStatusSigned:
		return o.SignUpOrder(ctx, &entity.OrderPayRequest{
			OrderID: req.OrderID,
			Amount:  req.Amount,
			Title:   req.Title,
			Content: req.Content,
		}, operator)
	case entity.OrderStatusRevoked:
		return o.RevokeOrder(ctx, req.OrderID, req.Content, operator)
	case entity.OrderStatusInvalid:
		return o.InvalidOrder(ctx, req.OrderID, req.Content, operator)
	case entity.OrderStatusDeposit:
		return o.DepositOrder(ctx, &entity.OrderPayRequest{
			OrderID: req.OrderID,
			Amount:  req.Amount,
			Title:   req.Title,
			Content: req.Content,
		}, operator)
	case entity.OrderStatusConsider:
		return o.ConsiderOrder(ctx, req.OrderID, req.Content, operator)
	}
	log.Warning.Printf("Invalid status, req: %#v\n", req)
	return ErrInvalidOrderStatus
}

//报名
func (o *OrderService) SignUpOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	if req.OrderID < 1 {
		log.Warning.Printf("Invalid orderId, orderId: %#v\n", req.OrderID)
		return ErrInvalidOrderId
	}
	//检查order状态
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		log.Warning.Printf("Get org failed, req: %#v, err: %v\n", req, err)
		return err
	}

	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check org order failed, req: %#v, operator: %#v, err: %v\n", req, operator, err)
		return err
	}

	if orderObj.Order.Status != entity.OrderStatusCreated &&
		orderObj.Order.Status != entity.OrderStatusDeposit &&
		orderObj.Order.Status != entity.OrderStatusConsider {
		log.Warning.Printf("Check order status failed, req: %#v, order: %#v, err: %v\n", req, orderObj, err)
		return ErrNoAuthToOperateOrder
	}

	log.Info.Printf("Sign order, req: %#v\n", req)
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		//修改order状态
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, req.OrderID, entity.OrderStatusSigned)
		if err != nil {
			log.Warning.Printf("Update order status failed, req: %#v, operator: %#v, err: %v\n", req, operator, err)
			return err
		}
		if req.Title == "" {
			req.Title = "报名费用"
		}
		//增加orderPay
		payData := &da.OrderPayRecord{
			OrderID: req.OrderID,
			Mode:    entity.OrderPayModePay,
			Title:   req.Title,
			Amount:  req.Amount,
			Status:  entity.OrderPayStatusPending,
		}
		_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, payData)
		if err != nil {
			log.Warning.Printf("Add order pay failed, req: %#v, payData: %#v, err: %v\n", req, payData, err)
			return err
		}

		//添加remarks
		content := fmt.Sprintf("订单已报名,缴费:%v元", req.Amount)
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  req.OrderID,
			InfoType: entity.OrderRemarkInfoTypeSignup,
			Info:     content,
			Content:  req.Content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, req: %#v, content: %#v, err: %v\n", req, content, err)
			return err
		}

		student, err := GetStudentService().GetStudentById(ctx, orderObj.Order.StudentID, operator)
		if err != nil {
			log.Warning.Printf("Get student failed, StudentID: %#v, err: %v\n", orderObj.Order.StudentID, err)
			return err
		}
		record := entity.OrderStatisticRecordEntity{
			Author:      student.AuthorID,
			OrgId:       orderObj.Order.ToOrgID,
			PublisherId: orderObj.Order.PublisherID,
			OrderSource: student.OrderSourceID,
		}
		err = GetOrderStatisticsService().AddSignupOrder(ctx, tx, record)
		if err != nil {
			log.Warning.Printf("Add statistics record failed, record: %#v, err: %v\n", record, err)
			return err
		}

		//添加通知
		err = GetOrderNotifyService().NotifyOrderSignup(ctx, tx, req.OrderID, content, operator)
		if err != nil {
			log.Warning.Printf("Add Notify, content: %#v, err: %v\n", content, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//交定金
func (o *OrderService) DepositOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	if req.OrderID < 1 {
		log.Warning.Printf("Invalid orderId, orderId: %#v\n", req.OrderID)
		return ErrInvalidOrderId
	}

	//检查order状态
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		log.Warning.Printf("Get org failed, req: %#v, err: %v\n", req, err)
		return err
	}

	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check org order failed, req: %#v, operator: %#v, err: %v\n", req, operator, err)
		return err
	}

	if orderObj.Order.Status != entity.OrderStatusCreated && orderObj.Order.Status != entity.OrderStatusConsider {
		log.Warning.Printf("Check order status failed, req: %#v, order: %#v, err: %v\n", req, orderObj, err)
		return ErrNoAuthToOperateOrder
	}

	log.Info.Printf("Deposit order, req: %#v\n", req)
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		//修改order状态
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, req.OrderID, entity.OrderStatusDeposit)
		if err != nil {
			log.Warning.Printf("Update order status failed, req: %#v, operator: %#v, err: %v\n", req, operator, err)
			return err
		}
		if req.Title == "" {
			req.Title = "报名定金"
		}
		//增加orderPay
		payData := &da.OrderPayRecord{
			OrderID: req.OrderID,
			Mode:    entity.OrderPayModePay,
			Title:   req.Title,
			Amount:  req.Amount,
			Status:  entity.OrderPayStatusPending,
		}
		_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, payData)
		if err != nil {
			log.Warning.Printf("Add order pay failed, req: %#v, payData: %#v, err: %v\n", req, payData, err)
			return err
		}

		content := fmt.Sprintf("订单交定金,金额:%v元", req.Amount)
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  req.OrderID,
			InfoType: entity.OrderRemarkInfoTypeDeposit,
			Info:     content,
			Content:  req.Content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, req: %#v, content: %#v, err: %v\n", req, content, err)
			return err
		}

		//添加通知
		err = GetOrderNotifyService().NotifyOrderDeposit(ctx, tx, req.OrderID, content, operator)
		if err != nil {
			log.Warning.Printf("Add Notify, content: %#v, err: %v\n", content, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//退费订单
func (o *OrderService) RevokeOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	if orderId < 1 {
		log.Warning.Printf("Invalid orderId, orderId: %#v\n", orderId)
		return ErrInvalidOrderId
	}
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		log.Warning.Printf("Revoke order failed, orderId: %#v, err: %v\n", orderId, err)
		return err
	}

	log.Info.Printf("Revoke order, order: %#v\n", orderObj)
	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check order org failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return err
	}

	if orderObj.Order.Status != entity.OrderStatusSigned &&
		orderObj.Order.Status != entity.OrderStatusDeposit {
		log.Warning.Printf("Check order status failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil
	}
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, orderId, entity.OrderStatusRevoked)
		if err != nil {
			log.Warning.Printf("Update order status failed, order: %#v, orderId: %#v, err: %v\n", orderObj, orderId, err)
			return err
		}

		info := fmt.Sprintf("订单已退学")
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  orderId,
			InfoType: entity.OrderRemarkInfoTypeRevoke,
			Info:     info,
			Content:  content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, orderId: %#v, content: %#v, err: %v\n", orderId, content, err)
			return err
		}

		student, err := GetStudentService().GetStudentById(ctx, orderObj.Order.StudentID, operator)
		if err != nil {
			log.Warning.Printf("Get student failed, StudentID: %#v, err: %v\n", orderObj.Order.StudentID, err)
			return err
		}
		record := entity.OrderStatisticRecordEntity{
			Author:      student.AuthorID,
			OrgId:       orderObj.Order.ToOrgID,
			PublisherId: orderObj.Order.PublisherID,
			OrderSource: student.OrderSourceID,
		}
		err = GetOrderStatisticsService().AddInvalidOrder(ctx, tx, record)
		if err != nil {
			log.Warning.Printf("Add statistics record failed, record: %#v, err: %v\n", record, err)
			return err
		}

		//添加通知
		err = GetOrderNotifyService().NotifyOrderDeposit(ctx, tx, orderId, info, operator)
		if err != nil {
			log.Warning.Printf("Add Notify, content: %#v, err: %v\n", info, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderService) ConsiderOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	if orderId < 1 {
		log.Warning.Printf("Consider orderId, orderId: %#v\n", orderId)
		return ErrInvalidOrderId
	}
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		log.Warning.Printf("Consider order failed, orderId: %#v, err: %v\n", orderId, err)
		return err
	}

	log.Info.Printf("Consider order, order: %#v\n", orderObj)
	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check order org failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return err
	}
	if orderObj.Order.Status != entity.OrderStatusCreated {
		log.Warning.Printf("Check order status failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil
	}

	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		//更新订单状态
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, db.Get(), orderId, entity.OrderStatusConsider)
		if err != nil {
			log.Warning.Printf("Update order status failed, order: %#v, orderId: %#v, err: %v\n", orderObj, orderId, err)
			return err
		}
		info := fmt.Sprintf("订单考虑中")
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  orderId,
			InfoType: entity.OrderRemarkInfoTypeConsider,
			Info:     info,
			Content:  content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, orderId: %#v, content: %#v, err: %v\n", orderId, content, err)
			return err
		}

		student, err := GetStudentService().GetStudentById(ctx, orderObj.Order.StudentID, operator)
		if err != nil {
			log.Warning.Printf("Get student failed, StudentID: %#v, err: %v\n", orderObj.Order.StudentID, err)
			return err
		}
		record := entity.OrderStatisticRecordEntity{
			Author:      student.AuthorID,
			OrgId:       orderObj.Order.ToOrgID,
			PublisherId: orderObj.Order.PublisherID,
			OrderSource: student.OrderSourceID,
		}
		err = GetOrderStatisticsService().AddConsiderOrder(ctx, tx, record)
		if err != nil {
			log.Warning.Printf("Add statistics record failed, record: %#v, err: %v\n", record, err)
			return err
		}

		//添加通知
		err = GetOrderNotifyService().NotifyOrderDeposit(ctx, tx, orderId, info, operator)
		if err != nil {
			log.Warning.Printf("Add Notify, content: %#v, err: %v\n", info, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//无效订单
func (o *OrderService) InvalidOrder(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	if orderId < 1 {
		log.Warning.Printf("Invalid orderId, orderId: %#v\n", orderId)
		return ErrInvalidOrderId
	}
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		log.Warning.Printf("Invalid order failed, orderId: %#v, err: %v\n", orderId, err)
		return err
	}

	log.Info.Printf("Invalid order, order: %#v\n", orderObj)
	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check order org failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return err
	}

	if orderObj.Order.Status != entity.OrderStatusCreated {
		log.Warning.Printf("Check order status failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil
	}
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, db.Get(), orderId, entity.OrderStatusInvalid)
		if err != nil {
			log.Warning.Printf("Update order status failed, order: %#v, orderId: %#v, err: %v\n", orderObj, orderId, err)
			return err
		}

		info := fmt.Sprintf("订单为无效订单")
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  orderId,
			InfoType: entity.OrderRemarkInfoTypeInvalid,
			Info:     info,
			Content:  content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, orderId: %#v, content: %#v, err: %v\n", orderId, content, err)
			return err
		}

		student, err := GetStudentService().GetStudentById(ctx, orderObj.Order.StudentID, operator)
		if err != nil {
			log.Warning.Printf("Get student failed, StudentID: %#v, err: %v\n", orderObj.Order.StudentID, err)
			return err
		}
		record := entity.OrderStatisticRecordEntity{
			Author:      student.AuthorID,
			OrgId:       orderObj.Order.ToOrgID,
			PublisherId: orderObj.Order.PublisherID,
			OrderSource: student.OrderSourceID,
		}
		err = GetOrderStatisticsService().AddInvalidOrder(ctx, tx, record)
		if err != nil {
			log.Warning.Printf("Add statistics record failed, record: %#v, err: %v\n", record, err)
			return err
		}

		err = GetStudentService().UpdateStudentOrderCount(ctx, tx, orderObj.Order.StudentID, -1)
		if err != nil {
			log.Warning.Printf("Update student order count failed, orderObj: %#v, err: %v\n", orderObj, err)
			return err
		}

		//添加通知
		err = GetOrderNotifyService().NotifyOrderDeposit(ctx, tx, orderId, info, operator)
		if err != nil {
			log.Warning.Printf("Add Notify, content: %#v, err: %v\n", info, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

//付款
func (o *OrderService) PayOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	return o.payOrder(ctx, entity.OrderPayModePay, req, operator)
}

//退费
func (o *OrderService) PaybackOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	return o.payOrder(ctx, entity.OrderPayModePayback, req, operator)
}

//确认
func (o *OrderService) ConfirmOrderPay(ctx context.Context, orderPayId int, status int, operator *entity.JWTUser) error {
	//所有payRecord都确认，才能将order确认
	//检查order状态
	//修改支付记录状态
	o.lock.Lock()
	defer o.lock.Unlock()
	log.Info.Printf("Comfirm order pay, payId: %#v, status: %#v\n", orderPayId, status)
	err := db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err := da.GetOrderModel().UpdateOrderPayRecordTx(ctx, tx, orderPayId, status)
		if err != nil {
			log.Warning.Printf("Update order pay record failed, orderPayId: %#v, err: %v\n", orderPayId, err)
			return err
		}
		if status == entity.OrderPayStatusChecked {
			payment, err := da.GetOrderModel().GetPayRecordById(ctx, orderPayId)
			if err != nil {
				log.Warning.Printf("Get pay record failed, orderPayId: %#v, err: %v\n", orderPayId, err)
				return err
			}
			performance := payment.Amount
			if payment.Mode == entity.OrderPayModePayback {
				performance = -performance
			}
			orderInfo, err := o.GetOrderById(ctx, payment.OrderID, operator)

			err = GetOrderStatisticsService().AddPerformance(ctx, tx, entity.OrderStatisticRecordEntity{
				Author:      orderInfo.StudentSummary.AuthorId,
				OrgId:       orderInfo.ToOrgID,
				PublisherId: orderInfo.PublisherID,
				OrderSource: orderInfo.OrderSource,
			}, performance)
			if err != nil {
				log.Warning.Printf("Add new performance failed, payment: %#v, performance: %#v, err: %v\n", payment, performance, err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//改价格
func (o *OrderService) UpdateOrderPayPrice(ctx context.Context, orderPayId int, price float64, operator *entity.JWTUser) error {
	//所有payRecord都确认，才能将order确认
	//检查order状态
	//修改支付记录状态
	o.lock.Lock()
	defer o.lock.Unlock()
	log.Info.Printf("Update order pay price, payId: %#v, status: %#v\n", orderPayId, price)
	err := db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		record, err := da.GetOrderModel().GetPayRecordById(ctx, orderPayId)
		if err != nil {
			log.Warning.Printf("Get order pay record failed, orderPayId: %#v, err: %v\n", orderPayId, err)
			return err
		}
		err = da.GetOrderModel().UpdateOrderPayRecordTx(ctx, tx, orderPayId, entity.OrderPayStatusUpdated)
		if err != nil {
			log.Warning.Printf("Update order pay record status failed, orderPayId: %#v, err: %v\n", orderPayId, err)
			return err
		}
		_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, &da.OrderPayRecord{
			OrderID: record.OrderID,
			Mode:    record.Mode,
			Title:   record.Title,
			Amount:  price,
			Content: record.Content,
			Status:  entity.OrderPayStatusPending,
		})
		if err != nil {
			log.Warning.Printf("Update order pay record failed, orderPayId: %#v, err: %v\n", orderPayId, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderService) MarkOrderRemark(ctx context.Context, remarkIDs []int, status int, operator *entity.JWTUser) error {
	if len(remarkIDs) < 1 {
		return nil
	}
	if status != entity.OrderRemarkRead && status != entity.OrderRemarkUnread {
		log.Error.Printf("Invalid remark status, status: %#v, err: %v\n", status, ErrInvalidOrderRemarkStatus)
		return ErrInvalidOrderRemarkStatus
	}
	_, records, err := da.GetOrderModel().SearchRemarkRecord(ctx, da.SearchRemarkRecordCondition{
		RemarkRecordIDList: remarkIDs,
	})
	if err != nil {
		log.Error.Printf("Get order remarks failed, remarkIDs: %#v, err: %v\n", remarkIDs, err)
		return err
	}
	//root只能设置client, client只能设置root
	mode := entity.OrderRemarkModeServer
	if operator.OrgId == entity.RootOrgId {
		mode = entity.OrderRemarkModeClient
	}
	//检查所有records方向是否正确
	for i := range records {
		if records[i].Mode != mode {
			log.Error.Printf("Invalid remark mode, remark: %#v, expect mode:%v err: %v\n", records[i], mode, err)
			return ErrInvalidRemarkID
		}
	}

	err = da.GetOrderModel().UpdateOrderRemarkRecordTx(ctx, db.Get(), remarkIDs, status)
	if err != nil {
		log.Error.Printf("update order remarks status failed, remarkIDs: %#v, status: %v err: %v\n", remarkIDs, status, err)
		return err
	}
	return nil
}

//TODO: search order pay record
func (o *OrderService) AddOrderRemark(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	if content == "" {
		log.Warning.Printf("No remark content, orderID: %v, content: %v\n", orderId, content)
		return ErrNoRemarkContent
	}
	err := db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		err := o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  orderId,
			InfoType: entity.OrderRemarkInfoTypeNormal,
			Info:     "文本消息",
			Content:  content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, orderId: %#v, content: %v err: %v\n", orderId, content, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderService) addOrderRemark(ctx context.Context, tx *gorm.DB, req entity.OrderRemarkRequest, operator *entity.JWTUser) error {
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		log.Warning.Printf("Get order failed, orderId: %#v, err: %v\n", req.OrderID, err)
		return err
	}
	log.Info.Printf("Add order remark, req: %#v\n", req)
	err = o.checkOrderAuthorize(ctx, orderObj, operator)
	if err != nil {
		log.Warning.Printf("Check order authorize failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return err
	}
	mode := entity.OrderRemarkModeClient
	if operator.OrgId == entity.RootOrgId {
		mode = entity.OrderRemarkModeServer
	}

	data := &da.OrderRemarkRecord{
		OrderID:  req.OrderID,
		Author:   operator.UserId,
		Mode:     mode,
		Content:  req.Content,
		Info:     req.Info,
		InfoType: req.InfoType,
		Status:   entity.OrderRemarkUnread,
	}
	_, err = da.GetOrderModel().AddRemarkRecordTx(ctx, tx, data)
	if err != nil {
		log.Warning.Printf("Add remark record failed, order: %#v, data: %#v, operator: %#v, err: %v\n", orderObj, data, operator, err)
		return err
	}
	_, records, err := da.GetOrderModel().SearchRemarkRecord(ctx, da.SearchRemarkRecordCondition{
		OrderIDList: []int{req.OrderID},
	})
	if err != nil {
		log.Error.Printf("search remarks failed, req: %#v, err:%v", req, err)
		return err
	}
	recordsIDs := make([]int, 0)
	for i := range records {
		if records[i].Mode != mode {
			recordsIDs = append(recordsIDs, records[i].ID)
		}
	}
	if len(recordsIDs) > 0 {
		err = da.GetOrderModel().UpdateOrderRemarkRecordTx(ctx, tx, recordsIDs, entity.OrderRemarkRead)
		if err != nil {
			log.Error.Printf("update remarks failed, recordsIDs: %v, err:%v", recordsIDs, err)
			return err
		}
	}

	return nil
}

func (o *OrderService) SearchOrderPayRecords(ctx context.Context, condition *entity.SearchPayRecordCondition, operator *entity.JWTUser) (*entity.PayRecordInfoList, error) {
	log.Info.Printf("Search order pay, condition: %#v\n", condition)
	total, records, err := da.GetOrderModel().SearchPayRecord(ctx, da.SearchPayRecordCondition{
		PayRecordIDList: condition.PayRecordIDList,
		OrderIDList:     condition.OrderIDList,
		AuthorIDList:    condition.AuthorIDList,
		Mode:            condition.Mode,
		StatusList:      condition.StatusList,

		OrderBy: condition.OrderBy,

		PageSize: condition.PageSize,
		Page:     condition.Page,
	})
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	if len(records) < 1 {
		return &entity.PayRecordInfoList{
			Total: 0,
		}, nil
	}
	res, err := o.getPayRecordInfo(ctx, records)
	if err != nil {
		log.Warning.Printf("Get order payment failed, condition: %#v, records: %#v, err: %v\n", condition, records, err)
		return nil, err
	}

	return &entity.PayRecordInfoList{
		Total:   total,
		Records: res,
	}, nil
}

func (o *OrderService) SearchOrderRemarks(ctx context.Context, condition *da.SearchRemarkRecordCondition, operator *entity.JWTUser) (*entity.OrderRemarkList, error) {
	if operator.OrgId == entity.RootOrgId {
		condition.Mode = entity.OrderRemarkModeClient
	} else {
		condition.Mode = entity.OrderRemarkModeServer
	}
	total, records, err := da.GetOrderModel().SearchRemarkRecord(ctx, *condition)
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	res := &entity.OrderRemarkList{
		Total:   total,
		Records: make([]*entity.OrderRemarkRecord, len(records)),
	}

	for i := range records {
		res.Records[i] = &entity.OrderRemarkRecord{
			ID:       records[i].ID,
			OrderID:  records[i].OrderID,
			Author:   records[i].Author,
			Mode:     records[i].Mode,
			Content:  records[i].Content,
			Status:   records[i].Status,
			Info:     records[i].Info,
			InfoType: records[i].InfoType,

			UpdatedAt: records[i].UpdatedAt,
			CreatedAt: records[i].CreatedAt,
		}
	}
	return res, nil
}

func (o *OrderService) ExportOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) ([]byte, error) {
	//查询订单
	log.Info.Printf("Search order, condition: %#v\n", condition)
	if len(condition.ToOrgIDList) > 0 {
		allIds, err := o.getOrgSubOrgs(ctx, condition.ToOrgIDList)
		if err != nil {
			log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}
		condition.ToOrgIDList = allIds
	}

	_, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		OrderIDList:     condition.IDs,
		StudentIDList:   condition.StudentIDList,
		ToOrgIDList:     condition.ToOrgIDList,
		IntentSubjects:  condition.IntentSubjects,
		PublisherIDList: condition.PublisherID,
		OrderSourceList: condition.OrderSourceList,
		CreateStartAt:   condition.CreateStartAt,
		StudentKeywords: condition.StudentsKeywords,
		AuthorIDList:    condition.AuthorID,
		Keywords:        condition.Keywords,
		CreateEndAt:     condition.CreateEndAt,
		UpdateStartAt:   condition.UpdateStartAt,
		UpdateEndAt:     condition.UpdateEndAt,
		Address:         condition.Address,
		Status:          condition.Status,
		OrderBy:         condition.OrderBy,
		Page:            condition.Page,
		PageSize:        condition.PageSize,
	})
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	//添加具体信息
	orderInfos, err := o.getOrderInfoRecords(ctx, orders)
	if err != nil {
		log.Warning.Printf("Get order detailed failed, condition: %#v, orders: %#v, err: %v\n", condition, orders, err)
		return nil, err
	}
	file, err := OrdersToXlsx(ctx, orderInfos)
	if err != nil {
		log.Warning.Printf("Build xlsx file failed, condition: %#v, orderInfos: %#v, err: %v\n", condition, orderInfos, err)
		return nil, err
	}
	buf := &bytes.Buffer{}
	file.Write(buf)
	return buf.Bytes(), nil
}

func (o *OrderService) SearchOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	//查询订单
	log.Info.Printf("Search order, condition: %#v\n", condition)
	if len(condition.ToOrgIDList) > 0 {
		allIds, err := o.getOrgSubOrgs(ctx, condition.ToOrgIDList)
		if err != nil {
			log.Warning.Printf("Search org failed, condition: %#v, err: %v\n", condition, err)
			return nil, err
		}
		condition.ToOrgIDList = allIds
	}

	total, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		OrderIDList:     condition.IDs,
		StudentIDList:   condition.StudentIDList,
		ToOrgIDList:     condition.ToOrgIDList,
		IntentSubjects:  condition.IntentSubjects,
		PublisherIDList: condition.PublisherID,
		OrderSourceList: condition.OrderSourceList,
		CreateStartAt:   condition.CreateStartAt,
		StudentKeywords: condition.StudentsKeywords,
		AuthorIDList:    condition.AuthorID,
		Keywords:        condition.Keywords,
		UpdateStartAt:   condition.UpdateStartAt,
		UpdateEndAt:     condition.UpdateEndAt,
		Address:         condition.Address,
		CreateEndAt:     condition.CreateEndAt,
		Status:          condition.Status,
		OrderBy:         condition.OrderBy,
		Page:            condition.Page,
		PageSize:        condition.PageSize,
	})
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, err: %v\n", condition, err)
		return nil, err
	}
	//添加具体信息
	orderInfos, err := o.getOrderInfoDetails(ctx, orders)
	if err != nil {
		log.Warning.Printf("Get order detailed failed, condition: %#v, orders: %#v, err: %v\n", condition, orders, err)
		return nil, err
	}

	return &entity.OrderInfoList{
		Total:  total,
		Orders: orderInfos,
	}, nil
}

func (o *OrderService) SearchOrderWithAuthor(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	condition.PublisherID = []int{operator.UserId}
	log.Info.Printf("Search order with author, condition: %#v\n", condition)
	//查询订单
	return o.SearchOrders(ctx, condition, operator)
}

func (o *OrderService) SearchOrderWithOrgId(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	orgIds := []int{operator.OrgId}

	condition.ToOrgIDList = orgIds
	log.Info.Printf("Search order with org, condition: %#v\n", condition)
	//查询订单
	return o.SearchOrders(ctx, condition, operator)
}

func (o *OrderService) GetOrderById(ctx context.Context, orderId int, operator *entity.JWTUser) (*entity.OrderInfoWithRecords, error) {
	log.Info.Printf("GetOrderById, orderId: %#v\n", orderId)
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		log.Warning.Printf("Get order failed, orderId: %#v, operator: %#v, err: %v\n", orderId, operator, err)
		return nil, err
	}
	err = o.checkOrderAuthorize(ctx, orderObj, operator)
	if err != nil {
		log.Warning.Printf("Check order authorize failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil, err
	}

	var student *da.Student
	var org *da.Org
	_, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		StudentIDList: []int{orderObj.Order.StudentID},
	})
	if err != nil {
		log.Warning.Printf("Search students failed, orderObj: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil, err
	}
	if len(students) < 1 {
		return nil, ErrInvalidStudentID
	}
	student = students[0]

	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		log.Warning.Printf("Search orgs failed, orderObj: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil, err
	}
	for i := range orgs {
		if orgs[i].ID == orderObj.Order.ToOrgID {
			org = orgs[i]
		}
	}
	if org == nil {
		log.Warning.Printf("Invalid to Org, orgs: %#v, org: %#v, err: %v\n", orgs, org, ErrInvalidToOrgID)
		return nil, ErrInvalidToOrgID
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: []int{orderObj.Order.PublisherID, student.AuthorID},
	})
	if err != nil {
		log.Warning.Printf("Search users failed, orderObj: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil, err
	}
	if len(users) < 1 {
		log.Warning.Printf("Invalid to users, users: %#v, err: %v\n", users, ErrInvalidPublisherID)
		return nil, ErrInvalidPublisherID
	}

	publisherName := ""
	authorName := ""
	for i := range users {
		if users[i].ID == student.AuthorID {
			authorName = users[i].Name
		}
		if users[i].ID == orderObj.Order.PublisherID {
			publisherName = users[i].Name
		}
	}

	orderSourceObj, err := da.GetOrderSourceModel().GetOrderSourceById(ctx, orderObj.Order.OrderSource)
	if err != nil {
		log.Warning.Printf("Get orderSource failed, orderObj: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return nil, err
	}

	res := &entity.OrderInfoWithRecords{
		OrderInfo: entity.OrderInfo{
			ID:            orderObj.Order.ID,
			StudentID:     orderObj.Order.StudentID,
			ToOrgID:       orderObj.Order.ToOrgID,
			IntentSubject: strings.Split(orderObj.Order.IntentSubjects, ","),
			PublisherID:   orderObj.Order.PublisherID,
			Status:        orderObj.Order.Status,
			OrderSource:   orderObj.Order.OrderSource,
			CreatedAt:     orderObj.Order.CreatedAt,
			UpdatedAt:     orderObj.Order.UpdatedAt,
		},
		StudentSummary: &entity.StudentSummaryInfo{
			ID:         student.ID,
			Name:       student.Name,
			Gender:     student.Gender,
			Telephone:  student.Telephone,
			Address:    student.Address,
			AddressExt: student.AddressExt,
			Note:       student.Note,
			AuthorId:   student.AuthorID,
			CreatedAt:  student.CreatedAt,
			UpdatedAt:  student.UpdatedAt,
		},
		OrgName:         org.Name,
		PublisherName:   publisherName,
		AuthorName:      authorName,
		OrderSourceName: orderSourceObj.Name,
	}

	//添加Payment和remark
	payRecords := make([]*entity.OrderPayRecord, len(orderObj.PaymentInfo))
	for i := range orderObj.PaymentInfo {
		payRecords[i] = &entity.OrderPayRecord{
			ID:        orderObj.PaymentInfo[i].ID,
			OrderID:   orderObj.PaymentInfo[i].OrderID,
			Mode:      orderObj.PaymentInfo[i].Mode,
			Amount:    orderObj.PaymentInfo[i].Amount,
			Status:    orderObj.PaymentInfo[i].Status,
			UpdatedAt: orderObj.PaymentInfo[i].UpdatedAt,
			CreatedAt: orderObj.PaymentInfo[i].CreatedAt,
			Title:     orderObj.PaymentInfo[i].Title,
		}
	}
	remarkRecords := make([]*entity.OrderRemarkRecord, len(orderObj.RemarkInfo))
	for i := range orderObj.RemarkInfo {
		remarkRecords[i] = &entity.OrderRemarkRecord{
			ID:        orderObj.RemarkInfo[i].ID,
			OrderID:   orderObj.RemarkInfo[i].OrderID,
			Author:    orderObj.RemarkInfo[i].Author,
			Mode:      orderObj.RemarkInfo[i].Mode,
			Content:   orderObj.RemarkInfo[i].Content,
			UpdatedAt: orderObj.RemarkInfo[i].UpdatedAt,
			CreatedAt: orderObj.RemarkInfo[i].CreatedAt,
			Status:    orderObj.RemarkInfo[i].Status,

			Info:     orderObj.RemarkInfo[i].Info,
			InfoType: orderObj.RemarkInfo[i].InfoType,
		}
	}
	res.PaymentInfo = payRecords
	res.RemarkInfo = remarkRecords

	return res, nil
}

func (o *OrderService) payOrder(ctx context.Context, mode int, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	//检查order状态
	log.Info.Printf("payOrder, req: %#v, mode: %#v\n", req, mode)
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		log.Warning.Printf("Get order failed, req: %#v, err: %v\n", req, err)
		return err
	}
	//检查是否为本机构的订单
	err = o.checkOrderOrg(ctx, operator.OrgId, orderObj.Order.ToOrgID)
	if err != nil {
		log.Warning.Printf("Check order org failed, order: %#v, operator: %#v, err: %v\n", orderObj, operator, err)
		return err
	}

	//增加orderPay
	payData := &da.OrderPayRecord{
		OrderID: req.OrderID,
		Mode:    mode,
		Title:   req.Title,
		Amount:  req.Amount,
		Content: req.Content,
		Status:  entity.OrderPayStatusPending,
	}
	err = db.GetTrans(ctx, func(ctx context.Context, tx *gorm.DB) error {
		_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, payData)
		if err != nil {
			log.Warning.Printf("Add order payment failed, order: %#v, payData: %#v, err: %v\n", orderObj, payData, err)
			return err
		}
		info := fmt.Sprintf("订单支付%v元", req.Amount)
		infoType := entity.OrderRemarkInfoTypePay
		if mode == entity.OrderPayModePayback {
			info = fmt.Sprintf("订单退费%v元", req.Amount)
			infoType = entity.OrderRemarkInfoTypePayback
		}
		err = o.addOrderRemark(ctx, tx, entity.OrderRemarkRequest{
			OrderID:  req.OrderID,
			InfoType: infoType,
			Info:     info,
			Content:  req.Content,
		}, operator)
		if err != nil {
			log.Warning.Printf("AddOrderRemark failed, orderId: %#v, err: %v\n", req, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderService) getPayRecordInfo(ctx context.Context, records []*da.OrderPayRecord) ([]*entity.PayRecordInfo, error) {
	res := make([]*entity.PayRecordInfo, len(records))
	orderIdsList := make([]int, len(records))
	for i := range records {
		orderIdsList[i] = records[i].OrderID
	}
	condition := da.SearchOrderCondition{
		OrderIDList: orderIdsList,
	}
	_, orders, err := da.GetOrderModel().SearchOrder(ctx, condition)
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, records: %#v, err: %v\n", condition, records, err)
		return nil, err
	}
	ordersInfo, err := o.getOrderInfoDetails(ctx, orders)
	if err != nil {
		log.Warning.Printf("Search order failed, condition: %#v, orders: %#v, err: %v\n", condition, orders, err)
		return nil, err
	}

	ordersMap := make(map[int]*entity.OrderInfoDetails)
	for i := range ordersInfo {
		ordersMap[ordersInfo[i].ID] = ordersInfo[i]
	}

	for i := range records {
		order := ordersMap[records[i].OrderID]
		res[i] = &entity.PayRecordInfo{
			ID:        records[i].ID,
			OrderID:   records[i].OrderID,
			Mode:      records[i].Mode,
			Title:     records[i].Title,
			Amount:    records[i].Amount,
			Content:   records[i].Content,
			CreatedAt: records[i].CreatedAt,
			UpdatedAt: records[i].UpdatedAt,

			StudentID:     order.StudentID,
			ToOrgID:       order.ToOrgID,
			IntentSubject: order.IntentSubject,
			PublisherID:   order.PublisherID,
			StudentName:   order.StudentName,
			OrgName:       order.OrgName,
			PublisherName: order.PublisherName,

			Status: records[i].Status,
		}
	}
	return res, nil
}

func (o *OrderService) getOrderInfoDetails(ctx context.Context, orders []*da.Order) ([]*entity.OrderInfoDetails, error) {
	studentIds := make([]int, len(orders))
	userIds := make([]int, len(orders))
	for i := range orders {
		studentIds[i] = orders[i].StudentID
		userIds[i] = orders[i].PublisherID
	}

	_, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		StudentIDList: studentIds,
	})
	if err != nil {
		log.Warning.Printf("Search students failed, students: %#v, orders: %#v, err: %v\n", studentIds, orders, err)
		return nil, err
	}
	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		log.Warning.Printf("Search students failed, orders: %#v, err: %v\n", orders, err)
		return nil, err
	}

	for i := range students {
		userIds = append(userIds, students[i].AuthorID)
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: userIds,
	})
	if err != nil {
		log.Warning.Printf("Search students failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}

	orderSources, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil {
		log.Warning.Printf("Search order source failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}

	studentMaps := make(map[int]*da.Student)
	orgMaps := make(map[int]*da.Org)
	userMaps := make(map[int]*da.User)
	orderSourceMaps := make(map[int]*da.OrderSource)

	for i := range students {
		studentMaps[students[i].ID] = students[i]
	}
	for i := range orgs {
		orgMaps[orgs[i].ID] = orgs[i]
	}
	for i := range users {
		userMaps[users[i].ID] = users[i]
	}
	for i := range orderSources {
		orderSourceMaps[orderSources[i].ID] = orderSources[i]
	}

	orderInfos := make([]*entity.OrderInfoDetails, len(orders))
	for i := range orders {
		orderSourceName := ""
		if orderSourceMaps[orders[i].OrderSource] != nil {
			orderSourceName = orderSourceMaps[orders[i].OrderSource].Name
		}
		orderInfos[i] = &entity.OrderInfoDetails{
			OrderInfo: entity.OrderInfo{
				ID:            orders[i].ID,
				StudentID:     orders[i].StudentID,
				ToOrgID:       orders[i].ToOrgID,
				IntentSubject: strings.Split(orders[i].IntentSubjects, ","),
				PublisherID:   orders[i].PublisherID,
				Status:        orders[i].Status,
				OrderSource:   orders[i].OrderSource,
				CreatedAt:     orders[i].CreatedAt,
				UpdatedAt:     orders[i].UpdatedAt,
			},
			StudentName:      studentMaps[orders[i].StudentID].Name,
			StudentTelephone: studentMaps[orders[i].StudentID].Telephone,
			OrgName:          orgMaps[orders[i].ToOrgID].Name,
			PublisherName:    userMaps[orders[i].PublisherID].Name,
			OrderSourceName:  orderSourceName,
			AuthorID:         studentMaps[orders[i].StudentID].AuthorID,
			AuthorName:       userMaps[studentMaps[orders[i].StudentID].AuthorID].Name,
		}
	}
	return orderInfos, nil
}

func (o *OrderService) getOrderInfoRecords(ctx context.Context, orders []*da.Order) ([]*entity.OrderInfoWithRecords, error) {
	studentIds := make([]int, len(orders))
	userIds := make([]int, len(orders))
	orderIds := make([]int, len(orders))
	for i := range orders {
		studentIds[i] = orders[i].StudentID
		userIds[i] = orders[i].PublisherID
		orderIds[i] = orders[i].ID
	}

	_, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		StudentIDList: studentIds,
	})
	if err != nil {
		log.Warning.Printf("Search students failed, students: %#v, orders: %#v, err: %v\n", studentIds, orders, err)
		return nil, err
	}
	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		log.Warning.Printf("Search students failed, orders: %#v, err: %v\n", orders, err)
		return nil, err
	}

	_, users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: userIds,
	})
	if err != nil {
		log.Warning.Printf("Search students failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}

	orderSources, err := da.GetOrderSourceModel().ListOrderSources(ctx)
	if err != nil {
		log.Warning.Printf("Search order source failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}

	//collect payment data
	paymentInfoMaps := make(map[int][]*entity.OrderPayRecord)
	_, payRecords, err := da.GetOrderModel().SearchPayRecord(ctx, da.SearchPayRecordCondition{
		OrderIDList: orderIds,
	})
	if err != nil {
		log.Warning.Printf("Search pay records failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}
	for i := range payRecords {
		payments := paymentInfoMaps[payRecords[i].OrderID]
		payments = append(payments, &entity.OrderPayRecord{
			ID:        payRecords[i].ID,
			OrderID:   payRecords[i].OrderID,
			Mode:      payRecords[i].Mode,
			Amount:    payRecords[i].Amount,
			Content:   payRecords[i].Content,
			Status:    payRecords[i].Status,
			UpdatedAt: payRecords[i].UpdatedAt,
			CreatedAt: payRecords[i].CreatedAt,
			Title:     payRecords[i].Title,
		})
		paymentInfoMaps[payRecords[i].OrderID] = payments
	}

	//collect remark data
	remarkInfoMaps := make(map[int][]*entity.OrderRemarkRecord)
	_, remarkRecords, err := da.GetOrderModel().SearchRemarkRecord(ctx, da.SearchRemarkRecordCondition{
		OrderIDList: orderIds,
	})
	if err != nil {
		log.Warning.Printf("Search remark records failed, userIds: %#v, orders: %#v, err: %v\n", userIds, orders, err)
		return nil, err
	}
	for i := range remarkRecords {
		remarks := remarkInfoMaps[remarkRecords[i].OrderID]
		remarks = append(remarks, &entity.OrderRemarkRecord{
			ID:        remarkRecords[i].ID,
			OrderID:   remarkRecords[i].OrderID,
			Author:    remarkRecords[i].Author,
			Mode:      remarkRecords[i].Mode,
			Content:   remarkRecords[i].Content,
			Status:    remarkRecords[i].Status,
			Info:      remarkRecords[i].Info,
			InfoType:  remarkRecords[i].InfoType,
			UpdatedAt: remarkRecords[i].UpdatedAt,
			CreatedAt: remarkRecords[i].CreatedAt,
		})
		remarkInfoMaps[remarkRecords[i].OrderID] = remarks
	}

	studentMaps := make(map[int]*da.Student)
	orgMaps := make(map[int]*da.Org)
	userMaps := make(map[int]*da.User)
	orderSourceMaps := make(map[int]*da.OrderSource)

	for i := range students {
		studentMaps[students[i].ID] = students[i]
	}
	for i := range orgs {
		orgMaps[orgs[i].ID] = orgs[i]
	}
	for i := range users {
		userMaps[users[i].ID] = users[i]
	}
	for i := range orderSources {
		orderSourceMaps[orderSources[i].ID] = orderSources[i]
	}

	orderInfos := make([]*entity.OrderInfoWithRecords, len(orders))
	for i := range orders {
		orderSourceName := ""
		if orderSourceMaps[orders[i].OrderSource] != nil {
			orderSourceName = orderSourceMaps[orders[i].OrderSource].Name
		}
		orderInfos[i] = &entity.OrderInfoWithRecords{
			OrderInfo: entity.OrderInfo{
				ID:            orders[i].ID,
				StudentID:     orders[i].StudentID,
				ToOrgID:       orders[i].ToOrgID,
				IntentSubject: strings.Split(orders[i].IntentSubjects, ","),
				PublisherID:   orders[i].PublisherID,
				Status:        orders[i].Status,
				OrderSource:   orders[i].OrderSource,
				CreatedAt:     orders[i].CreatedAt,
				UpdatedAt:     orders[i].UpdatedAt,
			},
			StudentSummary: &entity.StudentSummaryInfo{
				ID:         studentMaps[orders[i].StudentID].ID,
				Name:       studentMaps[orders[i].StudentID].Name,
				Gender:     studentMaps[orders[i].StudentID].Gender,
				Telephone:  studentMaps[orders[i].StudentID].Telephone,
				Address:    studentMaps[orders[i].StudentID].Address,
				AddressExt: studentMaps[orders[i].StudentID].AddressExt,
				Note:       studentMaps[orders[i].StudentID].Note,
				AuthorId:   studentMaps[orders[i].StudentID].AuthorID,
				CreatedAt:  studentMaps[orders[i].StudentID].CreatedAt,
				UpdatedAt:  studentMaps[orders[i].StudentID].UpdatedAt,
			},
			RemarkInfo:      remarkInfoMaps[orders[i].ID],
			PaymentInfo:     paymentInfoMaps[orders[i].ID],
			AuthorName:      userMaps[studentMaps[orders[i].StudentID].AuthorID].Name,
			OrgName:         orgMaps[orders[i].ToOrgID].Name,
			PublisherName:   userMaps[orders[i].PublisherID].Name,
			OrderSourceName: orderSourceName,
		}
	}
	return orderInfos, nil
}

func (o *OrderService) checkEntity(ctx context.Context, org *da.Org, orderEntity orderEntity) error {
	_, err := da.GetStudentModel().GetStudentById(ctx, orderEntity.StudentID)
	if err != nil {
		log.Warning.Println("Can't find student when check entity, orderEntity:", orderEntity)
		return err
	}

	if org.Status != entity.OrgStatusCertified {
		return ErrInvalidOrgStatus
	}

	return nil
}

func (o *OrderService) checkOrderAuthorize(ctx context.Context, order *da.OrderInfo, operator *entity.JWTUser) error {
	if operator.OrgId == entity.RootOrgId {
		return nil
	}
	subOrgs, err := da.GetOrgModel().GetOrgsByParentId(ctx, operator.OrgId)
	if err != nil {
		return err
	}

	if operator.OrgId != order.Order.ToOrgID {
		flag := false
		for i := range subOrgs {
			if subOrgs[i].ID == order.Order.ToOrgID {
				flag = true
				break
			}
		}

		if !flag {
			log.Warning.Printf("checkOrderAuthorize failed, order: %#v, operator: %#v, err: %v\n", order, operator, ErrNoAuthorizeToOperate)
			return ErrNoAuthorizeToOperate
		}
	}

	return nil
}

func (o *OrderService) checkOrderOrg(ctx context.Context, orgId, toOrgId int) error {
	if toOrgId != orgId {
		orgs, err := GetOrgService().GetSubOrgs(ctx, orgId)
		if err != nil {
			log.Warning.Printf("Get sub orgs failed, orgId: %#v, toOrgId: %#v, err: %v\n", orgId, toOrgId, err)
			return err
		}
		flag := false
		for i := range orgs {
			if orgs[i].ID == toOrgId {
				flag = true
				break
			}
		}
		if !flag {
			log.Warning.Printf("Get sub orgs failed, orgs: %#v, toOrgId: %#v, err: %v\n", orgs, toOrgId, ErrNoAuthorizeToOperate)
			return ErrNoAuthorizeToOperate
		}
	}
	return nil
}

func (o *OrderService) getOrgSubOrgs(ctx context.Context, orgIds []int) ([]int, error) {
	orgs, err := da.GetOrgModel().ListOrgsByIDs(ctx, orgIds)
	if err != nil {
		return nil, err
	}
	ret := make([]int, 0)
	parentIds := make([]int, len(orgs))
	for i := range orgs {
		if orgs[i].ParentID == 0 {
			parentIds[i] = orgs[i].ID
		}
	}
	if len(parentIds) > 0 {
		_, subOrgs, err := da.GetOrgModel().SearchOrgs(ctx, da.SearchOrgsCondition{
			Status:    []int{entity.OrgStatusCertified},
			ParentIDs: parentIds,
		})
		if err != nil {
			return nil, err
		}
		for i := range subOrgs {
			ret = append(ret, subOrgs[i].ID)
		}
	}
	ret = append(ret, orgIds...)
	return ret, nil
}

var (
	_orderService     IOrderService
	_orderServiceOnce sync.Once
)

func GetOrderService() IOrderService {
	_orderServiceOnce.Do(func() {
		if _orderService == nil {
			_orderService = new(OrderService)
		}
	})
	return _orderService
}
