package route

import (
	"errors"
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) login(c *gin.Context) {
	userLoginReq := new(entity.UserLoginRequest)
	err := c.ShouldBind(userLoginReq)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	if userLoginReq.Name == "" || userLoginReq.Password == "" {
		s.responseErr(c, http.StatusBadRequest, errors.New("no name or password"))
		return
	}

	data, err := service.GetUserService().Login(c.Request.Context(), userLoginReq.Name, userLoginReq.Password)
	if err != nil {
		s.responseErr(c, http.StatusNotAcceptable, err)
		return
	}
	s.responseSuccessWithData(c, "data", data)
}

func (s *Server) updatePassword(c *gin.Context) {
	newPasswordReq := new(entity.UserUpdatePasswordRequest)
	err := c.ShouldBind(newPasswordReq)

	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	err = service.GetUserService().UpdatePassword(c.Request.Context(), newPasswordReq.NewPassword, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) listUserAuthority(c *gin.Context) {
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	auth, err := service.GetUserService().ListUserAuthority(c.Request.Context(), user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "authority", auth)
}
