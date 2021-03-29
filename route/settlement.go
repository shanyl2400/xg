package route

import (
	"net/http"
	"strconv"
	"time"
	"xg/da"
	"xg/entity"
	"xg/service"

	"github.com/gin-gonic/gin"
)

// @Summary createSettlement
// @Description create a new settlement
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param request body entity.CreateSettlementRequest true "create settlement request"
// @Tags role
// @Success 200 {string} "success"
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/settlement [post]
func (s *Server) createSettlement(c *gin.Context) {
	req := new(entity.CreateSettlementRequest)
	err := c.ShouldBind(req)
	if err != nil {
		s.responseErr(c, http.StatusBadRequest, err)
		return
	}
	user := s.getJWTUser(c)
	err = service.GetSettlementService().CreateSettlement(c.Request.Context(), *req, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccess(c)
}

// @Summary searchSettlements
// @Description search Settlements
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer"
// @Param start_at query string false "search settlements with start_at"
// @Param end_at query string false "search settlements with end_at"
// @Param order_by query string false "search settlements order by column name"
// @Param page_size query int true "settlements list page size"
// @Param page query int false "settlements list page index"
// @Tags role
// @Success 200 {array} entity.SettlementData
// @Failure 500 {object} Response
// @Failure 400 {object} Response
// @Router /api/settlements [get]
func (s *Server) searchSettlements(c *gin.Context) {
	condition := buildSettlementCondition(c)
	user := s.getJWTUser(c)
	total, settlements, err := service.GetSettlementService().SearchSettlements(c.Request.Context(), *condition, user)
	if err != nil {
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, SettlementsResponse{
		Total:       total,
		Settlements: settlements,
		ErrMsg:      "success",
	})
}

func buildSettlementCondition(c *gin.Context) *da.SearchSettlementsCondition {
	orderBy := c.Query("order_by")
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize := 0
	page := 0

	temp, err := strconv.Atoi(pageSizeStr)
	if err == nil {
		pageSize = temp
	}
	temp, err = strconv.Atoi(pageStr)
	if err == nil {
		page = temp
	}

	//开始截至时间筛选
	startAt := parseInt(c.Query("start_at"))
	endAt := parseInt(c.Query("end_at"))
	var startAtObj *time.Time
	var endAtObj *time.Time
	if startAt > 0 {
		s := time.Unix(int64(startAt), 0)
		startAtObj = &s
	}
	if endAt > 0 {
		e := time.Unix(int64(endAt), 0)
		endAtObj = &e
	}

	return &da.SearchSettlementsCondition{
		StartAt: startAtObj,
		EndAt:   endAtObj,

		OrderBy:  orderBy,
		PageSize: pageSize,
		Page:     page,
	}
}
