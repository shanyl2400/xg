package route

import (
	"errors"
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

var(
	ErrNoAuth = errors.New("no authorization")
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
	user := s.getJWTUser(c)
	err = service.GetUserService().UpdatePassword(c.Request.Context(), newPasswordReq.NewPassword, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) listUserAuthority(c *gin.Context) {
	user := s.getJWTUser(c)
	auth, err := service.GetUserService().ListUserAuthority(c.Request.Context(), user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "authority", auth)
}

func (s *Server) listUsers(c *gin.Context) {
	users, err := service.GetUserService().ListUsers(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "users", users)
}

func (s *Server) resetPassword(c *gin.Context) {
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	user := s.getJWTUser(c)
	err := service.GetUserService().ResetPassword(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) createUser(c *gin.Context) {
	req := new(entity.CreateUserRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, err := service.GetUserService().CreateUser(c.Request.Context(), req)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}
