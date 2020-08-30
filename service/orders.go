package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"xg/da"
	"xg/db"
	"xg/entity"
)

type OrderService struct {
}
type orderEntity struct {
	StudentID int `json:"student_id"`
	ToOrgID   int `json:"to_org_id"`
}

func (o *OrderService) checkEntity(ctx context.Context, orderEntity orderEntity) error {
	_, err := da.GetStudentModel().GetStudentById(ctx, orderEntity.StudentID)
	if err != nil {
		fmt.Println("Can't find student when check entity")
		return err
	}

	_, err = da.GetOrgModel().GetOrgById(ctx, orderEntity.ToOrgID)
	if err != nil {
		fmt.Println("Can't find org when check entity")
		return err
	}

	return nil
}
func (o *OrderService) CreateOrder(ctx context.Context, req *entity.CreateOrderRequest, operator *entity.JWTUser) (int, error) {
	//Check entity
	err := o.checkEntity(ctx, orderEntity{
		StudentID: req.StudentID,
		ToOrgID:   req.ToOrgID,
	})
	if err != nil {
		return -1, err
	}

	//TODO:检查重复订单？

	id, err := da.GetOrderModel().CreateOrder(ctx, da.Order{
		StudentID:      req.StudentID,
		ToOrgID:        req.ToOrgID,
		IntentSubjects: strings.Join(req.IntentSubjects, ","),
		PublisherID:    operator.UserId,
		Status:         entity.OrderStatusCreated,
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

//报名
func (o *OrderService) SignUpOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	//检查order状态
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		return err
	}
	if orderObj.Order.ToOrgID != operator.OrgId {
		return ErrNoAuthorizeToOperate
	}
	if orderObj.Order.Status != entity.OrderStatusCreated {
		return ErrNoAuthToOperateOrder
	}
	tx := db.Get().Begin()
	//修改order状态
	err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, req.OrderID, entity.OrderStatusPendingSigned)
	if err != nil {
		tx.Rollback()
		return err
	}
	//增加orderPay
	_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, &da.OrderPayRecord{
		OrderID: req.OrderID,
		Mode:    entity.OrderPayModePay,
		Title:   req.Title,
		Amount:  req.Amount,
		Status:  entity.OrderPayStatusPendingCheck,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (o *OrderService) payOrder(ctx context.Context, mode int, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	//检查order状态
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, req.OrderID)
	if err != nil {
		return err
	}
	if orderObj.Order.ToOrgID != operator.OrgId {
		return ErrNoAuthorizeToOperate
	}
	if orderObj.Order.Status == entity.OrderStatusCreated {
		return ErrNoAuthToOperateOrder
	}

	tx := db.Get().Begin()
	//若order不为提醒状态，则修改order状态
	if orderObj.Order.Status != entity.OrderStatusPendingSigned && orderObj.Order.Status != entity.OrderStatusPendingCheck {
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, req.OrderID, entity.OrderStatusPendingCheck)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//增加orderPay
	_, err = da.GetOrderModel().AddOrderPayRecordTx(ctx, tx, &da.OrderPayRecord{
		OrderID: req.OrderID,
		Mode:    mode,
		Title:   req.Title,
		Amount:  req.Amount,
		Status:  entity.OrderPayStatusPendingCheck,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//付款
func (o *OrderService) PayOrder(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	return o.payOrder(ctx, entity.OrderPayModePay, req, operator)
}

//退费
func (o *OrderService) Payback(ctx context.Context, req *entity.OrderPayRequest, operator *entity.JWTUser) error {
	return o.payOrder(ctx, entity.OrderPayModePayback, req, operator)
}

//确认
func (o *OrderService) ConfirmOrderPay(ctx context.Context, orderId, orderPayId int, operator *entity.JWTUser) error {
	//所有payRecord都确认，才能将order确认
	//检查order状态
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	if orderObj.Order.Status == entity.OrderStatusCreated ||
		orderObj.Order.Status == entity.OrderStatusChecked {
		return ErrNoNeedToOperate
	}

	tx := db.Get().Begin()
	//修改支付记录状态
	err = da.GetOrderModel().UpdateOrderPayRecordTx(ctx, tx, orderPayId, entity.OrderPayStatusChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	//若没有其他支付，则修改订单状态
	hasPending := false
	for i := range orderObj.PaymentInfo {
		if orderObj.PaymentInfo[i].ID != orderPayId &&
			orderObj.PaymentInfo[i].Status == entity.OrderPayStatusPendingCheck {
			hasPending = true
		}
	}
	//当前订单已无pending支付记录,将订单设置为checked
	if !hasPending {
		err = da.GetOrderModel().UpdateOrderStatusTx(ctx, tx, orderId, entity.OrderStatusChecked)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

//TODO: search order pay record

func (o *OrderService) AddOrderRemark(ctx context.Context, orderId int, content string, operator *entity.JWTUser) error {
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		return err
	}
	err = o.checkOrderAuthorize(ctx, orderObj, operator)
	if err != nil {
		return err
	}
	mode := entity.OrderRemarkModeClient
	if operator.OrgId == entity.RootOrgId {
		mode = entity.OrderRemarkModeServer
	}
	_, err = da.GetOrderModel().AddRemarkRecord(ctx, &da.OrderRemarkRecord{
		OrderID: orderId,
		Author:  operator.UserId,
		Mode:    mode,
		Content: content,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderService) SearchOrderPayRecords(ctx context.Context, condition *entity.SearchPayRecordCondition, operator *entity.JWTUser) (*entity.PayRecordInfoList, error) {
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
		return nil, err
	}

	return &entity.PayRecordInfoList{
		Total:   total,
		Records: o.getPayRecordInfo(ctx, records),
	}, nil
}

func (o *OrderService) SearchOrders(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	//查询订单
	total, orders, err := da.GetOrderModel().SearchOrder(ctx, da.SearchOrderCondition{
		StudentIDList:  condition.StudentIDList,
		ToOrgIDList:    condition.ToOrgIDList,
		IntentSubjects: condition.IntentSubjects,
		PublisherID:    condition.PublisherID,
		Status:         condition.Status,
		OrderBy:        condition.OrderBy,
		Page:           condition.Page,
		PageSize:       condition.PageSize,
	})
	if err != nil {
		return nil, err
	}
	//添加具体信息
	orderInfos, err := o.getOrderInfoDetails(ctx, orders)
	if err != nil {
		return nil, err
	}

	return &entity.OrderInfoList{
		Total:  total,
		Orders: orderInfos,
	}, nil
}

func (o *OrderService) SearchOrderWithAuthor(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	condition.PublisherID = operator.UserId
	//查询订单
	return o.SearchOrders(ctx, condition, operator)
}

func (o *OrderService) SearchOrderWithOrgId(ctx context.Context, condition *entity.SearchOrderCondition, operator *entity.JWTUser) (*entity.OrderInfoList, error) {
	condition.ToOrgIDList = []int{operator.OrgId}
	//查询订单
	return o.SearchOrders(ctx, condition, operator)
}

func (o *OrderService) getPayRecordInfo(ctx context.Context, records []*da.OrderPayRecord) []*entity.PayRecordInfo {
	res := make([]*entity.PayRecordInfo, len(records))
	for i := range records {
		res[i] = &entity.PayRecordInfo{
			ID:      records[i].ID,
			OrderID: records[i].OrderID,
			Mode:    records[i].Mode,
			Title:   records[i].Title,
			Amount:  records[i].Amount,

			Status: records[i].Status,
		}
	}
	return res
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
		return nil, err
	}
	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		return nil, err
	}

	users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: userIds,
	})
	if err != nil {
		return nil, err
	}

	studentMaps := make(map[int]*da.Student)
	orgMaps := make(map[int]*da.Org)
	userMaps := make(map[int]*da.User)

	for i := range students {
		studentMaps[students[i].ID] = students[i]
	}
	for i := range orgs {
		orgMaps[orgs[i].ID] = orgs[i]
	}
	for i := range users {
		userMaps[users[i].ID] = users[i]
	}

	orderInfos := make([]*entity.OrderInfoDetails, len(orders))
	for i := range orders {
		orderInfos[i] = &entity.OrderInfoDetails{
			OrderInfo: entity.OrderInfo{
				ID:            orders[i].ID,
				StudentID:     orders[i].StudentID,
				ToOrgID:       orders[i].ToOrgID,
				IntentSubject: strings.Split(orders[i].IntentSubjects, ","),
				PublisherID:   orders[i].PublisherID,
				Status:        orders[i].Status,
			},
			StudentName:      studentMaps[orders[i].StudentID].Name,
			StudentTelephone: studentMaps[orders[i].StudentID].Telephone,
			OrgName:          orgMaps[orders[i].ToOrgID].Name,
			PublisherName:    userMaps[orders[i].PublisherID].Name,
		}
	}
	return orderInfos, nil
}

func (o *OrderService) GetOrderById(ctx context.Context, orderId int, operator *entity.JWTUser) (*entity.OrderInfoWithRecords, error) {
	orderObj, err := da.GetOrderModel().GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}
	err = o.checkOrderAuthorize(ctx, orderObj, operator)
	if err != nil {
		return nil, err
	}

	var student *da.Student
	var org *da.Org
	var user *da.User
	_, students, err := da.GetStudentModel().SearchStudents(ctx, da.SearchStudentCondition{
		StudentIDList: []int{orderObj.Order.StudentID},
	})
	if err != nil {
		return nil, err
	}
	if len(students) < 1 {
		return nil, ErrInvalidStudentID
	}
	student = students[0]

	orgs, err := da.GetOrgModel().ListOrgs(ctx)
	if err != nil {
		return nil, err
	}
	for i := range orgs {
		if orgs[i].ID == orderObj.Order.ToOrgID {
			org = orgs[i]
		}
	}
	if org == nil {
		return nil, ErrInvalidToOrgID
	}

	users, err := da.GetUsersModel().SearchUsers(ctx, da.SearchUserCondition{
		IDList: []int{orderObj.Order.PublisherID},
	})
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return nil, ErrInvalidPublisherID
	}
	user = users[0]

	res := &entity.OrderInfoWithRecords{
		OrderInfo: entity.OrderInfo{
			ID:            orderObj.Order.ID,
			StudentID:     orderObj.Order.StudentID,
			ToOrgID:       orderObj.Order.ToOrgID,
			IntentSubject: strings.Split(orderObj.Order.IntentSubjects, ","),
			PublisherID:   orderObj.Order.PublisherID,
			Status:        orderObj.Order.Status,
		},
		StudentSummary: &entity.StudentSummaryInfo{
			ID:        student.ID,
			Name:      student.Name,
			Gender:    student.Gender,
			Telephone: student.Telephone,
			Address:   student.Address,
			Note:      student.Note,
		},
		OrgName:       org.Name,
		PublisherName: user.Name,
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
		}
	}
	res.PaymentInfo = payRecords
	res.RemarkInfo = remarkRecords

	return res, nil
}

func (o *OrderService) checkOrderAuthorize(ctx context.Context, order *da.OrderInfo, operator *entity.JWTUser) error {
	if operator.OrgId == entity.RootOrgId {
		return nil
	}
	if operator.OrgId != order.Order.ToOrgID {
		return ErrNoAuthorizeToOperate
	}
	return nil
}

var (
	_orderService     *OrderService
	_orderServiceOnce sync.Once
)

func GetOrderService() *OrderService {
	_orderServiceOnce.Do(func() {
		if _orderService == nil {
			_orderService = new(OrderService)
		}
	})
	return _orderService
}
