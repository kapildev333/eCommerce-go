package engine

import (
	"eCommerce-go/controllers"
	"eCommerce-go/utils"
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/v1")
	controllers.ConfigAuthController(public)

	// Protected routes
	protected := router.Group("/v1")
	protected.Use(utils.AuthMiddleware())
	controllers.ConfigAddressController(protected)
	// Add other protected routes here
}
