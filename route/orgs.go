package route

import (
	"net/http"
	"strconv"
	"strings"
	"xg/da"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

type OrgListResponse struct {
	Orgs  []*entity.Org `json:"orgs"`
	Total int           `json:"total"`
}

func (s *Server) listOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().ListOrgs(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", OrgListResponse{
		Total: count,
		Orgs:  orgs,
	})
}

func (s *Server) listPendingOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().ListOrgsByStatus(c.Request.Context(), []int{entity.OrgStatusCreated})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", OrgListResponse{
		Total: count,
		Orgs:  orgs,
	})
}

func (s *Server) searchSubOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().SearchSubOrgs(c.Request.Context(), buildOrgsSearchCondition(c))
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "data", OrgListResponse{
		Total: count,
		Orgs:  orgs,
	})
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

func (s *Server) getOrgSubjectsById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	org, err := service.GetOrgService().GetOrgSubjectsById(c.Request.Context(), id)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "subjects", org)
}

func (s *Server) createOrg(c *gin.Context) {
	req := new(entity.CreateOrgWithSubOrgsRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	id, err := service.GetOrgService().CreateOrgWithSubOrgs(c.Request.Context(), req, user)
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
	err = service.GetOrgService().CheckOrgById(c.Request.Context(), id, entity.OrgStatusRejected, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func (s *Server) RevokeOrg(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetOrgService().RevokeOrgById(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

func buildOrgsSearchCondition(c *gin.Context) da.SearchOrgsCondition {
	subjects := make([]string, 0)
	subjectsParam := c.Query("subjects")
	if subjectsParam != "" {
		subjects = strings.Split(subjectsParam, ",")
	}
	return da.SearchOrgsCondition{
		Subjects: subjects,
		Address:  c.Query("address"),
	}
}
