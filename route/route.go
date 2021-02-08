package route

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"xg/entity"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
}

func Get() *gin.Engine {
	route := gin.Default()
	s := new(Server)

	allowOrigin := os.Getenv("allow_origin")
	allowOriginParts := strings.Split(allowOrigin, ",")
	fmt.Println("Allow Origin:", allowOriginParts)
	if len(allowOriginParts) < 1 {
		panic("invalid allow origin")
	}
	route.Use(cors.New(cors.Config{
		AllowOrigins:     allowOriginParts,
		AllowMethods:     []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "Accept", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			for i := range allowOriginParts {
				if origin == allowOriginParts[i] {
					return true
				}
			}
			return false
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
		users.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageUser, entity.AuthListAllOrder}), s.listUsers)
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
		subjects.POST("/", s.mustLogin, s.hasPermission([]int{entity.AuthManageSubject}), s.batchCreateSubject)
		subjects.GET("/details/:parent_id", s.mustLogin, s.listSubjects)
		subjects.GET("/tree", s.mustLogin, s.listSubjectsTree)
		subjects.GET("/all", s.mustLogin, s.listSubjectsAll)
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
		orders.GET("/remarks", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder, entity.AuthListAllOrder, entity.AuthListOrgOrder}), s.searchOrderRemarks)
		orders.PUT("/marks", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchSelfOrder, entity.AuthDispatchOrder, entity.AuthListAllOrder, entity.AuthListOrgOrder}), s.markOrderRemarks)
		orders.GET("/", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.searchOrder)
		orders.GET("/export", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.exportOrder)
		orders.GET("/author", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchOrder, entity.AuthDispatchSelfOrder}), s.searchOrderWithAuthor)
		orders.GET("/org", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.searchOrderWithOrgID)
	}

	notify := api.Group("/notify")
	{
		notify.PUT("/orders/:id", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder, entity.AuthEnterStudent}), s.markOrderNotify)
	}

	notifies := api.Group("/notifies")
	{
		notifies.GET("/orders/author", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.searchAuthorNotifies)
		notifies.GET("/orders", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.searchNotifies)
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
		statistics.GET("/table", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.statisticsTable)
		statistics.GET("/group", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.statisticsGroup)
		statistics.GET("/graph", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.graph)
		statistics.GET("/graph/org", s.mustLogin, s.hasPermission([]int{entity.AuthListOrgOrder}), s.orgGraph)
		statistics.GET("/graph/dispatch", s.mustLogin, s.hasPermission([]int{entity.AuthDispatchOrder}), s.dispatchGraph)
		statistics.GET("/graph/enter", s.mustLogin, s.hasPermission([]int{entity.AuthEnterStudent}), s.enterGraph)
		statistics.GET("/graph/order_source/:id", s.mustLogin, s.hasPermission([]int{entity.AuthListAllOrder}), s.orderSourceGraph)
	}

	socks := api.Group("/socks")
	{
		socks.GET("/register", s.mustLogin, s.registerSocks)
	}

	uploader := api.Group("/upload")
	{
		uploader.POST("/:partition", s.mustLogin, s.uploadFile)
	}

	//访问上传文件
	route.StaticFS("/data", http.Dir(os.Getenv("xg_upload_path")))

	return route
}
