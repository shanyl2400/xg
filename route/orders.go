package route

import (
	"net/http"
	"strconv"
	"time"
	"xg/da"
	"xg/db"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary createOrder
// @Description create an order
// @Accept json
// @Produce json
// @Param request body entity.CreateOrderRequest true "create request"
// @Param Authorization header string true "With the bearer"
// @Tags order
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order [post]
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
	c.JSON(http.StatusOK, IdResponse{
		ID:     id,
		ErrMsg: "success",
	})
}

// @Summary searchOrder
// @Description search order with condition
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param student_ids query string false "search order with student_ids"
// @Param to_org_ids query string false "search order with to_org_ids"
// @Param students_keywords query string false "search order by students info"
// @Param intent_subjects query string false "search order with intent_subjects"
// @Param publisher_id query int  false "search order with publisher_id"
// @Param status query string  false "search order with status"
// @Param order_by query string false "search order order by column name"
// @Param page_size query int true "order list page size"
// @Param page query int false "order list page index"
// @Tags order
// @Success 200 {object} entity.OrderInfoList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders [get]
func (s *Server) searchOrder(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrders(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderInfoListResponse{
		Data:   orders,
		ErrMsg: "success",
	})
}

// @Summary searchOrderWithAuthor
// @Description search order in author with condition
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param student_ids query string false "search order with student_ids"
// @Param to_org_ids query string false "search order with to_org_ids"
// @Param students_keywords query string false "search order by students info"
// @Param intent_subjects query string false "search order with intent_subjects"
// @Param status query string  false "search order with status"
// @Param order_by query string false "search order order by column name"
// @Param page_size query int true "order list page size"
// @Param page query int false "order list page index"
// @Tags order
// @Success 200 {object} entity.OrderInfoList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders/author [get]
func (s *Server) searchOrderWithAuthor(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrderWithAuthor(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderInfoListResponse{
		Data:   orders,
		ErrMsg: "success",
	})
}

// @Summary searchOrderWithOrgID
// @Description search order in org with condition
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param student_ids query string false "search order with student_ids"
// @Param intent_subjects query string false "search order with intent_subjects"
// @Param students_keywords query string false "search order by students info"
// @Param status query string  false "search order with status"
// @Param publisher_id query int  false "search order with publisher_id"
// @Param order_by query string false "search order order by column name"
// @Param page_size query int true "order list page size"
// @Param page query int false "order list page index"
// @Tags order
// @Success 200 {object} entity.OrderInfoList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders/org [get]
func (s *Server) searchOrderWithOrgID(c *gin.Context) {
	req := buildOrderCondition(c)
	user := s.getJWTUser(c)

	orders, err := service.GetOrderService().SearchOrderWithOrgId(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderInfoListResponse{
		Data:   orders,
		ErrMsg: "success",
	})
}

// @Summary getOrderByID
// @Description get order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Tags order
// @Success 200 {object} entity.OrderInfoWithRecords
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id} [get]
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
	c.JSON(http.StatusOK, OrderRecordResponse{
		Data:   order,
		ErrMsg: "success",
	})
}

// @Summary searchPendingPayRecord
// @Description search pending order pay record with condition
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param pay_record_ids query string false "search order pay record with pay_record_ids"
// @Param order_ids query string false "search order pay record with order_ids"
// @Param author_ids query string false "search order pay record with author_ids"
// @Param mode query int false "search order pay record with mode"
// @Param status query string  false "search order pay record with status"
// @Param order_by query string false "search order pay record by column name"
// @Param page_size query int true "order pay record list page size"
// @Param page query int true "order pay record list page index"
// @Tags order
// @Success 200 {object} entity.OrderRemarkList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders/pending [get]
func (s *Server) searchPendingPayRecord(c *gin.Context) {
	req := buildSearchPayRecordCondition(c)
	user := s.getJWTUser(c)
	req.StatusList = []int{entity.OrderPayStatusPending}

	records, err := service.GetOrderService().SearchOrderPayRecords(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderPaymentRecordListResponse{
		Data:   records,
		ErrMsg: "success",
	})
}

// @Summary searchOrderRemarks
// @Description search order remarks
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param author_ids query string false "search order pay record with author_ids"
// @Param mode query int false "search order pay record with mode"
// @Param status query string  false "search order pay record with status"
// @Param order_by query string false "search order pay record by column name"
// @Param page_size query int true "order pay record list page size"
// @Param page query int true "order pay record list page index"
// @Tags order
// @Success 200 {object} entity.OrderRemarkList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders/remarks [get]
func (s *Server) searchOrderRemarks(c *gin.Context) {
	req := buildOrderRemarkCondition(c)
	user := s.getJWTUser(c)

	records, err := service.GetOrderService().SearchOrderRemarks(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, records)
}

// @Summary signupOrder
// @Description signup order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Param request body entity.OrderPayRequest true "order signup request"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/signup [put]
func (s *Server) signupOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	req.OrderID = id
	user := s.getJWTUser(c)
	err = service.GetOrderService().SignUpOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary depositOrder
// @Description deposit order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Param request body entity.OrderPayRequest true "order signup request"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/deposit [put]
func (s *Server) depositOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	req.OrderID = id
	user := s.getJWTUser(c)
	err = service.GetOrderService().DepositOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary revokeOrder
// @Description revoke order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/revoke [put]
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

// @Summary invalidOrder
// @Description invalid order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/invalid [put]
func (s *Server) invalidOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().InvalidOrder(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary payOrder
// @Description pay order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Param request body entity.OrderPayRequest true "order pay request"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/pay [put]
func (s *Server) payOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	req.OrderID = id
	user := s.getJWTUser(c)
	err = service.GetOrderService().PayOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary paybackOrder
// @Description payback order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Param request body entity.OrderPayRequest true "order payback request"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/payback [put]
func (s *Server) paybackOrder(c *gin.Context) {
	req := new(entity.OrderPayRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	req.OrderID = id
	user := s.getJWTUser(c)
	err = service.GetOrderService().PaybackOrder(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary acceptPayment
// @Description accept payment by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "payment id"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/payment/{id}/review/accept [put]
func (s *Server) acceptPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().ConfirmOrderPay(c.Request.Context(), id, entity.OrderPayStatusChecked, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary rejectPayment
// @Description reject payment by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "payment id"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/payment/{id}/review/reject [put]
func (s *Server) rejectPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrderService().ConfirmOrderPay(c.Request.Context(), id, entity.OrderPayStatusRejected, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary addOrderMark
// @Description add mark to order by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Param request body entity.OrderMarkRequest true "create mark"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order/{id}/mark [post]
func (s *Server) addOrderMark(c *gin.Context) {
	req := new(entity.OrderMarkRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	req.OrderID = id
	user := s.getJWTUser(c)
	err = service.GetOrderService().AddOrderRemark(c.Request.Context(), req.OrderID, req.Content, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary searchAuthorNotifies
// @Description search author notifies
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param classify query string false "search order notify with classify"
// @Param status query string  false "search order notify with status"
// @Param order_by query string false "search order notify by column name"
// @Param page_size query int true "order notify list page size"
// @Param page query int true " order notify list page index"
// @Tags order
// @Success 200 {object} entity.OrderRemarkList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/notifies/orders/author [get]
func (s *Server) searchAuthorNotifies(c *gin.Context) {
	req := buildOrderNotifyCondition(c)
	user := s.getJWTUser(c)

	total, records, err := service.GetOrderNotifyService().SearchAuthorNotifies(c.Request.Context(), *req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderNotifyResponse{
		Total:  total,
		Data:   records,
		ErrMsg: "success",
	})
}

// @Summary searchAuthorNotifies
// @Description search author notifies
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param classify query string false "search order notify with classify"
// @Param status query string  false "search order notify with status"
// @Param order_by query string false "search order notify by column name"
// @Param page_size query int true "order notify list page size"
// @Param page query int true " order notify list page index"
// @Tags order
// @Success 200 {object} entity.OrderRemarkList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/notifies/orders [get]
func (s *Server) searchNotifies(c *gin.Context) {
	req := buildOrderNotifyCondition(c)
	user := s.getJWTUser(c)

	total, records, err := service.GetOrderNotifyService().SearchNotifies(c.Request.Context(), *req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrderNotifyResponse{
		Total:  total,
		Data:   records,
		ErrMsg: "success",
	})
}

// @Summary markNotify
// @Description mark notify remarks
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order id"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/notify/orders/{id} [put]
func (s *Server) markOrderNotify(c *gin.Context) {
	user := s.getJWTUser(c)
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	err := service.GetOrderNotifyService().MarkNotifyRead(c.Request.Context(), db.Get(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	s.responseSuccess(c)
}

// @Summary markOrderRemark
// @Description mark order remarks
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.MarkOrderRemarkRequest true "mark order remarks requests"
// @Tags order
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orders/marks [put]
func (s *Server) markOrderRemarks(c *gin.Context) {
	req := new(entity.MarkOrderRemarkRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	user := s.getJWTUser(c)
	err = service.GetOrderService().MarkOrderRemark(c.Request.Context(), req.IDs, req.Status, user)
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

func buildOrderNotifyCondition(c *gin.Context) *da.OrderNotifiesCondition {
	classifies := c.Query("classify")
	status := c.Query("status")
	condition := &da.OrderNotifiesCondition{
		Classifies: parseInts(classifies),
		Status:     parseInts(status),
	}
	return condition
}

func buildOrderRemarkCondition(c *gin.Context) *da.SearchRemarkRecordCondition {
	status := c.Query("status")
	orderIds := c.Query("order_ids")
	authorIds := c.Query("author_ids")
	orderBy := c.Query("order_by")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	condition := &da.SearchRemarkRecordCondition{
		OrderIDList:  parseInts(orderIds),
		AuthorIDList: parseInts(authorIds),
		StatusList:   parseInts(status),
		OrderBy:      orderBy,

		PageSize: parseInt(pageSize),
		Page:     parseInt(page),
	}
	return condition
}

func buildOrderCondition(c *gin.Context) *entity.SearchOrderCondition {
	studentIds := c.Query("student_ids")
	toOrgIds := c.Query("to_org_ids")
	intentSubjects := c.Query("intent_subjects")
	publisherID := c.Query("publisher_id")
	status := c.Query("status")
	orderSources := c.Query("order_sources")
	orderBy := c.Query("order_by")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	studentsKeywords := c.Query("students_keywords")
	keywords := c.Query("keywords")

	condition := &entity.SearchOrderCondition{
		StudentIDList:    parseInts(studentIds),
		ToOrgIDList:      parseInts(toOrgIds),
		StudentsKeywords: studentsKeywords,
		IntentSubjects:   intentSubjects,
		PublisherID:      parseInts(publisherID),
		OrderSourceList:  parseInts(orderSources),
		Keywords:         keywords,

		Status:  parseInts(status),
		OrderBy: orderBy,

		PageSize: parseInt(pageSize),
		Page:     parseInt(page),
	}

	createStartAtStr := c.Query("create_start_at")
	createEndAtStr := c.Query("create_end_at")
	if createStartAtStr != "" && createEndAtStr != "" {
		createStartAt := time.Unix(int64(parseInt(createStartAtStr)), 0)
		createEndAt := time.Unix(int64(parseInt(createEndAtStr)), 0)
		condition.CreateStartAt = &createStartAt
		condition.CreateEndAt = &createEndAt
	}

	return condition
}
