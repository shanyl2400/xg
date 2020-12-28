package route

import (
	"net/http"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary listSubjects
// @Description list all subjects
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param parent_id path string true "subject parent id"
// @Tags subject
// @Success 200 {object} SubjectsObjResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subjects/details/{parent_id} [get]
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
	c.JSON(http.StatusOK, SubjectsObjResponse{
		Subjects: subjects,
		ErrMsg:   "success",
	})
}

// @Summary listSubjectsTree
// @Description list all subjects tree
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags subject
// @Success 200 {array} SubjectsTreeObjResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subjects/tree [get]
func (s *Server) listSubjectsTree(c *gin.Context) {
	subjects, err := service.GetSubjectService().ListSubjectsTree(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, SubjectsTreeObjResponse{
		Subjects: subjects,
		ErrMsg:   "success",
	})
}

// @Summary listSubjectsAll
// @Description list all subjects
// @Accept mpfd
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags subject
// @Success 200 {array} SubjectsTreeObjResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subjects/all [get]
func (s *Server) listSubjectsAll(c *gin.Context) {
	subjects, err := service.GetSubjectService().ListSubjectsAll(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, SubjectsTreeObjResponse{
		Subjects: subjects,
		ErrMsg:   "success",
	})
}

// @Summary createSubject
// @Description create a new subject
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
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
	c.JSON(http.StatusOK, IdResponse{
		ID:     id,
		ErrMsg: "success",
	})
}

// @Summary batchCreateSubject
// @Description batch create new subjects
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.BatchCreateSubjectRequest true "create subject request"
// @Tags subject
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/subjects [post]
func (s *Server) batchCreateSubject(c *gin.Context) {
	req := new(entity.BatchCreateSubjectRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}

	err = service.GetSubjectService().BatchCreateSubject(c.Request.Context(), req.Data)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		ErrMsg: "success",
	})
}
