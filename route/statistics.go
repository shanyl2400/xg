package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xg/entity"
	"xg/service"
)

// @Summary summary
// @Description get system data summary
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags statistics
// @Success 200 {object} entity.SummaryInfo
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/summary [get]
func (s *Server) summary(c *gin.Context){
	summary, err := service.GetStatisticsService().Summary(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, SummaryResponse{
		Summary: summary,
		ErrMsg:   "success",
	})
}

// @Summary graph
// @Description get system data graph
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags statistics
// @Success 200 {object} entity.StatisticGraph
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/graph [get]
func (s *Server) graph(c *gin.Context) {
	studentsRecords, err := service.GetStatisticsService().SearchYearRecords(c.Request.Context(), entity.StudentStatisticsKey)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	performanceRecords, err := service.GetStatisticsService().SearchYearRecords(c.Request.Context(), entity.PerformanceStatisticsKey)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, GraphResponse{
		Graph: &entity.StatisticGraph{
			StudentsGraph:     studentsRecords,
			PerformancesGraph: performanceRecords,
		},
		ErrMsg:   "success",
	})
}