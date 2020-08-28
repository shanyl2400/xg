package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xg/service"
)

func (s *Server) listOrgs(c *gin.Context){
	subjects, err := service.GetOrgService().ListOrgs(c.Request.Context())
	if err != nil{
		s.responseErr(c, http.StatusInternalServerError, err)
		return
	}
	s.responseSuccessWithData(c, "orgs", subjects)
}