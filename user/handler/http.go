package handler

import (
	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/util/echox"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

const (
	tag = "[USER] "
)

func NewUserHttpHandler(useCase domain.UserUseCase) *HttpHandler {
	return &HttpHandler{useCase: useCase}
}

type HttpHandler struct {
	useCase domain.UserUseCase
}

type CreatedUserResponse struct {
	Id uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedUserResponse

func (h *HttpHandler) Bind(e *echo.Echo) {
	// get token
	e.POST("/sign-in", h.signInUser)

	// ADMIN
	// Self control
	// Update my info
	e.PUT("/admin/me", echox.UserID(h.updateAdminMyInfo), debug.JwtBypassOnDebug())
	// Update admin password
	e.PATCH("/admin/me/pw", echox.UserID(h.updateAdminMyPassword), debug.JwtBypassOnDebug())

	// CUSTOMER
	// Customer control
	// Create customer
	e.POST("/customer", h.createCustomer, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	// Update customer
	e.PUT("/customer/:userId", h.updateCustomer, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))
	// Delete customer
	e.DELETE("/customer/:userId", h.deleteCustomerUser, debug.JwtBypassOnDebugWithRole(domain.AdminUserRole))

	// SUPER_ADMIN
	// Create admin
	e.POST("/admin", h.createAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
	// Update admin info
	e.PUT("/admin/:userId", h.updateAdminBySuperAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
	// Delete admin
	e.DELETE("/admin/:userId", h.deleteAdminBySuperAdmin, debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
}
