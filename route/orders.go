package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xg/entity"
	"xg/service"
)

func (s *Server) createOrder(c *gin.Context) {
	req := new(entity.CreateOrderRequest)
	err := c.ShouldBind(req)
	if err != nil{
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok{
		return
	}
	id, err := service.GetOrderService().CreateOrder(c.Request.Context(), req, user)
	if err != nil{
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}