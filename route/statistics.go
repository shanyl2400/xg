package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xg/entity"
	"xg/service"
)

func (s *Server) summary(c *gin.Context){
	summary, err := service.GetStatisticsService().Summary(c.Request.Context())
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "summary", summary)
}

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
	s.responseSuccessWithData(c, "graph", entity.StatisticGraph{
		StudentsGraph:     studentsRecords,
		PerformancesGraph: performanceRecords,
	})
}