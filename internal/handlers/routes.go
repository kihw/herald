package handlers

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes enregistre toutes les routes de l'application
func RegisterRoutes(router *gin.Engine, testHandler *TestHandler) {
	api := router.Group("/api")
	{
		// Routes de validation
		RegisterValidationRoutes(api)

		// Routes d'analytics
		RegisterAnalyticsRoutes(api)

		// Routes de dashboard
		RegisterDashboardRoutes(api)

		// Routes de test (en mode dev uniquement)
		RegisterTestRoutes(api, testHandler)
	}
}

// RegisterValidationRoutes enregistre les routes de validation
func RegisterValidationRoutes(api *gin.RouterGroup) {
	validation := api.Group("/validation")
	{
		// Ces routes seraient définies dans validation_handler.go
		_ = validation
	}
}

// RegisterAnalyticsRoutes enregistre les routes d'analytics
func RegisterAnalyticsRoutes(api *gin.RouterGroup) {
	analytics := api.Group("/analytics")
	{
		// Ces routes seraient définies dans analytics_handler.go
		_ = analytics
	}
}

// RegisterDashboardRoutes enregistre les routes de dashboard
func RegisterDashboardRoutes(api *gin.RouterGroup) {
	dashboard := api.Group("/dashboard")
	{
		// Ces routes seraient définies dans dashboard_handler.go
		_ = dashboard
	}
}

// Note: RegisterTestRoutes is implemented in test_handler.go

// Note: RegisterNotificationRoutes is implemented in notification_handler.go
