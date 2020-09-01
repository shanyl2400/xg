package route

import (
	"net/http"
	"strconv"
	"strings"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) createOrder(c *gin.Context) {
	req := new(entity.CreateOrderRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	id, err := service.GetOrderService().CreateOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}

func (s *Server) searchOrder(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrders(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", orders)
}

func (s *Server) searchOrderWithAuthor(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrderWithAuthor(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", orders)
}

func (s *Server) searchOrderWithOrgID(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrderWithOrgId(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", orders)
}

func (s *Server) getOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)

	order, err := service.GetOrderService().GetOrderById(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", order)
}

func (s *Server) searchPendingPayRecord(c *gin.Context) {
	req := buildSearchPayRecordCondition(c)
	user := s.getJWTUser(c)
	req.StatusList = []int{entity.OrderPayStatusPending}

	records, err := service.GetOrderService().SearchOrderPayRecords(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", records)
}

func (s *Server) signupOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().SignUpOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) revokeOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().RevokeOrder(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}


func (s *Server) payOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().PayOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) paybackOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().PaybackOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) acceptPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().ConfirmOrderPay(c.Request.Context(), id, entity.OrderPayStatusChecked , user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) rejectPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().ConfirmOrderPay(c.Request.Context(), id, entity.OrderPayStatusRejected , user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) addOrderMark(c *gin.Context) {
	req := new(entity.OrderMarkRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().AddOrderRemark(c.Request.Context(), req.OrderID, req.Content, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func buildSearchPayRecordCondition(c *gin.Context) *entity.SearchPayRecordCondition {
	payRecordIds := c.Query("pay_record_ids")
	orderIds := c.Query("order_ids")
	authorIds := c.Query("author_ids")
	mode := c.Query("mode")
	status := c.Query("status")
	orderBy := c.Query("order_by")
	page := c.Query("page")
	pageSize := c.Query("page_size")

	return &entity.SearchPayRecordCondition{
		PayRecordIDList: parseInts(payRecordIds),
		OrderIDList:     parseInts(orderIds),
		AuthorIDList:    parseInts(authorIds),
		Mode:            parseInt(mode),
		StatusList:      parseInts(status),

		OrderBy: orderBy,

		PageSize: parseInt(pageSize),
		Page:     parseInt(page),
	}
}

func buildOrderCondition(c *gin.Context) *entity.SearchOrderCondition {
	studentIds := c.Query("student_ids")
	toOrgIds := c.Query("to_org_ids")
	intentSubjects := c.Query("intent_subjects")
	publisherID := c.Query("publisher_id")
	status := c.Query("status")
	orderBy := c.Query("order_by")
	page := c.Query("page")
	pageSize := c.Query("page_size")

	return &entity.SearchOrderCondition{
		StudentIDList:  parseInts(studentIds),
		ToOrgIDList:    parseInts(toOrgIds),
		IntentSubjects: intentSubjects,
		PublisherID:    parseInt(publisherID),

		Status:  parseInt(status),
		OrderBy: orderBy,

		PageSize: parseInt(pageSize),
		Page:     parseInt(page),
	}
}

func parseInt(str string) int {
	id, err := strconv.Atoi(str)
	if err == nil {
		return 0
	}
	return id
}
func parseInts(str string) []int {
	strList := strings.Split(str, ",")
	ret := make([]int, 0)
	for i := range strList {
		id, err := strconv.Atoi(strList[i])
		if err == nil {
			ret = append(ret, id)
		}
	}
	if len(ret) < 1 {
		return nil
	}
	return ret
}
