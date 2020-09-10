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

type OrgListInfo struct {
	Orgs  []*entity.Org `json:"orgs"`
	Total int           `json:"total"`
}
// @Summary listOrgs
// @Description list all organizations
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags organization
// @Success 200 {object} OrgListInfo
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orgs [get]
func (s *Server) listOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().ListOrgs(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrgsListResponse{
		Sources: &OrgListInfo{
			Orgs:  orgs,
			Total: count,
		},
		ErrMsg:  "success",
	})
}

// @Summary listPendingOrgs
// @Description list pending organizations to check
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags organization
// @Success 200 {object} OrgListInfo
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orgs/pending [get]
func (s *Server) listPendingOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().ListOrgsByStatus(c.Request.Context(), []int{entity.OrgStatusCreated})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrgsListResponse{
		Sources: &OrgListInfo{
			Orgs:  orgs,
			Total: count,
		},
		ErrMsg:  "success",
	})
}

// @Summary searchSubOrgs
// @Description search sub organizations with condition
// @Accept json
// @Produce json
// @Tags organization
// @Param Authorization header string true "With the bearer"
// @Param subjects query string false "search organizations with subjects"
// @Param address query string false "search organizations with address"
// @Success 200 {object} OrgListInfo
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/orgs/campus [get]
func (s *Server) searchSubOrgs(c *gin.Context) {
	count, orgs, err := service.GetOrgService().SearchSubOrgs(c.Request.Context(), buildOrgsSearchCondition(c))
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrgsListResponse{
		Sources: &OrgListInfo{
			Orgs:  orgs,
			Total: count,
		},
		ErrMsg:  "success",
	})
}

// @Summary getOrgById
// @Description get org by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "org id"
// @Tags organization
// @Success 200 {object} entity.Org
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org/{id} [get]
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
	c.JSON(http.StatusOK, OrgInfoResponse{
		Org:    org,
		ErrMsg: "success",
	})
}


// @Summary getOrgSubjectsById
// @Description get org subjects by id
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "org id"
// @Tags organization
// @Success 200 {object} SubjectsResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org/{id}/subjects [get]
func (s *Server) getOrgSubjectsById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	subjects, err := service.GetOrgService().GetOrgSubjectsById(c.Request.Context(), id)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, OrgSubjectsResponse{
		Subjects: subjects,
		ErrMsg:   "success",
	})
}

// @Summary createOrg
// @Description create org
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.CreateOrgWithSubOrgsRequest true "create org request"
// @Tags organization
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org [post]
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
	c.JSON(http.StatusOK, IdResponse{
		ID: id,
		ErrMsg:   "success",
	})
}

// @Summary approveOrg
// @Description approve org
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "org id"
// @Tags organization
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org/{id}/review/approve [put]
func (s *Server) approveOrg(c *gin.Context) {
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

// @Summary rejectOrg
// @Description reject org
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "org id"
// @Tags organization
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org/{id}/review/reject [put]
func (s *Server) rejectOrg(c *gin.Context) {
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

// @Summary revokeOrg
// @Description revoke org
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "org id"
// @Tags organization
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/org/{id}/revoke [put]
func (s *Server) revokeOrg(c *gin.Context) {
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
