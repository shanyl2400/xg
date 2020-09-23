package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xg/da"
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
	summary, err := service.GetOrderStatisticsService().Summary(c.Request.Context())
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
	studentsRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(), da.SearchOrderStatisticsRecordCondition{
		Key:         entity.OrderStatisticKeyStudent,
	})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	performanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(), da.SearchOrderStatisticsRecordCondition{
		Key: entity.OrderStatisticKeyOrder,
	})
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

// @Summary orgGraph
// @Description get org data graph
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags statistics
// @Success 200 {object} PerformanceGraphResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/graph/org [get]
func (s *Server) orgGraph(c *gin.Context) {
	user := s.getJWTUser(c)
	subOrgs, err := service.GetOrgService().GetSubOrgs(c.Request.Context(), user.OrgId)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	orgIds := []int{user.OrgId}
	for i := range subOrgs {
		orgIds = append(orgIds, subOrgs[i].ID)
	}

	performanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(),
		da.SearchOrderStatisticsRecordCondition{
			Key: entity.OrderStatisticKeyOrder,
			OrgId: orgIds,
		})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, PerformanceGraphResponse{
		Graph: &entity.PerformancesGraph{
			PerformancesGraph: performanceRecords,
		},
		ErrMsg:   "success",
	})
}

// @Summary dispatchGraph
// @Description get dispatch data graph
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags statistics
// @Success 200 {object} PerformanceGraphResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/graph/dispatch [get]
func (s *Server) dispatchGraph(c *gin.Context) {
	user := s.getJWTUser(c)
	performanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(),
		da.SearchOrderStatisticsRecordCondition{
			Key: entity.OrderStatisticKeyOrder,
			PublisherId: []int{user.UserId},
		})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, PerformanceGraphResponse{
		Graph: &entity.PerformancesGraph{
			PerformancesGraph: performanceRecords,
		},
		ErrMsg:   "success",
	})
}


// @Summary orderSourceGraph
// @Description get order source data graph
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param id path string true "order source id"
// @Tags statistics
// @Success 200 {object} PerformanceGraphResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/graph/order_source/{id} [get]
func (s *Server) orderSourceGraph(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	performanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(),
		da.SearchOrderStatisticsRecordCondition{
			Key: entity.OrderStatisticKeyOrder,
			OrderSource: []int{id},
		})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, PerformanceGraphResponse{
		Graph: &entity.PerformancesGraph{
			PerformancesGraph: performanceRecords,
		},
		ErrMsg:   "success",
	})
}

// @Summary enterGraph
// @Description get enter data graph
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Tags statistics
// @Success 200 {object} AuthorPerformanceGraphResponse
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/statistics/graph/enter [get]
func (s *Server) enterGraph(c *gin.Context) {
	user := s.getJWTUser(c)
	publisherPerformanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(),
		da.SearchOrderStatisticsRecordCondition{
			Key: entity.OrderStatisticKeyOrder,
			PublisherId: []int{user.UserId},
		})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	authorPerformanceRecords, err := service.GetOrderStatisticsService().SearchRecordsMonth(c.Request.Context(),
		da.SearchOrderStatisticsRecordCondition{
		Key: entity.OrderStatisticKeyOrder,
		Author: []int{user.UserId},
	})
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, AuthorPerformanceGraphResponse{
		Graph: &entity.AuthorPerformancesGraph{
			PublisherPerformancesGraph: publisherPerformanceRecords,
			AuthorPerformancesGraph: authorPerformanceRecords,
		},
		ErrMsg:   "success",
	})
}