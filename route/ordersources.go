package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary listOrderSources
// @Description list all order sources
// @Accept json
// @Produce json
// @Tags orderSource
// @Success 200 {array} entity.OrderSource
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/order_sources [get]
func (s *Server) listOrderSources(c *gin.Context) {
	sources, err := service.GetOrderSourceService().ListOrderService(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "sources", sources)
}

// @Summary createOrderSource
// @Description create an order source
// @Accept json
// @Produce json
// @Tags orderSource
// @Param request body entity.CreateOrderSourceRequest true "create order source"
// @Success 200 {object} Response
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
	s.responseSuccessWithData(c, "id", id)
}
