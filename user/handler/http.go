package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/core/debug"
	"github.com/stockfolioofficial/back-editfolio/domain"
	"github.com/stockfolioofficial/back-editfolio/util/echox"
	"net/http"
)

const (
	tag = "[USER] "
)

func NewUserController(useCase domain.UserUseCase) *UserController {
	return &UserController{useCase: useCase}
}

type UserController struct {
	useCase domain.UserUseCase
}

type CreatedUserResponse struct {
	Id uuid.UUID `json:"userId" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedUserResponse

func (c *UserController) createSuperAdmin(ctx echo.Context) error {
	var req CreateAdminRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "createSuperAdmin, request data bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	newId, err := c.useCase.CreateSuperAdminUser(ctx.Request().Context(), domain.CreateSuperAdminUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, CreatedUserResponse{Id: newId})
	case domain.ErrItemAlreadyExist:
		return ctx.JSON(http.StatusConflict, domain.ItemExist)
	default:
		log.WithError(err).Error(tag, "createSuperAdmin, unhandled error useCase.CreateSuperAdminUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

func (c *UserController) Bind(e *echo.Echo) {
	// get token
	e.POST("/sign-in", c.signInUser)

	// ===== INIT ====
	e.POST("/sa", c.createSuperAdmin)

	// ===== ADMIN =====
	// Fetch admin
	// v1, todo refactor
	e.GET("/admin", c.fetchAdmin,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	// v1, todo refactor
	e.GET("/admin/creator", c.fetchAdminCreator,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))

	// Self control
	// Get my info (admin)
	e.GET("/admin/me", echox.UserID(c.getAdminMyInfo), debug.JwtBypassOnDebug())
	// Update my info
	e.PUT("/admin/me", echox.UserID(c.updateAdminMyInfo), debug.JwtBypassOnDebug())
	// Update admin password
	e.PATCH("/admin/me/pw", echox.UserID(c.updateAdminMyPassword), debug.JwtBypassOnDebug())

	// ===== CUSTOMER =====
	// Customer control
	// Fetch customer
	// v1, todo refactor
	e.GET("/customer", c.fetchCustomer,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))

	// Create customer
	e.POST("/customer", c.createCustomer,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	// Get Customer
	e.GET("/customer/:userId", c.getCustomerDetailInfo,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))

	// Update customer
	e.PUT("/customer/:userId", c.updateCustomer,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))
	// Delete customer
	e.DELETE("/customer/:userId", c.deleteCustomerUser,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole, domain.AdminUserRole))

	e.GET("/customer/me", echox.UserID(c.getMyCustomerInfo),
		debug.JwtBypassOnDebugWithRole(domain.CustomerUserRole))

	// ===== SUPER_ADMIN =====
	// Create admin
	e.POST("/admin", c.createAdmin,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
	// Update admin info
	e.PUT("/admin/:userId", c.updateAdminBySuperAdmin,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
	// Update admin info
	e.PATCH("/admin/:userId/pw", c.updateAdminPasswordBySuperAdmin,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
	// Delete admin
	e.DELETE("/admin/:userId", c.deleteAdminBySuperAdmin,
		debug.JwtBypassOnDebugWithRole(domain.SuperAdminUserRole))
}
