package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/match"
	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Match Processing Service
// Comprehensive match data processing with analytics integration

// MatchProcessingService handles end-to-end match processing workflows
type MatchProcessingService struct {
	riotService     *RiotService
	analyticsEngine *analytics.AnalyticsEngine
	matchAnalyzer   *match.MatchAnalyzer

	// Processing configuration
	config *MatchProcessingConfig

	// Processing state
	processingQueue chan *MatchProcessingJob
	workers         int
	shutdown        chan bool
	wg              sync.WaitGroup
}

// MatchProcessingConfig contains service configuration
type MatchProcessingConfig struct {
	// Worker configuration
	WorkerCount       int           `json:"worker_count"`
	QueueSize         int           `json:"queue_size"`
	ProcessingTimeout time.Duration `json:"processing_timeout"`
	RetryAttempts     int           `json:"retry_attempts"`
	RetryDelay        time.Duration `json:"retry_delay"`

	// Analysis configuration
	DefaultAnalysisDepth string `json:"default_analysis_depth"`
	EnablePhaseAnalysis  bool   `json:"enable_phase_analysis"`
	EnableKeyMoments     bool   `json:"enable_key_moments"`
	EnableTeamAnalysis   bool   `json:"enable_team_analysis"`

	// Performance settings
	CacheResults    bool          `json:"cache_results"`
	CacheExpiration time.Duration `json:"cache_expiration"`
	BatchSize       int           `json:"batch_size"`

	// Rate limiting
	ProcessingRateLimit int `json:"processing_rate_limit"` // matches per minute
}

// MatchProcessingJob represents a match processing job
type MatchProcessingJob struct {
	ID          string                         `json:"id"`
	MatchID     string                         `json:"match_id"`
	PlayerPUUID string                         `json:"player_puuid"`
	Priority    int                            `json:"priority"` // 1-10 scale
	Options     *MatchProcessingOptions        `json:"options"`
	CreatedAt   time.Time                      `json:"created_at"`
	StartedAt   *time.Time                     `json:"started_at"`
	CompletedAt *time.Time                     `json:"completed_at"`
	Status      string                         `json:"status"` // "pending", "processing", "completed", "failed"
	Result      *MatchProcessingResult         `json:"result"`
	Error       string                         `json:"error"`
	RetryCount  int                            `json:"retry_count"`
	Callbacks   []func(*MatchProcessingResult) `json:"-"`
}

// MatchProcessingOptions contains processing options
type MatchProcessingOptions struct {
	AnalysisDepth           string   `json:"analysis_depth"`
	IncludePhaseAnalysis    bool     `json:"include_phase_analysis"`
	IncludeKeyMoments       bool     `json:"include_key_moments"`
	IncludeTeamAnalysis     bool     `json:"include_team_analysis"`
	IncludeOpponentAnalysis bool     `json:"include_opponent_analysis"`
	FocusAreas              []string `json:"focus_areas"`
	CompareWithAverage      bool     `json:"compare_with_average"`
	GenerateInsights        bool     `json:"generate_insights"`
	CacheResult             bool     `json:"cache_result"`
}

// MatchProcessingResult contains the complete processing result
type MatchProcessingResult struct {
	JobID             string                     `json:"job_id"`
	MatchAnalysis     *match.MatchAnalysisResult `json:"match_analysis"`
	AnalyticsData     *analytics.AnalyticsResult `json:"analytics_data"`
	ProcessingMetrics *ProcessingMetrics         `json:"processing_metrics"`
	Status            string                     `json:"status"`
	ProcessedAt       time.Time                  `json:"processed_at"`
}

// ProcessingMetrics contains metrics about the processing job
type ProcessingMetrics struct {
	ProcessingTime    time.Duration `json:"processing_time"`
	QueueTime         time.Duration `json:"queue_time"`
	AnalysisTime      time.Duration `json:"analysis_time"`
	DataRetrievalTime time.Duration `json:"data_retrieval_time"`
	RetryCount        int           `json:"retry_count"`
	WorkerID          string        `json:"worker_id"`
}

// MatchBatchProcessingRequest contains batch processing parameters
type MatchBatchProcessingRequest struct {
	PlayerPUUID string                  `json:"player_puuid"`
	MatchIDs    []string                `json:"match_ids"`
	Options     *MatchProcessingOptions `json:"options"`
	Priority    int                     `json:"priority"`
	CallbackURL string                  `json:"callback_url"`
}

// MatchBatchProcessingResult contains batch processing results
type MatchBatchProcessingResult struct {
	BatchID         string                   `json:"batch_id"`
	TotalMatches    int                      `json:"total_matches"`
	ProcessedCount  int                      `json:"processed_count"`
	SuccessCount    int                      `json:"success_count"`
	FailureCount    int                      `json:"failure_count"`
	Results         []*MatchProcessingResult `json:"results"`
	StartedAt       time.Time                `json:"started_at"`
	CompletedAt     *time.Time               `json:"completed_at"`
	Status          string                   `json:"status"`
	ProcessingStats *BatchProcessingStats    `json:"processing_stats"`
}

// BatchProcessingStats contains batch processing statistics
type BatchProcessingStats struct {
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	FastestProcess        time.Duration `json:"fastest_process"`
	SlowestProcess        time.Duration `json:"slowest_process"`
	TotalProcessingTime   time.Duration `json:"total_processing_time"`
	SuccessRate           float64       `json:"success_rate"`
	RetryRate             float64       `json:"retry_rate"`
}

// NewMatchProcessingService creates a new match processing service
func NewMatchProcessingService(
	riotService *RiotService,
	analyticsEngine *analytics.AnalyticsEngine,
	config *MatchProcessingConfig,
) *MatchProcessingService {
	if config == nil {
		config = DefaultMatchProcessingConfig()
	}

	// Create match analyzer
	analyzerConfig := match.DefaultMatchAnalysisConfig()
	analyzerConfig.EnablePhaseAnalysis = config.EnablePhaseAnalysis
	analyzerConfig.EnableKeyMomentDetection = config.EnableKeyMoments
	analyzerConfig.EnableTeamAnalysis = config.EnableTeamAnalysis

	matchAnalyzer := match.NewMatchAnalyzer(analyzerConfig, analyticsEngine)

	service := &MatchProcessingService{
		riotService:     riotService,
		analyticsEngine: analyticsEngine,
		matchAnalyzer:   matchAnalyzer,
		config:          config,
		processingQueue: make(chan *MatchProcessingJob, config.QueueSize),
		workers:         config.WorkerCount,
		shutdown:        make(chan bool),
	}

	// Start worker goroutines
	service.startWorkers()

	return service
}

// ProcessMatch processes a single match asynchronously
func (s *MatchProcessingService) ProcessMatch(ctx context.Context, matchID, playerPUUID string, options *MatchProcessingOptions) (*MatchProcessingJob, error) {
	if matchID == "" || playerPUUID == "" {
		return nil, fmt.Errorf("match_id and player_puuid are required")
	}

	if options == nil {
		options = s.defaultProcessingOptions()
	}

	job := &MatchProcessingJob{
		ID:          fmt.Sprintf("match_%s_%s_%d", matchID, playerPUUID, time.Now().UnixNano()),
		MatchID:     matchID,
		PlayerPUUID: playerPUUID,
		Priority:    5, // Default priority
		Options:     options,
		CreatedAt:   time.Now(),
		Status:      "pending",
	}

	select {
	case s.processingQueue <- job:
		log.Printf("Match processing job queued: %s", job.ID)
		return job, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("processing queue full or context cancelled")
	default:
		return nil, fmt.Errorf("processing queue is full")
	}
}

// ProcessMatchSync processes a match synchronously
func (s *MatchProcessingService) ProcessMatchSync(ctx context.Context, matchID, playerPUUID string, options *MatchProcessingOptions) (*MatchProcessingResult, error) {
	job, err := s.ProcessMatch(ctx, matchID, playerPUUID, options)
	if err != nil {
		return nil, err
	}

	// Wait for completion with timeout
	timeout := time.After(s.config.ProcessingTimeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("processing timeout for job %s", job.ID)
		case <-ticker.C:
			if job.Status == "completed" {
				return job.Result, nil
			}
			if job.Status == "failed" {
				return nil, fmt.Errorf("processing failed: %s", job.Error)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// ProcessMatchBatch processes multiple matches as a batch
func (s *MatchProcessingService) ProcessMatchBatch(ctx context.Context, request *MatchBatchProcessingRequest) (*MatchBatchProcessingResult, error) {
	if len(request.MatchIDs) == 0 {
		return nil, fmt.Errorf("no match IDs provided")
	}

	batchResult := &MatchBatchProcessingResult{
		BatchID:      fmt.Sprintf("batch_%d", time.Now().UnixNano()),
		TotalMatches: len(request.MatchIDs),
		Results:      make([]*MatchProcessingResult, 0, len(request.MatchIDs)),
		StartedAt:    time.Now(),
		Status:       "processing",
	}

	// Process matches in batches to avoid overwhelming the system
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, s.config.BatchSize)
	resultsChan := make(chan *MatchProcessingResult, len(request.MatchIDs))
	errorsChan := make(chan error, len(request.MatchIDs))

	for _, matchID := range request.MatchIDs {
		wg.Add(1)
		go func(mID string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			result, err := s.ProcessMatchSync(ctx, mID, request.PlayerPUUID, request.Options)
			if err != nil {
				errorsChan <- err
				return
			}
			resultsChan <- result
		}(matchID)
	}

	// Wait for all jobs to complete
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorsChan)
	}()

	// Collect results
	for result := range resultsChan {
		batchResult.Results = append(batchResult.Results, result)
		batchResult.ProcessedCount++
		batchResult.SuccessCount++
	}

	// Count errors
	for range errorsChan {
		batchResult.ProcessedCount++
		batchResult.FailureCount++
	}

	// Finalize batch result
	completedAt := time.Now()
	batchResult.CompletedAt = &completedAt
	batchResult.Status = "completed"
	batchResult.ProcessingStats = s.calculateBatchStats(batchResult)

	return batchResult, nil
}

// GetJobStatus returns the status of a processing job
func (s *MatchProcessingService) GetJobStatus(jobID string) (*MatchProcessingJob, error) {
	// In a real implementation, this would query a database or cache
	// For now, we'll return a mock response
	return nil, fmt.Errorf("job status tracking not yet implemented")
}

// Worker implementation

func (s *MatchProcessingService) startWorkers() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(fmt.Sprintf("worker-%d", i))
	}
}

func (s *MatchProcessingService) worker(workerID string) {
	defer s.wg.Done()

	for {
		select {
		case job := <-s.processingQueue:
			s.processJob(job, workerID)
		case <-s.shutdown:
			return
		}
	}
}

func (s *MatchProcessingService) processJob(job *MatchProcessingJob, workerID string) {
	startTime := time.Now()
	job.Status = "processing"
	job.StartedAt = &startTime

	ctx, cancel := context.WithTimeout(context.Background(), s.config.ProcessingTimeout)
	defer cancel()

	// Process the match
	result, err := s.doProcessMatch(ctx, job, workerID)

	completedAt := time.Now()
	job.CompletedAt = &completedAt

	if err != nil {
		job.Error = err.Error()
		if job.RetryCount < s.config.RetryAttempts {
			job.RetryCount++
			job.Status = "pending"

			// Retry after delay
			go func() {
				time.Sleep(s.config.RetryDelay)
				select {
				case s.processingQueue <- job:
					log.Printf("Retrying job %s (attempt %d)", job.ID, job.RetryCount+1)
				default:
					log.Printf("Failed to queue retry for job %s", job.ID)
					job.Status = "failed"
				}
			}()
			return
		}
		job.Status = "failed"
		log.Printf("Job failed permanently: %s - %s", job.ID, err.Error())
	} else {
		job.Status = "completed"
		job.Result = result
		log.Printf("Job completed successfully: %s", job.ID)
	}

	// Execute callbacks
	if job.Result != nil {
		for _, callback := range job.Callbacks {
			go func(cb func(*MatchProcessingResult)) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Callback panic for job %s: %v", job.ID, r)
					}
				}()
				cb(job.Result)
			}(callback)
		}
	}
}

func (s *MatchProcessingService) doProcessMatch(ctx context.Context, job *MatchProcessingJob, workerID string) (*MatchProcessingResult, error) {
	startTime := time.Now()

	// Retrieve match data from Riot API
	dataRetrievalStart := time.Now()
	matchData, err := s.riotService.GetMatchByID(ctx, job.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve match data: %w", err)
	}
	dataRetrievalTime := time.Since(dataRetrievalStart)

	// Perform match analysis
	analysisStart := time.Now()
	analysisRequest := &match.MatchAnalysisRequest{
		Match:                   matchData,
		PlayerPUUID:             job.PlayerPUUID,
		AnalysisDepth:           job.Options.AnalysisDepth,
		IncludePhaseAnalysis:    job.Options.IncludePhaseAnalysis,
		IncludeKeyMoments:       job.Options.IncludeKeyMoments,
		IncludeTeamAnalysis:     job.Options.IncludeTeamAnalysis,
		IncludeOpponentAnalysis: job.Options.IncludeOpponentAnalysis,
		CompareWithAverage:      job.Options.CompareWithAverage,
		FocusAreas:              job.Options.FocusAreas,
	}

	matchAnalysis, err := s.matchAnalyzer.AnalyzeMatch(ctx, analysisRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze match: %w", err)
	}
	analysisTime := time.Since(analysisStart)

	// Generate analytics data if analytics engine is available
	var analyticsData *analytics.AnalyticsResult
	if s.analyticsEngine != nil && job.Options.GenerateInsights {
		analyticsData, _ = s.generateAnalyticsData(ctx, matchData, matchAnalysis)
	}

	// Create processing result
	result := &MatchProcessingResult{
		JobID:         job.ID,
		MatchAnalysis: matchAnalysis,
		AnalyticsData: analyticsData,
		ProcessingMetrics: &ProcessingMetrics{
			ProcessingTime:    time.Since(startTime),
			QueueTime:         startTime.Sub(job.CreatedAt),
			AnalysisTime:      analysisTime,
			DataRetrievalTime: dataRetrievalTime,
			RetryCount:        job.RetryCount,
			WorkerID:          workerID,
		},
		Status:      "completed",
		ProcessedAt: time.Now(),
	}

	return result, nil
}

// Helper methods

func (s *MatchProcessingService) defaultProcessingOptions() *MatchProcessingOptions {
	return &MatchProcessingOptions{
		AnalysisDepth:           s.config.DefaultAnalysisDepth,
		IncludePhaseAnalysis:    s.config.EnablePhaseAnalysis,
		IncludeKeyMoments:       s.config.EnableKeyMoments,
		IncludeTeamAnalysis:     s.config.EnableTeamAnalysis,
		IncludeOpponentAnalysis: false,
		CompareWithAverage:      true,
		GenerateInsights:        true,
		CacheResult:             s.config.CacheResults,
	}
}

func (s *MatchProcessingService) generateAnalyticsData(ctx context.Context, matchData *riot.Match, analysis *match.MatchAnalysisResult) (*analytics.AnalyticsResult, error) {
	// This would integrate with the analytics engine to generate additional insights
	// For now, return a placeholder
	return &analytics.AnalyticsResult{
		PlayerID:   analysis.PlayerPUUID,
		MatchID:    analysis.MatchID,
		Timestamp:  time.Now(),
		Metrics:    make(map[string]interface{}),
		Insights:   []string{"Advanced analytics integration pending"},
		Confidence: 0.85,
	}, nil
}

func (s *MatchProcessingService) calculateBatchStats(batchResult *MatchBatchProcessingResult) *BatchProcessingStats {
	if len(batchResult.Results) == 0 {
		return &BatchProcessingStats{}
	}

	var totalTime, fastestTime, slowestTime time.Duration
	fastestTime = time.Hour // Initialize with a large value

	for _, result := range batchResult.Results {
		if result.ProcessingMetrics != nil {
			processingTime := result.ProcessingMetrics.ProcessingTime
			totalTime += processingTime

			if processingTime < fastestTime {
				fastestTime = processingTime
			}
			if processingTime > slowestTime {
				slowestTime = processingTime
			}
		}
	}

	averageTime := totalTime / time.Duration(len(batchResult.Results))
	successRate := float64(batchResult.SuccessCount) / float64(batchResult.TotalMatches)

	return &BatchProcessingStats{
		AverageProcessingTime: averageTime,
		FastestProcess:        fastestTime,
		SlowestProcess:        slowestTime,
		TotalProcessingTime:   totalTime,
		SuccessRate:           successRate,
		RetryRate:             0.0, // Would be calculated based on retry statistics
	}
}

// Shutdown gracefully shuts down the match processing service
func (s *MatchProcessingService) Shutdown(ctx context.Context) error {
	log.Println("Shutting down match processing service...")

	close(s.shutdown)

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Match processing service shut down successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout exceeded")
	}
}

// DefaultMatchProcessingConfig returns default configuration
func DefaultMatchProcessingConfig() *MatchProcessingConfig {
	return &MatchProcessingConfig{
		WorkerCount:          5,
		QueueSize:            1000,
		ProcessingTimeout:    30 * time.Second,
		RetryAttempts:        3,
		RetryDelay:           5 * time.Second,
		DefaultAnalysisDepth: "standard",
		EnablePhaseAnalysis:  true,
		EnableKeyMoments:     true,
		EnableTeamAnalysis:   true,
		CacheResults:         true,
		CacheExpiration:      24 * time.Hour,
		BatchSize:            10,
		ProcessingRateLimit:  60, // 60 matches per minute
	}
}
