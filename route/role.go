package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary createRole
// @Description create a new role
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.CreateRoleRequest true "create role request"
// @Tags role
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/role [post]
func (s *Server) createRole(c *gin.Context) {
	req := new(entity.CreateRoleRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	id, err := service.GetRoleService().CreateRole(c.Request.Context(), req.Name, req.AuthIds)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, IdResponse{
		ID: id,
		ErrMsg:   "success",
	})
}

// @Summary listRoles
// @Description list all roles
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags role
// @Success 200 {array} entity.Role
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/roles [get]
func (s *Server) listRoles(c *gin.Context) {
	roles, err := service.GetRoleService().ListRole(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, RolesResponse{
		Roles: roles,
		ErrMsg:   "success",
	})
}