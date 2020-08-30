package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

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
