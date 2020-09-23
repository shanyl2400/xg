package route

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
	"xg/entity"
)

type Server struct {
}

func Get() *gin.Engine {
	route := gin.Default()
	s := new(Server)

	allowOrigin := os.Getenv("allow_origin")
	fmt.Println("Allow Origin:", allowOrigin)
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowOrigin},
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == allowOrigin
		},
		MaxAge: 12 * time.Hour,
	}))
	api := route.Group("/api")
	user := api.Group("/user")
	{
		user.GET("/authority", s.mustLogin, s.listUserAuthority)
		user.POST("/login", s.login)
		user.PUT("/password", s.mustLogin, s.updatePassword)
		user.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.createUser)
		user.PUT("/reset/:id", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.resetPassword)
	}
	users := api.Group("/users")
	{
		users.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser}), s.listUsers)
	}

	role := api.Group("/role")
	{
		role.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageRole}), s.createRole)
	}
	roles := api.Group("/roles")
	{
		roles.GET("/", s.mustLogin, s.listRoles)
	}

	auth := api.Group("/auths")
	{
		auth.GET("/", s.mustLogin, s.listAuth)
	}

	student := api.Group("/student")
	{
		student.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.createStudent)
		student.GET("/:id", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent, entity.AuthListAllOrder, entity.AuthDispatchOrder}), s.getStudentById)
	}
	students := api.Group("/students")
	{
		students.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder, entity.AuthDispatchOrder}), s.searchStudents)
		students.GET("/private", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.searchPrivateStudents)
	}

	subject := api.Group("/subject")
	{
		subject.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageSubject}), s.createSubject)
	}
	subjects := api.Group("/subjects")
	{
		subjects.GET("/:parent_id", s.mustLogin, s.listSubjects)
	}
	org := api.Group("/org")
	{
		org.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrg}), s.createOrg)
		org.GET("/:id", s.mustLogin, s.getOrgById)
		org.GET("/:id/subjects", s.mustLogin, s.getOrgSubjectsById)
		org.PUT("/:id/revoke", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.revokeOrg)
		org.PUT("/:id/review/reject", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.rejectOrg)
		org.PUT("/:id/review/approve", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.approveOrg)

		org.PUT("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageSelfOrg}), s.updateSelfOrgById)
		org.PUT("/:id", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrg}), s.updateOrgById)
	}
	orgs := api.Group("/orgs")
	{
		orgs.GET("/", s.mustLogin, s.listOrgs)
		orgs.GET("/pending", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrg}), s.listPendingOrgs)
		orgs.GET("/campus", s.mustLogin, s.searchSubOrgs)
	}

	order := api.Group("/order")
	{
		order.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder}), s.createOrder)
		order.GET("/:id", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder, entity.AuthListAllOrder, entity.AuthListOrgOrder}), s.getOrderByID)
		order.PUT("/:id/signup", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.signupOrder)
		order.PUT("/:id/deposit", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.depositOrder)
		order.PUT("/:id/revoke", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.revokeOrder)
		order.PUT("/:id/invalid", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.invalidOrder)
		order.POST("/:id/mark", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder, entity.AuthListAllOrder, entity.AuthListOrgOrder}), s.addOrderMark)
	}
	orders := api.Group("/orders")
	{
		orders.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.searchOrder)
		orders.GET("/author", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchOrder, entity.AuthDispatchSelfOrder}), s.searchOrderWithAuthor)
		orders.GET("/org", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.searchOrderWithOrgID)
	}

	payment := api.Group("/payment")
	{
		payment.POST("/:id/pay", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.payOrder)
		payment.POST("/:id/payback", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.paybackOrder)
		payment.PUT("/:id/review/accept", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.acceptPayment)
		payment.PUT("/:id/review/reject", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.rejectPayment)
	}
	payments := api.Group("/payments")
	{
		payments.GET("/pending", s.mustLogin, s.hasPermission([]int{entity.AuthCheckOrder}), s.searchPendingPayRecord)
	}
	orderSource := api.Group("/order_source")
	{
		orderSource.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageOrderSource}), s.createOrderSource)
	}
	orderSources := api.Group("/order_sources")
	{
		orderSources.GET("/", s.mustLogin, s.listOrderSources)
	}

	statistics := api.Group("/statistics")
	{
		statistics.GET("/summary", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.summary)
		statistics.GET("/graph", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.graph)
		statistics.GET("/graph/org", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.orgGraph)
		statistics.GET("/graph/dispatch", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchOrder}), s.dispatchGraph)
		statistics.GET("/graph/enter", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.enterGraph)
		statistics.GET("/graph/order_source/:id", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.orderSourceGraph)
	}

	uploader := api.Group("/upload")
	{
		uploader.POST("/:partition", s.uploadFile)
	}

	//访问上传文件
	route.StaticFS("/data", http.Dir(os.Getenv("xg_upload_path")))

	return route
}
