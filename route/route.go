package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	"xg/entity"
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
		user.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.createUser)
		user.PUT("/reset", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.resetPassword)
	}
	users := route.Group("/users")
	{
		users.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.listUsers)
	}

	role := route.Group("/role")
	{
		role.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageRole}), s.createRole)
	}
	roles := route.Group("/roles")
	{
		roles.GET("/", s.mustLogin, s.listRoles)
	}

	student := route.Group("/student")
	{
		student.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.createStudent)
		student.GET("/:id", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent, entity.AuthListAllOrder}), s.getStudentById)
	}
	students := route.Group("/students")
	{
		student.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.searchStudents)
		students.GET("/private", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.searchPrivateStudents)
	}

	subject := route.Group("/subject")
	{
		subject.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageSubject}), s.createSubject)
	}
	subjects := route.Group("/subjects")
	{
		subjects.GET("/:parent_id", s.mustLogin, s.listSubjects)
	}
	org := route.Group("/org")
	{
		org.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrg}), s.createOrg)
		org.GET("/:id", s.mustLogin, s.getOrgById)
		org.PUT("/:id/review/reject", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.RejectOrg)
		org.PUT("/:id/review/approve", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.ApproveOrg)
	}
	orgs := route.Group("/orgs")
	{
		orgs.GET("/", s.mustLogin, s.listOrgs)
		orgs.GET("/pending", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.listPendingOrgs)
	}

	order := route.Group("/order")
	{
		order.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder}), s.createOrder)
		order.GET("/:id", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder, entity.AuthListAllOrder}), s.getOrderByID)
		order.PUT("/:id/signup", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.signupOrder)
		order.PUT("/:id/revoke", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.revokeOrder)
	}
	orders := route.Group("/orders")
	{
		orders.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.searchOrder)
		orders.GET("/author", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchOrder, entity.AuthDispatchSelfOrder}), s.searchOrderWithAuthor)
		orders.GET("/org", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.searchOrderWithOrgID)
	}

	payment := route.Group("/payment")
	{
		payment.POST("/:id/pay", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.payOrder)
		payment.POST("/:id/payback", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.paybackOrder)
		payment.PUT("/:id/review/accept", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.acceptPayment)
		payment.PUT("/:id/review/reject", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.rejectPayment)
	}
	payments := route.Group("/payments")
	{
		payments.GET("/pending", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.searchPendingPayRecord)
	}
	orderSource := route.Group("/order_sources")
	{
		orderSource.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrderSource}), s.createOrderSource)
	}
	orderSources := route.Group("/order_sources")
	{
		orderSources.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrderSource}), s.listOrderSources)
	}

	return route
}
