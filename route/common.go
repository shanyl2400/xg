package route

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"xg/log"

	"github.com/gin-gonic/gin"
)

type Response struct {
	ErrMsg string `json:"err_msg"`
}

type IdResponse struct {
	ID string `json:"id"`
	ErrMsg string `json:"err_msg"`
}

type SubjectsResponse struct {
	Subjects []string `json:"subjects"`
	ErrMsg string `json:"err_msg"`
}

func (s *Server) responseErr(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"err_msg": err.Error(),
	})
	c.Abort()
}

func (s *Server) responseSuccessWithData(c *gin.Context, key string, value interface{}) {
	c.JSON(http.StatusOK, gin.H{
		key:       value,
		"err_msg": "success",
	})
	c.Abort()
}

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
