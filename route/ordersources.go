package route

import (
	"net/http"
	"strconv"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary listOrderSources
// @Description list all order sources
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags orderSource
// @Success 200 {object} OrderSourcesListResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order_sources [get]
func (s *Server) listOrderSources(c *gin.Context) {
	sources, err := service.GetOrderSourceService().ListOrderService(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, OrderSourcesListResponse{
		Sources:     sources,
		ErrMsg: "success",
	})
}

// @Summary createOrderSource
// @Description create an order source
// @Accept json
// @Produce json
// @Tags orderSource
// @Param Authorization header string true "With the bearer"
// @Param request body entity.CreateOrderSourceRequest true "create order source"
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order_source [post]
func (s *Server) createOrderSource(c *gin.Context) {
	req := new(entity.CreateOrderSourceRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	id, err := service.GetOrderSourceService().CreateOrderService(c.Request.Context(), req.Name)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, IdResponse{
		ID:     id,
		ErrMsg: "success",
	})
}


// @Summary deleteOrderSource
// @Description delete an order source
// @Accept json
// @Produce json
// @Tags orderSource
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order source id"
// @Success 200 {string} string "success"
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order_source [delete]
func (s *Server) deleteOrderSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	err = service.GetOrderSourceService().DeleteOrderService(c.Request.Context(),id)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	s.responseSuccess(c)
}
