package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) listOrderSources(c *gin.Context) {
	sources, err := service.GetOrderSourceService().ListOrderService(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "sources", sources)
}

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
