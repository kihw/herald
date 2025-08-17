package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"lol-match-exporter/internal/services"
)

type OptimizedAnalyticsHandler struct {
	OptimizedService *services.OptimizedAnalyticsService
}

func NewOptimizedAnalyticsHandler(optimizedService *services.OptimizedAnalyticsService) *OptimizedAnalyticsHandler {
	return &OptimizedAnalyticsHandler{
		OptimizedService: optimizedService,
	}
}

// GetPeriodStatsOptimized returns analytics for a specific time period using optimized service
// GET /api/analytics/v2/period/:period
func (h *OptimizedAnalyticsHandler) GetPeriodStatsOptimized(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	period := c.Param("period")
	if period == "" {
		period = "week"
	}

	// Validate period
	validPeriods := map[string]bool{
		"today":  true,
		"week":   true,
		"month":  true,
		"season": true,
		"all":    true,
	}

	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_period",
			"message": "Valid periods: today, week, month, season, all",
		})
		return
	}

	// Use optimized async method
	stats, err := h.OptimizedService.GetPeriodStatsAsync(uid, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to retrieve period statistics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"meta": gin.H{
			"period":      period,
			"user_id":     uid,
			"cached":      true, // Indicate this might be cached
			"generated_at": time.Now(),
		},
	})
}

// GetMMRTrajectoryOptimized returns MMR trajectory using optimized service
// GET /api/analytics/v2/mmr?days=30
func (h *OptimizedAnalyticsHandler) GetMMRTrajectoryOptimized(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	// Parse days parameter
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_days",
			"message": "Days must be between 1 and 365",
		})
		return
	}

	// Use optimized async method
	trajectory, err := h.OptimizedService.GetMMRTrajectoryAsync(uid, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to retrieve MMR trajectory",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trajectory,
		"meta": gin.H{
			"days":        days,
			"user_id":     uid,
			"cached":      true,
			"generated_at": time.Now(),
		},
	})
}

// GetRecommendationsOptimized returns AI recommendations using optimized service
// GET /api/analytics/v2/recommendations
func (h *OptimizedAnalyticsHandler) GetRecommendationsOptimized(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	// Use optimized async method
	recommendations, err := h.OptimizedService.GetRecommendationsAsync(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to retrieve recommendations",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    recommendations,
		"meta": gin.H{
			"count":       len(recommendations),
			"user_id":     uid,
			"cached":      true,
			"generated_at": time.Now(),
		},
	})
}

// GetBatchAnalytics returns multiple analytics types in a single request
// POST /api/analytics/v2/batch
func (h *OptimizedAnalyticsHandler) GetBatchAnalytics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	var requestBody struct {
		Requests []string `json:"requests" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate requests
	validRequests := map[string]bool{
		"period_stats_week":  true,
		"period_stats_month": true,
		"mmr_trajectory":     true,
		"recommendations":    true,
	}

	for _, req := range requestBody.Requests {
		if !validRequests[req] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_request_type",
				"message": "Invalid request type: " + req,
				"valid_types": []string{
					"period_stats_week", "period_stats_month", 
					"mmr_trajectory", "recommendations",
				},
			})
			return
		}
	}

	// Use optimized batch method
	results, err := h.OptimizedService.GetBatchAnalytics(uid, requestBody.Requests)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "analytics_error",
			"message": "Failed to retrieve batch analytics",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"meta": gin.H{
			"requests_count": len(requestBody.Requests),
			"results_count":  len(results),
			"user_id":        uid,
			"processed_at":   time.Now(),
		},
	})
}

// InvalidateUserCache invalidates all cache entries for the current user
// POST /api/analytics/v2/cache/invalidate
func (h *OptimizedAnalyticsHandler) InvalidateUserCache(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	err := h.OptimizedService.InvalidateUserCache(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "cache_error",
			"message": "Failed to invalidate user cache",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User cache invalidated successfully",
		"meta": gin.H{
			"user_id":       uid,
			"invalidated_at": time.Now(),
		},
	})
}

// WarmupUserCache pre-calculates and caches analytics for the current user
// POST /api/analytics/v2/cache/warmup
func (h *OptimizedAnalyticsHandler) WarmupUserCache(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "user_not_found",
			"message": "User ID not found in context",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_user_id",
			"message": "Invalid user ID format",
		})
		return
	}

	err := h.OptimizedService.WarmupUserCache(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "cache_error",
			"message": "Failed to warmup user cache",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache warmup initiated successfully",
		"meta": gin.H{
			"user_id":     uid,
			"initiated_at": time.Now(),
		},
	})
}

// GetPerformanceStats returns performance statistics for the optimized service
// GET /api/analytics/v2/performance
func (h *OptimizedAnalyticsHandler) GetPerformanceStats(c *gin.Context) {
	stats := h.OptimizedService.GetPerformanceStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"meta": gin.H{
			"retrieved_at": time.Now(),
		},
	})
}

// HealthCheckOptimized returns the health status of the optimized analytics service
// GET /api/analytics/v2/health
func (h *OptimizedAnalyticsHandler) HealthCheckOptimized(c *gin.Context) {
	perfStats := h.OptimizedService.GetPerformanceStats()
	
	// Determine service health based on stats
	healthy := true
	var issues []string

	if serviceStats, ok := perfStats["service"].(map[string]interface{}); ok {
		if running, ok := serviceStats["running"].(bool); ok && !running {
			healthy = false
			issues = append(issues, "Service not running")
		}
	}

	status := "healthy"
	if !healthy {
		status = "unhealthy"
	}

	response := gin.H{
		"service": "Optimized Analytics Service",
		"status":  status,
		"healthy": healthy,
		"features": gin.H{
			"go_native_analytics": true,
			"redis_cache":         perfStats["cache"],
			"worker_pool":         perfStats["worker_pool"],
			"async_processing":    true,
		},
		"performance": perfStats,
		"timestamp":   time.Now(),
	}

	if len(issues) > 0 {
		response["issues"] = issues
	}

	statusCode := http.StatusOK
	if !healthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}