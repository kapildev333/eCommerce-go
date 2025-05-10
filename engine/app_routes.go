package engine

import (
	"eCommerce-go/controllers"
	middleware "eCommerce-go/middleware"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/v1")
	controllers.ConfigAuthController(public)

	// Protected routes
	protected := router.Group("/v1")
	protected.Use(middleware.AuthMiddleware())
	controllers.ConfigAddressController(protected)
	controllers.ConfigPaymentController(protected)
	controllers.ConfigCategoryController(protected)
	controllers.ConfigProductController(protected)
}
