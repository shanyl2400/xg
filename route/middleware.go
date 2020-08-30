package route

import (
	"net/http"
	"strings"
	"xg/crypto"
	"xg/entity"

	"github.com/gin-gonic/gin"
)

func (s *Server) mustLogin(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	user, err := crypto.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_msg": "please login",
		})
		c.Abort()
		return
	}
	c.Set("user", user)
}

func (s *Server) mustOutOrg(c *gin.Context) {
	rawUser, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err_msg": "please login",
		})
		c.Abort()
	}
	user := rawUser.(*entity.JWTUser)
	if user.OrgId == entity.RootOrgId {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"err_msg": "only org can access this api",
		})
		c.Abort()
	}
}

func (s *Server) getJWTUser(c *gin.Context) *entity.JWTUser {
	user, ok := c.Get("user")
	if !ok {
		return nil
	}

	return user.(*entity.JWTUser)
}

func (s *Server) corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
