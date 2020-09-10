package route

import (
	"net/http"
	"strconv"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary createStudent
// @Description create a new student
// @Accept json
// @Produce json
// @Param request body entity.CreateStudentRequest true "create student request"
// @Tags student
// @Success 200 {object} IdResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/student [post]
func (s *Server) createStudent(c *gin.Context) {
	req := new(entity.CreateStudentRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
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

// @Summary createStudent
// @Description create a new student
// @Accept json
// @Produce json
// @Param id path string true "student id"
// @Tags student
// @Success 200 {object} entity.StudentInfosWithOrders
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/student/{id} [get]
func (s *Server) getStudentById(c *gin.Context) {
	id, ok := s.getParamInt(c, "id")
	if !ok {
		return
	}
	user := s.getJWTUser(c)
	student, err := service.GetStudentService().GetStudentById(c.Request.Context(), id, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "student", student)
}

// @Summary searchPrivateStudents
// @Description search private students of user with condition
// @Accept json
// @Produce json
// @Param name query string false "search student with name"
// @Param telephone query string false "search student with telephone"
// @Param address query string false "search student with address"
// @Param order_by query string false "search student order by column name"
// @Param page_size query int true "student list page size"
// @Param page query int false "student list page index"
// @Tags student
// @Success 200 {object} entity.StudentInfoList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/students/private [get]
func (s *Server) searchPrivateStudents(c *gin.Context) {
	condition := buildCondition(c)
	user := s.getJWTUser(c)
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


// @Summary searchStudents
// @Description search students with condition
// @Accept json
// @Produce json
// @Param name query string false "search student with name"
// @Param telephone query string false "search student with telephone"
// @Param address query string false "search student with address"
// @Param author_id query string false "search student with author_id"
// @Param order_by query string false "search student order by column name"
// @Param page_size query int true "student list page size"
// @Param page query int false "student list page index"
// @Tags student
// @Success 200 {object} entity.StudentInfoList
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/students/private [get]
func (s *Server) searchStudents(c *gin.Context) {
	condition := buildCondition(c)
	user := s.getJWTUser(c)
	total, students, err := service.GetStudentService().SearchStudents(c.Request.Context(), condition, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "result", entity.StudentInfoList{
		Total:    total,
		Students: students,
	})
}

func buildCondition(c *gin.Context) *entity.SearchStudentRequest {
	name := c.Query("name")
	telephone := c.Query("telephone")
	address := c.Query("address")
	authorIdStr := c.Query("author_id")

	orderBy := c.Query("order_by")

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

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

	authorIdList := make([]int, 0)
	if authorId > 0 {
		authorIdList = append(authorIdList, authorId)
	}

	return &entity.SearchStudentRequest{
		Name:         name,
		Telephone:    telephone,
		Address:      address,
		AuthorIDList: authorIdList,
		OrderBy:      orderBy,
		PageSize:     pageSize,
		Page:         page,
	}
}