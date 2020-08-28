package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xg/service"
)

func (s *Server) listSubjects(c *gin.Context){
	parentId, ok := s.getParamInt(c, "parent_id")
	if !ok {
		return
	}

	subjects, err := service.GetSubjectService().ListSubjects(c.Request.Context(), parentId)
	if err != nil{
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "subjects", subjects)
}