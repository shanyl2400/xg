package route

import (
	"net/http"
	"strconv"
	"xg/da"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary handleStudentConflict
// @Description handle student conflict
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.HandleStudentConflictRequest true "handle student conflict request"
// @Tags student
// @Success 200 {string} string "success"
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/students_conflicts [put]
func (s *Server) handleStudentConflict(c *gin.Context) {
	req := new(entity.HandleStudentConflictRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	err = service.GetStudentConflictService().HandleStudentConflict(c.Request.Context(), *req)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	s.responseSuccess(c)
}

// @Summary searchStudentConflicts
// @Description search student conflict records
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param name query string false "search student with name"
// @Param telephone query string false "search student with telephone"
// @Param intent_subjects query string false "search student with intent_subjects"
// @Param address query string false "search student with address"
// @Param order_by query string false "search student order by column name"
// @Param page_size query int true "student list page size"
// @Param page query int false "student list page index"
// @Tags student
// @Success 200 {object} StudentConflictListResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/student_conflicts [get]
func (s *Server) searchStudentConflicts(c *gin.Context) {
	condition := buildStudentConflictsCondition(c)
	total, records, err := service.GetStudentConflictService().SearchStudentConflicts(c.Request.Context(), condition)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, StudentConflictListResponse{
		Result: &entity.StudentConflictsInfoList{
			Total:   total,
			Records: records,
		},
		ErrMsg: "success",
	})
}

// @Summary updateConflictStudentStatus
// @Description update conflict student status
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.HandleUpdateConflictStudentStatusRequest true "handle student conflict request"
// @Tags student
// @Success 200 {string} string "success"
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/students_conflicts [put]
func (s *Server) updateConflictStudentStatus(c *gin.Context) {
	req := new(entity.HandleUpdateConflictStudentStatusRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	err = service.GetStudentConflictService().UpdateConflictStudentStatus(c.Request.Context(), *req)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	s.responseSuccess(c)
}

// @Summary handleStudentConflictRecordStatus
// @Description handle Student Conflict Record Status
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.HandleStudentConflictStatusRequest true "handle student conflict request"
// @Tags student
// @Success 200 {string} string "success"
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/students_conflicts [put]
func (s *Server) handleStudentConflictRecordStatus(c *gin.Context) {
	req := new(entity.HandleStudentConflictStatusRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	err = service.GetStudentConflictService().HandleStudentConflictRecordStatus(c.Request.Context(), *req)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	s.responseSuccess(c)
}

func buildStudentConflictsCondition(c *gin.Context) da.SearchStudentConflictCondition {
	telephone := c.Query("telephone")
	authorIdStr := c.Query("author_id")

	orderBy := c.Query("order_by")

	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	status := parseInts(c.Query("status"))

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

	return da.SearchStudentConflictCondition{
		Telephone:    telephone,
		AuthorIDList: authorIdList,
		Status:       status,
		OrderBy:      orderBy,
		PageSize:     pageSize,
		Page:         page,
	}
}
