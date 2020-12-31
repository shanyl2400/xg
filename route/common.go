package route

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"xg/entity"
	"xg/log"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ErrMsg string `json:"err_msg"`
}

func (r Response) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type IdResponse struct {
	ID     int    `json:"id"`
	ErrMsg string `json:"err_msg"`
}

func (r IdResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type FileNameResponse struct {
	Name   string `json:"name"`
	ErrMsg string `json:"err_msg"`
}

func (r FileNameResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type SubjectsResponse struct {
	Subjects []string `json:"subjects"`
	ErrMsg   string   `json:"err_msg"`
}

func (r SubjectsResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type SubjectsObjResponse struct {
	Subjects []*entity.Subject `json:"subjects"`
	ErrMsg   string            `json:"err_msg"`
}

type SubjectsTreeObjResponse struct {
	Subjects []*entity.SubjectTreeNode `json:"subjects"`
	ErrMsg   string                    `json:"err_msg"`
}

func (r SubjectsObjResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type AuthsListResponse struct {
	ErrMsg   string         `json:"err_msg"`
	AuthList []*entity.Auth `json:"auths"`
}

func (r AuthsListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type AuthorizationListResponse struct {
	ErrMsg   string         `json:"err_msg"`
	AuthList []*entity.Auth `json:"authority"`
}

func (r AuthorizationListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrderInfoListResponse struct {
	Data   *entity.OrderInfoList `json:"data"`
	ErrMsg string                `json:"err_msg"`
}

func (r OrderInfoListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrderRecordResponse struct {
	Data   *entity.OrderInfoWithRecords `json:"data"`
	ErrMsg string                       `json:"err_msg"`
}

func (r OrderRecordResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrderPaymentRecordListResponse struct {
	Data   *entity.PayRecordInfoList `json:"data"`
	ErrMsg string                    `json:"err_msg"`
}

type OrderNotifyResponse struct {
	Data   []*entity.OrderNotify `json:"data"`
	Total  int                   `json:"total"`
	ErrMsg string                `json:"err_msg"`
}

func (r OrderPaymentRecordListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrderSourcesListResponse struct {
	Sources []*entity.OrderSource `json:"sources"`
	ErrMsg  string                `json:"err_msg"`
}

func (r OrderSourcesListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrgsListResponse struct {
	Data   *OrgListInfo `json:"data"`
	ErrMsg string       `json:"err_msg"`
}

type SubOrgsListResponse struct {
	Data   *SubOrgListInfo `json:"data"`
	ErrMsg string          `json:"err_msg"`
}

func (r SubOrgsListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}
func (r OrgsListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrgInfoResponse struct {
	Org    *entity.Org `json:"org"`
	ErrMsg string      `json:"err_msg"`
}

func (r OrgInfoResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type RolesResponse struct {
	Roles  []*entity.Role `json:"roles"`
	ErrMsg string         `json:"err_msg"`
}

func (r RolesResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type SummaryResponse struct {
	Summary *entity.SummaryInfo `json:"summary"`
	ErrMsg  string              `json:"err_msg"`
}

func (r SummaryResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type StatisticTableResponse struct {
	Data   *entity.OrderStatisticTable `json:"data"`
	ErrMsg string                      `json:"err_msg"`
}

type StatisticTimeTableResponse struct {
	Data   []*entity.OrderStatisticGroupTableItem `json:"data"`
	ErrMsg string                                 `json:"err_msg"`
}

func (r StatisticTableResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type GraphResponse struct {
	Graph  *entity.StatisticGraph `json:"graph"`
	ErrMsg string                 `json:"err_msg"`
}

func (r GraphResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type PerformanceGraphResponse struct {
	Graph  *entity.PerformancesGraph `json:"graph"`
	ErrMsg string                    `json:"err_msg"`
}

func (r PerformanceGraphResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type AuthorPerformanceGraphResponse struct {
	Graph  *entity.AuthorPerformancesGraph `json:"graph"`
	ErrMsg string                          `json:"err_msg"`
}

func (r AuthorPerformanceGraphResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type OrgSubjectsResponse struct {
	Subjects []string `json:"subjects"`
	ErrMsg   string   `json:"err_msg"`
}

func (r OrgSubjectsResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type IDStatusResponse struct {
	Result entity.CreateStudentResponse `json:"result"`
	ErrMsg string                       `json:"err_msg"`
}

func (r IDStatusResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type StudentWithDetailsListResponse struct {
	Student *entity.StudentInfosWithOrders `json:"student"`
	ErrMsg  string                         `json:"err_msg"`
}

func (r StudentWithDetailsListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type StudentListResponse struct {
	Result *entity.StudentInfoList `json:"result"`
	ErrMsg string                  `json:"err_msg"`
}

func (r StudentListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type UserLoginResponse struct {
	Data   *entity.UserLoginResponse `json:"data"`
	ErrMsg string                    `json:"err_msg"`
}

func (r UserLoginResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

type UserListResponse struct {
	Total  int                `json:"total"`
	Users  []*entity.UserInfo `json:"users"`
	ErrMsg string             `json:"err_msg"`
}

func (r UserListResponse) Error() error {
	if r.ErrMsg == "success" {
		return nil
	}
	return errors.New(r.ErrMsg)
}

func (s *Server) responseErr(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"err_msg": err.Error(),
	})
	c.Abort()
}

//func (s *Server) responseSuccessWithData(c *gin.Context, key string, value interface{}) {
//	c.JSON(http.StatusOK, gin.H{
//		key:       value,
//		"err_msg": "success",
//	})
//	c.Abort()
//}

func (s *Server) responseSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"err_msg": "success",
	})
	c.Abort()
}

func (s *Server) getParamInt(c *gin.Context, key string) (int, bool) {
	valueStr := c.Param(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Warning.Println("Parse ", key, " failed, err:", err)
		s.responseErr(c, http.StatusBadRequest, fmt.Errorf("parse %v failed", key))
		return -1, false
	}
	return value, true
}

func parseInt(str string) int {
	id, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return id
}
func parseInts(str string) []int {
	strList := strings.Split(str, ",")
	ret := make([]int, 0)
	for i := range strList {
		id, err := strconv.Atoi(strList[i])
		if err == nil {
			ret = append(ret, id)
		}
	}
	if len(ret) < 1 {
		return nil
	}
	return ret
}
