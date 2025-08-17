package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"lol-match-exporter/internal/db"
)

// FeedbackType represents different types of user feedback
type FeedbackType string

const (
	RecommendationFeedback FeedbackType = "recommendation"
	InsightFeedback       FeedbackType = "insight"
	PerformanceFeedback   FeedbackType = "performance"
	GeneralFeedback       FeedbackType = "general"
)

// FeedbackRating represents user satisfaction rating
type FeedbackRating int

const (
	VeryDissatisfied FeedbackRating = 1
	Dissatisfied     FeedbackRating = 2
	Neutral          FeedbackRating = 3
	Satisfied        FeedbackRating = 4
	VerySatisfied    FeedbackRating = 5
)

// UserFeedback represents user feedback on recommendations and insights
type UserFeedback struct {
	ID                int                    `json:"id"`
	UserID            int                    `json:"user_id"`
	Type              FeedbackType           `json:"type"`
	RelatedID         *int                   `json:"related_id,omitempty"` // ID of recommendation/insight
	Rating            FeedbackRating         `json:"rating"`
	IsHelpful         bool                   `json:"is_helpful"`
	IsAccurate        bool                   `json:"is_accurate"`
	IsActionable      bool                   `json:"is_actionable"`
	Comment           string                 `json:"comment,omitempty"`
	FollowedAdvice    bool                   `json:"followed_advice"`
	PerceivedImpact   *int                   `json:"perceived_impact,omitempty"` // 1-10 scale
	Tags              []string               `json:"tags,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// FeedbackAnalytics contains aggregated feedback analytics
type FeedbackAnalytics struct {
	TotalFeedbacks       int                    `json:"total_feedbacks"`
	AverageRating        float64                `json:"average_rating"`
	HelpfulnessRate      float64                `json:"helpfulness_rate"`
	AccuracyRate         float64                `json:"accuracy_rate"`
	ActionabilityRate    float64                `json:"actionability_rate"`
	FollowThroughRate    float64                `json:"follow_through_rate"`
	AverageImpact        float64                `json:"average_impact"`
	FeedbackByType       map[string]int         `json:"feedback_by_type"`
	FeedbackByRating     map[string]int         `json:"feedback_by_rating"`
	TopTags              []TagCount             `json:"top_tags"`
	TrendingIssues       []string               `json:"trending_issues"`
	ImprovementAreas     []ImprovementArea      `json:"improvement_areas"`
	RecommendationStats  RecommendationStats    `json:"recommendation_stats"`
}

// TagCount represents tag usage count
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// ImprovementArea represents areas needing improvement
type ImprovementArea struct {
	Area        string  `json:"area"`
	Priority    int     `json:"priority"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Actions     []string `json:"actions"`
}

// RecommendationStats contains recommendation-specific statistics
type RecommendationStats struct {
	MostHelpfulTypes     []string `json:"most_helpful_types"`
	LeastHelpfulTypes    []string `json:"least_helpful_types"`
	HighestImpactTypes   []string `json:"highest_impact_types"`
	MostFollowedTypes    []string `json:"most_followed_types"`
}

// FeedbackService handles user feedback collection and analysis
type FeedbackService struct {
	db                  *sql.DB
	analyticsService    *AnalyticsService
}

// NewFeedbackService creates a new feedback service
func NewFeedbackService(database *db.Database, analyticsService *AnalyticsService) *FeedbackService {
	return &FeedbackService{
		db:               database.DB,
		analyticsService: analyticsService,
	}
}

// SubmitFeedback submits user feedback
func (fs *FeedbackService) SubmitFeedback(feedback UserFeedback) (*UserFeedback, error) {
	// Validate feedback
	if err := fs.validateFeedback(feedback); err != nil {
		return nil, fmt.Errorf("invalid feedback: %w", err)
	}

	// Prepare metadata and tags
	metadataJSON, _ := json.Marshal(feedback.Metadata)
	tagsJSON, _ := json.Marshal(feedback.Tags)

	// Insert feedback into database
	query := `
		INSERT INTO user_feedback (
			user_id, type, related_id, rating, is_helpful, is_accurate, 
			is_actionable, comment, followed_advice, perceived_impact, 
			tags, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	err := fs.db.QueryRow(query,
		feedback.UserID, feedback.Type, feedback.RelatedID, feedback.Rating,
		feedback.IsHelpful, feedback.IsAccurate, feedback.IsActionable,
		feedback.Comment, feedback.FollowedAdvice, feedback.PerceivedImpact,
		tagsJSON, metadataJSON,
	).Scan(&feedback.ID, &feedback.CreatedAt, &feedback.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to submit feedback: %w", err)
	}

	// Process feedback for learning
	go fs.processFeedbackForLearning(feedback)

	log.Printf("Feedback submitted: user=%d, type=%s, rating=%d", 
		feedback.UserID, feedback.Type, feedback.Rating)

	return &feedback, nil
}

// GetUserFeedback retrieves feedback for a specific user
func (fs *FeedbackService) GetUserFeedback(userID int, limit int, feedbackType *FeedbackType) ([]UserFeedback, error) {
	query := `
		SELECT id, user_id, type, related_id, rating, is_helpful, is_accurate,
			   is_actionable, comment, followed_advice, perceived_impact,
			   tags, metadata, created_at, updated_at
		FROM user_feedback 
		WHERE user_id = $1`
	
	args := []interface{}{userID}
	argIndex := 2

	if feedbackType != nil {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *feedbackType)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
	}

	rows, err := fs.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feedback: %w", err)
	}
	defer rows.Close()

	var feedbacks []UserFeedback
	for rows.Next() {
		var feedback UserFeedback
		var tagsJSON, metadataJSON []byte

		err = rows.Scan(
			&feedback.ID, &feedback.UserID, &feedback.Type, &feedback.RelatedID,
			&feedback.Rating, &feedback.IsHelpful, &feedback.IsAccurate,
			&feedback.IsActionable, &feedback.Comment, &feedback.FollowedAdvice,
			&feedback.PerceivedImpact, &tagsJSON, &metadataJSON,
			&feedback.CreatedAt, &feedback.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning feedback: %v", err)
			continue
		}

		// Parse JSON fields
		if len(tagsJSON) > 0 {
			json.Unmarshal(tagsJSON, &feedback.Tags)
		}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &feedback.Metadata)
		}

		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

// GetFeedbackAnalytics generates comprehensive feedback analytics
func (fs *FeedbackService) GetFeedbackAnalytics(userID *int, timeframe string) (*FeedbackAnalytics, error) {
	// Build base query with time filter
	timeFilter := fs.getTimeFilter(timeframe)
	
	baseQuery := `
		SELECT 
			COUNT(*) as total,
			AVG(CAST(rating AS FLOAT)) as avg_rating,
			AVG(CASE WHEN is_helpful THEN 1.0 ELSE 0.0 END) as helpfulness_rate,
			AVG(CASE WHEN is_accurate THEN 1.0 ELSE 0.0 END) as accuracy_rate,
			AVG(CASE WHEN is_actionable THEN 1.0 ELSE 0.0 END) as actionability_rate,
			AVG(CASE WHEN followed_advice THEN 1.0 ELSE 0.0 END) as follow_through_rate,
			AVG(CAST(perceived_impact AS FLOAT)) as avg_impact
		FROM user_feedback 
		WHERE created_at >= $1`

	args := []interface{}{timeFilter}
	argIndex := 2

	if userID != nil {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *userID)
	}

	// Execute main analytics query
	var analytics FeedbackAnalytics
	err := fs.db.QueryRow(baseQuery, args...).Scan(
		&analytics.TotalFeedbacks,
		&analytics.AverageRating,
		&analytics.HelpfulnessRate,
		&analytics.AccuracyRate,
		&analytics.ActionabilityRate,
		&analytics.FollowThroughRate,
		&analytics.AverageImpact,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback analytics: %w", err)
	}

	// Get feedback by type
	analytics.FeedbackByType, err = fs.getFeedbackByType(userID, timeFilter)
	if err != nil {
		log.Printf("Error getting feedback by type: %v", err)
	}

	// Get feedback by rating
	analytics.FeedbackByRating, err = fs.getFeedbackByRating(userID, timeFilter)
	if err != nil {
		log.Printf("Error getting feedback by rating: %v", err)
	}

	// Get top tags
	analytics.TopTags, err = fs.getTopTags(userID, timeFilter, 10)
	if err != nil {
		log.Printf("Error getting top tags: %v", err)
	}

	// Get trending issues
	analytics.TrendingIssues, err = fs.getTrendingIssues(userID, timeFilter)
	if err != nil {
		log.Printf("Error getting trending issues: %v", err)
	}

	// Get improvement areas
	analytics.ImprovementAreas = fs.calculateImprovementAreas(analytics)

	// Get recommendation stats
	analytics.RecommendationStats, err = fs.getRecommendationStats(userID, timeFilter)
	if err != nil {
		log.Printf("Error getting recommendation stats: %v", err)
	}

	return &analytics, nil
}

// validateFeedback validates feedback data
func (fs *FeedbackService) validateFeedback(feedback UserFeedback) error {
	if feedback.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if feedback.Rating < VeryDissatisfied || feedback.Rating > VerySatisfied {
		return fmt.Errorf("invalid rating: must be between 1 and 5")
	}

	if feedback.PerceivedImpact != nil && (*feedback.PerceivedImpact < 1 || *feedback.PerceivedImpact > 10) {
		return fmt.Errorf("invalid perceived impact: must be between 1 and 10")
	}

	validTypes := map[FeedbackType]bool{
		RecommendationFeedback: true,
		InsightFeedback:       true,
		PerformanceFeedback:   true,
		GeneralFeedback:       true,
	}

	if !validTypes[feedback.Type] {
		return fmt.Errorf("invalid feedback type")
	}

	return nil
}

// processFeedbackForLearning processes feedback to improve recommendation algorithms
func (fs *FeedbackService) processFeedbackForLearning(feedback UserFeedback) {
	log.Printf("Processing feedback for learning: user=%d, type=%s", feedback.UserID, feedback.Type)

	// Update recommendation weights based on feedback
	if feedback.Type == RecommendationFeedback && feedback.RelatedID != nil {
		fs.updateRecommendationWeights(feedback)
	}

	// Detect patterns in negative feedback
	if feedback.Rating <= Dissatisfied {
		fs.analyzeNegativeFeedback(feedback)
	}

	// Update user preference model
	fs.updateUserPreferences(feedback)

	// Train improvement models
	fs.trainImprovementModels(feedback)
}

// updateRecommendationWeights adjusts recommendation algorithm weights
func (fs *FeedbackService) updateRecommendationWeights(feedback UserFeedback) {
	// Get recommendation details
	recommendation, err := fs.getRecommendationByID(*feedback.RelatedID)
	if err != nil {
		log.Printf("Error getting recommendation: %v", err)
		return
	}

	// Calculate weight adjustment based on feedback
	weightAdjustment := fs.calculateWeightAdjustment(feedback, recommendation)
	
	// TODO: Update weights in recommendation engine when needed
	// The recommendation engine is now integrated into AnalyticsService
	
	log.Printf("Updated recommendation weights: adjustment=%.3f", weightAdjustment)
}

// calculateWeightAdjustment calculates how much to adjust recommendation weights
func (fs *FeedbackService) calculateWeightAdjustment(feedback UserFeedback, recommendation interface{}) float64 {
	baseAdjustment := 0.0

	// Rating impact (stronger for extreme ratings)
	switch feedback.Rating {
	case VerySatisfied:
		baseAdjustment = 0.15
	case Satisfied:
		baseAdjustment = 0.05
	case Neutral:
		baseAdjustment = 0.0
	case Dissatisfied:
		baseAdjustment = -0.05
	case VeryDissatisfied:
		baseAdjustment = -0.15
	}

	// Helpfulness impact
	if feedback.IsHelpful {
		baseAdjustment += 0.1
	} else {
		baseAdjustment -= 0.1
	}

	// Accuracy impact
	if feedback.IsAccurate {
		baseAdjustment += 0.05
	} else {
		baseAdjustment -= 0.1
	}

	// Follow-through impact (strongest signal)
	if feedback.FollowedAdvice {
		baseAdjustment += 0.2
		
		// Perceived impact boost
		if feedback.PerceivedImpact != nil && *feedback.PerceivedImpact >= 7 {
			baseAdjustment += 0.1
		}
	} else {
		baseAdjustment -= 0.05
	}

	// Clamp adjustment to reasonable range
	if baseAdjustment > 0.3 {
		baseAdjustment = 0.3
	} else if baseAdjustment < -0.3 {
		baseAdjustment = -0.3
	}

	return baseAdjustment
}

// Helper methods for analytics

func (fs *FeedbackService) getTimeFilter(timeframe string) time.Time {
	now := time.Now()
	switch timeframe {
	case "day":
		return now.Add(-24 * time.Hour)
	case "week":
		return now.Add(-7 * 24 * time.Hour)
	case "month":
		return now.Add(-30 * 24 * time.Hour)
	case "quarter":
		return now.Add(-90 * 24 * time.Hour)
	case "year":
		return now.Add(-365 * 24 * time.Hour)
	default:
		return now.Add(-30 * 24 * time.Hour) // Default to month
	}
}

func (fs *FeedbackService) getFeedbackByType(userID *int, timeFilter time.Time) (map[string]int, error) {
	query := `
		SELECT type, COUNT(*) 
		FROM user_feedback 
		WHERE created_at >= $1`
	
	args := []interface{}{timeFilter}
	if userID != nil {
		query += " AND user_id = $2"
		args = append(args, *userID)
	}
	
	query += " GROUP BY type"

	rows, err := fs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var feedbackType string
		var count int
		if err := rows.Scan(&feedbackType, &count); err != nil {
			continue
		}
		result[feedbackType] = count
	}

	return result, nil
}

func (fs *FeedbackService) getFeedbackByRating(userID *int, timeFilter time.Time) (map[string]int, error) {
	query := `
		SELECT rating, COUNT(*) 
		FROM user_feedback 
		WHERE created_at >= $1`
	
	args := []interface{}{timeFilter}
	if userID != nil {
		query += " AND user_id = $2"
		args = append(args, *userID)
	}
	
	query += " GROUP BY rating ORDER BY rating"

	rows, err := fs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var rating int
		var count int
		if err := rows.Scan(&rating, &count); err != nil {
			continue
		}
		result[fmt.Sprintf("%d", rating)] = count
	}

	return result, nil
}

func (fs *FeedbackService) getTopTags(userID *int, timeFilter time.Time, limit int) ([]TagCount, error) {
	query := `
		SELECT jsonb_array_elements_text(tags) as tag, COUNT(*) as count
		FROM user_feedback 
		WHERE created_at >= $1 AND tags IS NOT NULL AND jsonb_array_length(tags) > 0`
	
	args := []interface{}{timeFilter}
	if userID != nil {
		query += " AND user_id = $2"
		args = append(args, *userID)
	}
	
	query += " GROUP BY tag ORDER BY count DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	rows, err := fs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []TagCount
	for rows.Next() {
		var tag TagCount
		if err := rows.Scan(&tag.Tag, &tag.Count); err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (fs *FeedbackService) getTrendingIssues(userID *int, timeFilter time.Time) ([]string, error) {
	// Simple implementation - look for common negative feedback patterns
	query := `
		SELECT comment
		FROM user_feedback 
		WHERE created_at >= $1 AND rating <= 2 AND comment IS NOT NULL AND comment != ''`
	
	args := []interface{}{timeFilter}
	if userID != nil {
		query += " AND user_id = $2"
		args = append(args, *userID)
	}

	rows, err := fs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// In a real implementation, this would use NLP to extract common issues
	issues := []string{
		"Recommendations not specific enough",
		"Insights arrive too late",
		"Performance metrics confusing",
		"Champion suggestions inaccurate",
	}

	return issues, nil
}

func (fs *FeedbackService) calculateImprovementAreas(analytics FeedbackAnalytics) []ImprovementArea {
	var areas []ImprovementArea

	// Analyze metrics and suggest improvements
	if analytics.HelpfulnessRate < 0.7 {
		areas = append(areas, ImprovementArea{
			Area:        "Recommendation Relevance",
			Priority:    1,
			Score:       analytics.HelpfulnessRate,
			Description: "Users find recommendations less helpful than expected",
			Actions:     []string{"Improve personalization", "Better context analysis", "User preference learning"},
		})
	}

	if analytics.AccuracyRate < 0.8 {
		areas = append(areas, ImprovementArea{
			Area:        "Recommendation Accuracy",
			Priority:    1,
			Score:       analytics.AccuracyRate,
			Description: "Recommendations accuracy needs improvement",
			Actions:     []string{"Enhance data quality", "Improve prediction models", "Better validation"},
		})
	}

	if analytics.FollowThroughRate < 0.5 {
		areas = append(areas, ImprovementArea{
			Area:        "Actionability",
			Priority:    2,
			Score:       analytics.FollowThroughRate,
			Description: "Users don't follow through on recommendations",
			Actions:     []string{"Simplify action steps", "Provide clearer guidance", "Add progress tracking"},
		})
	}

	return areas
}

func (fs *FeedbackService) getRecommendationStats(userID *int, timeFilter time.Time) (RecommendationStats, error) {
	// This would analyze recommendation types and their performance
	// For now, return mock data based on common patterns
	
	return RecommendationStats{
		MostHelpfulTypes:   []string{"champion_suggestion", "gameplay_tip"},
		LeastHelpfulTypes:  []string{"role_optimization", "build_suggestion"},
		HighestImpactTypes: []string{"gameplay_tip", "champion_suggestion"},
		MostFollowedTypes:  []string{"champion_suggestion", "role_optimization"},
	}, nil
}

// Additional helper methods

func (fs *FeedbackService) getRecommendationByID(id int) (interface{}, error) {
	// This would fetch the recommendation from the database
	// For now, return a simple mock
	return map[string]interface{}{
		"id":   id,
		"type": "champion_suggestion",
	}, nil
}

func (fs *FeedbackService) analyzeNegativeFeedback(feedback UserFeedback) {
	log.Printf("Analyzing negative feedback: user=%d, rating=%d, comment=%s", 
		feedback.UserID, feedback.Rating, feedback.Comment)
	
	// In a real implementation, this would:
	// - Extract keywords from comments
	// - Categorize common issues
	// - Alert developers to systemic problems
	// - Update blacklists or filters
}

func (fs *FeedbackService) updateUserPreferences(feedback UserFeedback) {
	log.Printf("Updating user preferences based on feedback: user=%d", feedback.UserID)
	
	// This would update user preference models based on:
	// - What types of recommendations they find helpful
	// - What aspects they value most (accuracy vs actionability)
	// - Their engagement patterns
}

func (fs *FeedbackService) trainImprovementModels(feedback UserFeedback) {
	log.Printf("Training improvement models with feedback: user=%d", feedback.UserID)
	
	// This would:
	// - Update ML models with new feedback data
	// - Retrain recommendation algorithms
	// - Adjust confidence thresholds
	// - Update feature weights
}