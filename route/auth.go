package route

import (
	"net/http"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) listAuth(c *gin.Context) {
	auths, err := service.GetAuthService().ListAuths(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "auths", auths)
}
