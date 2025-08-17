package services

// AnalyticsServiceWrapper wraps AnalyticsService to implement the worker interface
type AnalyticsServiceWrapper struct {
	service *AnalyticsService
}

// NewAnalyticsServiceWrapper creates a wrapper for the analytics service
func NewAnalyticsServiceWrapper(service *AnalyticsService) *AnalyticsServiceWrapper {
	return &AnalyticsServiceWrapper{service: service}
}

// GetPeriodStats implements the worker interface
func (asw *AnalyticsServiceWrapper) GetPeriodStats(userID int, period string) (interface{}, error) {
	return asw.service.GetPeriodStats(userID, period)
}

// GetMMRTrajectory implements the worker interface
func (asw *AnalyticsServiceWrapper) GetMMRTrajectory(userID int, days int) (interface{}, error) {
	return asw.service.GetMMRTrajectory(userID, days)
}

// GetRecommendations implements the worker interface
func (asw *AnalyticsServiceWrapper) GetRecommendations(userID int) (interface{}, error) {
	return asw.service.GetRecommendations(userID)
}

// GetChampionAnalysis implements the worker interface
func (asw *AnalyticsServiceWrapper) GetChampionAnalysis(userID int, championID int, period string) (interface{}, error) {
	return asw.service.GetChampionAnalysis(userID, championID, period)
}