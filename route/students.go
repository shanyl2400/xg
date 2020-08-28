package route

import (
	"net/http"
	"strconv"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) createStudent(c *gin.Context) {
	req := new(entity.CreateStudentRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	id, status, err := service.GetStudentService().CreateStudent(c.Request.Context(), req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "result", entity.CreateStudentResponse{
		ID:     id,
		Status: status,
	})
}

func (s *Server) getStudentById(c *gin.Context) {
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	student, err := service.GetStudentService().GetStudentById(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "student", student)
}

func buildCondition(c *gin.Context) *entity.SearchStudentRequest {
	name := c.Query("name")
	telephone := c.Query("telephone")
	address := c.Query("address")

	orderBy := c.Query("order_by")

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	authorIdStr := c.Query("author_id")

	pageSize := 0
	page := 0
	authorId := 0

	temp, err := strconv.Atoi(pageSizeStr)
	if err == nil {
		pageSize = temp
	}
	temp, err = strconv.Atoi(pageStr)
	if err == nil {
		page = temp
	}
	temp, err = strconv.Atoi(authorIdStr)
	if err == nil {
		authorId = temp
	}

	return &entity.SearchStudentRequest{
		Name:         name,
		Telephone:    telephone,
		Address:      address,
		AuthorIDList: []int{authorId},
		OrderBy:      orderBy,
		PageSize:     pageSize,
		Page:         page,
	}
}
func (s *Server) searchPrivateStudents(c *gin.Context) {
	condition := buildCondition(c)
	user, ok := s.getJWTUser(c)
	if !ok {
		return
	}
	total, students, err := service.GetStudentService().SearchPrivateStudents(c.Request.Context(), condition, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "result", entity.StudentInfoList{
		Total:    total,
		Students: students,
	})
}
