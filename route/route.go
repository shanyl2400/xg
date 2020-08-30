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
		user.POST("/", s.mustLogin, s.createUser)
		user.PUT("/reset", s.mustLogin, s.resetPassword)
	}
	users := route.Group("/users")
	{
		users.GET("/", s.mustLogin, s.listUsers)
	}

	role := route.Group("/role")
	{
		role.POST("/", s.mustLogin, s.createRole)
	}
	roles := route.Group("/roles")
	{
		roles.GET("/", s.mustLogin, s.listRoles)
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
		subject.POST("/", s.mustLogin, s.createSubject)
	}
	org := route.Group("/org")
	{
		org.POST("/", s.mustLogin, s.createOrg)
		org.GET("/:id", s.mustLogin, s.getOrgById)
		org.PUT("/:id/reject", s.mustLogin, s.RejectOrg)
		org.PUT("/:id/approve", s.mustLogin, s.ApproveOrg)
	}
	orgs := route.Group("/orgs")
	{
		orgs.GET("/", s.mustLogin, s.listOrgs)
		orgs.GET("/pending", s.mustLogin, s.listPendingOrgs)
	}

	order := route.Group("/order")
	{
		order.POST("/", s.mustLogin, s.createOrder)
		order.GET("/:id", s.mustLogin, s.getOrderByID)
	}
	orders := route.Group("/orders")
	{
		orders.GET("/", s.mustLogin, s.searchOrder)
		orders.GET("/author", s.mustLogin, s.searchOrderWithAuthor)
		orders.GET("/org", s.mustLogin, s.searchOrderWithOrgID)
		orders.GET("/payment", s.mustLogin, s.searchPendingPayRecord)
	}

	orderSources := route.Group("/order_sources")
	{
		orderSources.GET("/", s.mustLogin, s.listOrderSources)
		orderSources.POST("/", s.mustLogin, s.createOrderSource)
	}

	return route
}
