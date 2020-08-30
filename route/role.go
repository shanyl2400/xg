package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

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
	s.responseSuccessWithData(c, "id", id)
}

func (s *Server) listRoles(c *gin.Context) {
	roles, err := service.GetRoleService().ListRole(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "roles", roles)
}