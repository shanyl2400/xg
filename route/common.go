package route

import (
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

type IdResponse struct {
	ID int `json:"id"`
	ErrMsg string `json:"err_msg"`
}
type FileNameResponse struct {
	Name string `json:"name"`
	ErrMsg string `json:"err_msg"`
}

type SubjectsResponse struct {
	Subjects []string `json:"subjects"`
	ErrMsg string `json:"err_msg"`
}
type SubjectsObjResponse struct {
	Subjects []*entity.Subject `json:"subjects"`
	ErrMsg string `json:"err_msg"`
}

type AuthsListResponse struct {
	ErrMsg string `json:"err_msg"`
	AuthList []*entity.Auth `json:"auths"`
}
type AuthorizationListResponse struct {
	ErrMsg string `json:"err_msg"`
	AuthList []*entity.Auth `json:"authority"`
}

type OrderInfoListResponse struct {
	Data *entity.OrderInfoList `json:"data"`
	ErrMsg string `json:"err_msg"`
}
type OrderRecordResponse struct {
	Data *entity.OrderInfoWithRecords `json:"data"`
	ErrMsg string `json:"err_msg"`
}

type OrderPaymentRecordListResponse struct {
	Data *entity.PayRecordInfoList `json:"data"`
	ErrMsg string `json:"err_msg"`
}

type OrderSourcesListResponse struct {
	Sources []*entity.OrderSource `json:"sources"`
	ErrMsg string `json:"err_msg"`
}

type OrgsListResponse struct {
	Sources *OrgListInfo `json:"data"`
	ErrMsg string        `json:"err_msg"`
}

type OrgInfoResponse struct {
	Org 	*entity.Org 	`json:"org"`
	ErrMsg 	string        `json:"err_msg"`
}
type RolesResponse struct {
	Roles []*entity.Role `json:"roles"`
	ErrMsg 	string        `json:"err_msg"`
}
type SummaryResponse struct {
	Summary *entity.SummaryInfo `json:"summary"`
	ErrMsg 	string        `json:"err_msg"`
}
type GraphResponse struct {
	 Graph *entity.StatisticGraph `json:"graph"`
	ErrMsg 	string        `json:"err_msg"`
}

type OrgSubjectsResponse struct {
	Subjects 	[]string	`json:"subjects"`
	ErrMsg 	string        `json:"err_msg"`
}

type IDStatusResponse struct {
	Result entity.CreateStudentResponse `json:"result"`
	ErrMsg string `json:"err_msg"`
}

type StudentWithDetailsListResponse struct {
	Student *entity.StudentInfosWithOrders `json:"student"`
	ErrMsg string `json:"err_msg"`
}
type StudentListResponse struct {
	Result *entity.StudentInfoList `json:"result"`
	ErrMsg string `json:"err_msg"`
}

type UserLoginResponse struct {
	Data *entity.UserLoginResponse `json:"data"`
	ErrMsg string `json:"err_msg"`
}

type UserListResponse struct {
	Users []*entity.UserInfo `json:"users"`
	ErrMsg string `json:"err_msg"`
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
	if err == nil {
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
