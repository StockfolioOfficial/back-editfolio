package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/echox"
)

const (
	tag = "[ORDER] "
)

func NewOrderController(useCase domain.OrderUseCase) *OrderController {
	return &OrderController{useCase: useCase}
}

type OrderController struct {
	useCase domain.OrderUseCase
}

func (c *OrderController) Bind(e *echo.Echo) {

	//CUSTOMER
	// 진행중인 주문 가져오기
	e.GET("/order/recent-processing", echox.UserID(c.getRecentProcessingOrder), debug.JwtBypassOnDebugWithRole(domain.CustomerUserRole))
	// 진행중인 주문 완료
	e.POST("/order/recent-processing/done", echox.UserID(c.myOrderDone), debug.JwtBypassOnDebugWithRole(domain.CustomerUserRole))
	// 수정 접수
	e.POST("/order/recent-processing/edit", echox.UserID(c.myOrderEdit), debug.JwtBypassOnDebugWithRole(domain.CustomerUserRole))
	// 주문 접수
	e.POST("/order", echox.UserID(c.createOrder), debug.JwtBypassOnDebugWithRole(domain.CustomerUserRole))

	//ADMIN
	e.GET("/order/:orderId", c.getOrderDetailInfo,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	e.POST("/order/:orderId/assign-self", echox.UserID(c.orderAssignSelf),
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	e.PUT("/order/:orderId", c.updateOrderInfo,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	e.POST("/order/:orderId/edit-done", nil,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole)) // 대기

	// v1 - fetch, todo refactor
	e.GET("/order/ready", c.fetchOrderToReady,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	e.GET("/order/processing", echox.UserID(c.fetchOrderToProcessing),
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	e.GET("/order/done", c.fetchOrderToDone,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
}
