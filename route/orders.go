package route

import (
	"net/http"
	"strconv"
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
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	id, err := service.GetOrderService().CreateOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}

func (s *Server) searchOrder(c *gin.Context) {
	req := new(entity.SearchOrderCondition)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}

	orders, err := service.GetOrderService().SearchOrders(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", orders)
}

func (s *Server) searchOrderWithAuthor(c *gin.Context) {
	req := new(entity.SearchOrderCondition)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}

	orders, err := service.GetOrderService().SearchOrderWithAuthor(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", orders)
}

func (s *Server) searchOrderWithOrgID(c *gin.Context) {
	req := new(entity.SearchOrderCondition)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}

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
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}

	order, err := service.GetOrderService().GetOrderById(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", order)
}

func (s *Server) searchPendingPayRecord(c *gin.Context) {
	req := new(entity.SearchPayRecordCondition)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	req.StatusList = []int{entity.OrderPayStatusPendingCheck}

	records, err := service.GetOrderService().SearchOrderPayRecords(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", records)
}
