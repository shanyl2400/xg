package route

import (
	"errors"
	"net/http"
	"xg/da"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoAuth = errors.New("no authorization")
)

// @Summary login
// @Description login system
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.UserLoginRequest true "user login request"
// @Tags user
// @Success 200 {object} entity.UserLoginResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Failure 406 {object} Response
// @Router /api/user/login [post]
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
	c.JSON(http.StatusOK, UserLoginResponse{
		Data: data,
		ErrMsg:  "success",
	})
}

// @Summary updatePassword
// @Description update user password
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.UserUpdatePasswordRequest true "password to update"
// @Tags user
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/user/password [put]
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

// @Summary listUserAuthority
// @Description list user all authority
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags user
// @Success 200 {array} entity.Auth
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/user/authority [get]
func (s *Server) listUserAuthority(c *gin.Context) {
	user := s.getJWTUser(c)
	auth, err := service.GetUserService().ListUserAuthority(c.Request.Context(), user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, AuthorizationListResponse{
		AuthList: auth,
		ErrMsg:  "success",
	})
}

// @Summary listUsers
// @Description list all users
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags user
// @Success 200 {array} entity.UserInfo
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/users [get]
func (s *Server) listUsers(c *gin.Context) {
	condition := buildUsersSearchCondition(c)

	total, users, err := service.GetUserService().ListUsers(c.Request.Context(), condition)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Total: total,
		Users: users,
		ErrMsg:  "success",
	})
}

// @Summary resetPassword
// @Description reset user password
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "user id"
// @Tags user
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/user/reset/{id} [put]
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

// @Summary createUser
// @Description create a new user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.CreateUserRequest true "create user request"
// @Tags user
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/user [post]
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
	c.JSON(http.StatusOK, IdResponse{
		ID: id,
		ErrMsg:  "success",
	})
}


func buildUsersSearchCondition(c *gin.Context) da.SearchUserCondition {
	orderBy := c.Query("order_by")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	return da.SearchUserCondition{
		OrderBy: orderBy,

		PageSize: parseInt(pageSize),
		Page:     parseInt(page),
	}
}
