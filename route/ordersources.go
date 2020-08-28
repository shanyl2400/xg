package route

import (
	"net/http"
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
