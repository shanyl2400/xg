package route

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
}

func Get() *gin.Engine {
	route := gin.Default()
	s := new(Server)

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	user := route.Group("/user")
	{
		user.GET("/authority", s.mustLogin, s.listUserAuthority)
		user.POST("/login", s.login)
		user.PUT("/password", s.mustLogin, s.updatePassword)
	}

	student := route.Group("/student")
	{
		student.POST("", s.mustLogin, s.createStudent)
		student.GET("/:id", s.mustLogin, s.getStudentById)
	}
	students := route.Group("/students")
	{
		students.GET("/private", s.mustLogin, s.searchPrivateStudents)
	}

	subject := route.Group("/subjects")
	{
		subject.GET("/:parent_id", s.mustLogin, s.listSubjects)
	}

	org := route.Group("/orgs")
	{
		org.GET("/", s.mustLogin, s.listOrgs)
	}

	order := route.Group("/order")
	{
		order.POST("/", s.mustLogin, s.createOrder)
	}

	orderSources := route.Group("/order_sources")
	{
		orderSources.GET("/", s.mustLogin, s.listOrderSources)
	}

	return route
}
