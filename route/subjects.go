package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary listRoles
// @Description list all roles
// @Accept mpfd
// @Produce json
// @Param parent_id path string true "subject parent id"
// @Tags subject
// @Success 200 {array} entity.Subject
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subjects [get]
func (s *Server) listSubjects(c *gin.Context) {
	parentID, ok := s.getParamInt(c, "parent_id")
	if !ok {
		return
	}

	subjects, err := service.GetSubjectService().ListSubjects(c.Request.Context(), parentID)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "subjects", subjects)
}
// @Summary createSubject
// @Description create a new subject
// @Accept json
// @Produce json
// @Param request body entity.CreateSubjectRequest true "create subject request"
// @Tags subject
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subject [post]
func (s *Server) createSubject(c *gin.Context) {
	req := new(entity.CreateSubjectRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	id, err := service.GetSubjectService().CreateSubject(c.Request.Context(), *req)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "id", id)
}
