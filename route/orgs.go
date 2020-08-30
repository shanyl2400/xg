package route

import (
	"net/http"
	"strconv"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) listOrgs(c *gin.Context) {
	orgs, err := service.GetOrgService().ListOrgs(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "orgs", orgs)
}

func (s *Server) listPendingOrgs(c *gin.Context) {
	orgs, err := service.GetOrgService().ListOrgsByStatus(c.Request.Context(), []int{entity.OrgStatusCreated})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "orgs", orgs)
}

func (s *Server) getOrgById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	org, err := service.GetOrgService().GetOrgById(c.Request.Context(), id)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "org", org)
}

func (s *Server) createOrg(c *gin.Context) {
	req := new(entity.CreateOrgRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	id, err := service.GetOrgService().CreateOrg(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}

func (s *Server) ApproveOrg(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrgService().CheckOrgById(c.Request.Context(), id, entity.OrgStatusCertified, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) RejectOrg(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrgService().CheckOrgById(c.Request.Context(), id, entity.OrgStatusRevoked, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}
