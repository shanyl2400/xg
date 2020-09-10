package route

import (
	"net/http"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary listAuth
// @Description list all authorization
// @Accept json
// @Produce json
// @Tags auth
// @Param Authorization header string true "With the bearer"
// @Success 200 {object} entity.Auth
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/auths [get]
func (s *Server) listAuth(c *gin.Context) {
	auths, err := service.GetAuthService().ListAuths(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &AuthsListResponse{
		ErrMsg:   "success",
		AuthList: auths,
	})
}
