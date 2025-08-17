package workers

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"lol-match-exporter/internal/cache"
)

// AnalyticsService interface to avoid import cycles
type AnalyticsService interface {
	GetPeriodStats(userID int, period string) (interface{}, error)
	GetMMRTrajectory(userID int, days int) (interface{}, error)
	GetRecommendations(userID int) (interface{}, error)
	GetChampionAnalysis(userID int, championID int, period string) (interface{}, error)
}


// AnalyticsTask represents a task to be processed by analytics workers
type AnalyticsTask struct {
	ID       string
	Type     TaskType
	UserID   int
	Data     map[string]interface{}
	Priority int
	Created  time.Time
	Retry    int
	MaxRetry int
}

// TaskType defines the type of analytics task
type TaskType string

const (
	TaskPeriodStats      TaskType = "period_stats"
	TaskMMRTrajectory    TaskType = "mmr_trajectory"
	TaskRecommendations  TaskType = "recommendations"
	TaskChampionAnalysis TaskType = "champion_analysis"
	TaskCacheWarmup      TaskType = "cache_warmup"
	TaskCacheInvalidate  TaskType = "cache_invalidate"
)

// TaskResult represents the result of a processed task
type TaskResult struct {
	TaskID    string
	Success   bool
	Error     error
	Data      interface{}
	Duration  time.Duration
	Timestamp time.Time
}

// WorkerPool manages a pool of goroutines for processing analytics tasks
type WorkerPool struct {
	taskQueue      chan AnalyticsTask
	resultQueue    chan TaskResult
	workers        []*Worker
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	running        bool
	mutex          sync.RWMutex
	
	// Services
	analyticsService AnalyticsService
	cacheService     *cache.CacheService
	
	// Configuration
	maxWorkers    int
	maxQueueSize  int
	
	// Statistics
	stats WorkerPoolStats
}

// WorkerPoolStats tracks pool performance metrics
type WorkerPoolStats struct {
	TasksProcessed   int64
	TasksSucceeded   int64
	TasksFailed      int64
	TasksInQueue     int64
	AverageTaskTime  time.Duration
	TotalProcessTime time.Duration
	WorkersActive    int
	QueueUtilization float64
	mutex            sync.RWMutex
}

// Worker represents a single worker goroutine
type Worker struct {
	id       int
	pool     *WorkerPool
	taskChan chan AnalyticsTask
	active   bool
	mutex    sync.RWMutex
}

// NewWorkerPool creates a new analytics worker pool
func NewWorkerPool(analyticsService AnalyticsService, cacheService *cache.CacheService) *WorkerPool {
	// Default to number of CPU cores, but allow override
	maxWorkers := runtime.NumCPU()
	if maxWorkers < 2 {
		maxWorkers = 2
	}
	if maxWorkers > 10 {
		maxWorkers = 10 // Cap at 10 workers
	}

	maxQueueSize := maxWorkers * 100 // 100 tasks per worker

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		taskQueue:        make(chan AnalyticsTask, maxQueueSize),
		resultQueue:      make(chan TaskResult, maxQueueSize/10),
		workers:          make([]*Worker, 0, maxWorkers),
		ctx:              ctx,
		cancel:           cancel,
		analyticsService: analyticsService,
		cacheService:     cacheService,
		maxWorkers:       maxWorkers,
		maxQueueSize:     maxQueueSize,
		stats:            WorkerPoolStats{},
	}
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start() error {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.running {
		return fmt.Errorf("worker pool is already running")
	}

	log.Printf("🚀 Starting analytics worker pool with %d workers", wp.maxWorkers)

	// Create and start workers
	for i := 0; i < wp.maxWorkers; i++ {
		worker := &Worker{
			id:       i + 1,
			pool:     wp,
			taskChan: wp.taskQueue,
		}
		wp.workers = append(wp.workers, worker)

		wp.wg.Add(1)
		go worker.start()
	}

	// Start result processor
	wp.wg.Add(1)
	go wp.processResults()

	// Start stats updater
	wp.wg.Add(1)
	go wp.updateStats()

	wp.running = true
	log.Printf("✅ Analytics worker pool started successfully")
	
	return nil
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() error {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if !wp.running {
		return fmt.Errorf("worker pool is not running")
	}

	log.Println("🛑 Stopping analytics worker pool...")

	// Cancel context to signal workers to stop
	wp.cancel()

	// Close task queue to prevent new tasks
	close(wp.taskQueue)

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close result queue
	close(wp.resultQueue)

	wp.running = false
	log.Println("✅ Analytics worker pool stopped")

	return nil
}

// SubmitTask adds a task to the worker pool queue
func (wp *WorkerPool) SubmitTask(task AnalyticsTask) error {
	wp.mutex.RLock()
	if !wp.running {
		wp.mutex.RUnlock()
		return fmt.Errorf("worker pool is not running")
	}
	wp.mutex.RUnlock()

	// Set default values
	if task.ID == "" {
		task.ID = fmt.Sprintf("%s_%d_%d", task.Type, task.UserID, time.Now().UnixNano())
	}
	if task.Created.IsZero() {
		task.Created = time.Now()
	}
	if task.MaxRetry == 0 {
		task.MaxRetry = 3
	}

	select {
	case wp.taskQueue <- task:
		wp.stats.incrementTasksInQueue()
		return nil
	case <-time.After(time.Second * 5):
		return fmt.Errorf("worker pool queue is full, task rejected")
	}
}

// SubmitHighPriorityTask submits a task with high priority
func (wp *WorkerPool) SubmitHighPriorityTask(taskType TaskType, userID int, data map[string]interface{}) error {
	task := AnalyticsTask{
		Type:     taskType,
		UserID:   userID,
		Data:     data,
		Priority: 1, // High priority
	}
	return wp.SubmitTask(task)
}

// SubmitLowPriorityTask submits a task with low priority
func (wp *WorkerPool) SubmitLowPriorityTask(taskType TaskType, userID int, data map[string]interface{}) error {
	task := AnalyticsTask{
		Type:     taskType,
		UserID:   userID,
		Data:     data,
		Priority: 3, // Low priority
	}
	return wp.SubmitTask(task)
}

// GetStats returns current worker pool statistics
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.stats.mutex.RLock()
	defer wp.stats.mutex.RUnlock()
	
	// Update queue utilization
	queueLength := float64(len(wp.taskQueue))
	wp.stats.QueueUtilization = (queueLength / float64(wp.maxQueueSize)) * 100

	// Count active workers
	activeWorkers := 0
	for _, worker := range wp.workers {
		if worker.isActive() {
			activeWorkers++
		}
	}
	wp.stats.WorkersActive = activeWorkers

	return wp.stats
}

// IsRunning returns true if the worker pool is currently running
func (wp *WorkerPool) IsRunning() bool {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	return wp.running
}

// Worker methods

// start begins the worker's task processing loop
func (w *Worker) start() {
	defer w.pool.wg.Done()
	
	log.Printf("📋 Analytics worker %d started", w.id)

	for {
		select {
		case task, ok := <-w.taskChan:
			if !ok {
				log.Printf("📋 Analytics worker %d stopping (channel closed)", w.id)
				return
			}
			
			w.setActive(true)
			result := w.processTask(task)
			w.setActive(false)
			
			// Send result
			select {
			case w.pool.resultQueue <- result:
			case <-w.pool.ctx.Done():
				return
			}

		case <-w.pool.ctx.Done():
			log.Printf("📋 Analytics worker %d stopping (context cancelled)", w.id)
			return
		}
	}
}

// processTask processes a single analytics task
func (w *Worker) processTask(task AnalyticsTask) TaskResult {
	startTime := time.Now()
	
	result := TaskResult{
		TaskID:    task.ID,
		Timestamp: startTime,
	}

	log.Printf("🔄 Worker %d processing task %s (type: %s, user: %d)", w.id, task.ID, task.Type, task.UserID)

	// Process based on task type
	switch task.Type {
	case TaskPeriodStats:
		result.Data, result.Error = w.processPeriodStats(task)
	case TaskMMRTrajectory:
		result.Data, result.Error = w.processMMRTrajectory(task)
	case TaskRecommendations:
		result.Data, result.Error = w.processRecommendations(task)
	case TaskChampionAnalysis:
		result.Data, result.Error = w.processChampionAnalysis(task)
	case TaskCacheWarmup:
		result.Error = w.processCacheWarmup(task)
	case TaskCacheInvalidate:
		result.Error = w.processCacheInvalidate(task)
	default:
		result.Error = fmt.Errorf("unknown task type: %s", task.Type)
	}

	result.Duration = time.Since(startTime)
	result.Success = result.Error == nil

	if result.Success {
		log.Printf("✅ Worker %d completed task %s in %v", w.id, task.ID, result.Duration)
	} else {
		log.Printf("❌ Worker %d failed task %s: %v", w.id, task.ID, result.Error)
	}

	return result
}

// Task processors

func (w *Worker) processPeriodStats(task AnalyticsTask) (interface{}, error) {
	period, ok := task.Data["period"].(string)
	if !ok {
		period = "week"
	}

	// Check cache first
	cacheKey := cache.AnalyticsCacheKey(task.UserID, period, "period_stats")
	if w.pool.cacheService.IsEnabled() {
		var cachedData interface{}
		if err := w.pool.cacheService.GetJSON(cacheKey, &cachedData); err == nil {
			log.Printf("📦 Cache hit for period stats (user: %d, period: %s)", task.UserID, period)
			return cachedData, nil
		}
	}

	// Generate stats
	stats, err := w.pool.analyticsService.GetPeriodStats(task.UserID, period)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if w.pool.cacheService.IsEnabled() {
		w.pool.cacheService.SetJSON(cacheKey, stats, cache.TTLMedium)
	}

	return stats, nil
}

func (w *Worker) processMMRTrajectory(task AnalyticsTask) (interface{}, error) {
	days, ok := task.Data["days"].(int)
	if !ok {
		days = 30
	}

	// Check cache first
	cacheKey := cache.MMRCacheKey(task.UserID, days)
	if w.pool.cacheService.IsEnabled() {
		var cachedData interface{}
		if err := w.pool.cacheService.GetJSON(cacheKey, &cachedData); err == nil {
			log.Printf("📦 Cache hit for MMR trajectory (user: %d, days: %d)", task.UserID, days)
			return cachedData, nil
		}
	}

	// Calculate trajectory
	trajectory, err := w.pool.analyticsService.GetMMRTrajectory(task.UserID, days)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if w.pool.cacheService.IsEnabled() {
		w.pool.cacheService.SetJSON(cacheKey, trajectory, cache.TTLLong)
	}

	return trajectory, nil
}

func (w *Worker) processRecommendations(task AnalyticsTask) (interface{}, error) {
	// Check cache first
	cacheKey := cache.RecommendationCacheKey(task.UserID)
	if w.pool.cacheService.IsEnabled() {
		var cachedData interface{}
		if err := w.pool.cacheService.GetJSON(cacheKey, &cachedData); err == nil {
			log.Printf("📦 Cache hit for recommendations (user: %d)", task.UserID)
			return cachedData, nil
		}
	}

	// Generate recommendations
	recommendations, err := w.pool.analyticsService.GetRecommendations(task.UserID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if w.pool.cacheService.IsEnabled() {
		w.pool.cacheService.SetJSON(cacheKey, recommendations, cache.TTLLong)
	}

	return recommendations, nil
}

func (w *Worker) processChampionAnalysis(task AnalyticsTask) (interface{}, error) {
	championID, ok := task.Data["champion_id"].(int)
	if !ok {
		return nil, fmt.Errorf("champion_id required for champion analysis")
	}

	period, ok := task.Data["period"].(string)
	if !ok {
		period = "week"
	}

	// Check cache first
	cacheKey := cache.ChampionCacheKey(task.UserID, championID, period)
	if w.pool.cacheService.IsEnabled() {
		var cachedData interface{}
		if err := w.pool.cacheService.GetJSON(cacheKey, &cachedData); err == nil {
			log.Printf("📦 Cache hit for champion analysis (user: %d, champion: %d)", task.UserID, championID)
			return cachedData, nil
		}
	}

	// Analyze champion
	analysis, err := w.pool.analyticsService.GetChampionAnalysis(task.UserID, championID, period)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if w.pool.cacheService.IsEnabled() {
		w.pool.cacheService.SetJSON(cacheKey, analysis, cache.TTLMedium)
	}

	return analysis, nil
}

func (w *Worker) processCacheWarmup(task AnalyticsTask) error {
	// Warm up cache for user by pre-calculating common analytics
	log.Printf("🔥 Warming up cache for user %d", task.UserID)

	// Submit multiple tasks for cache warmup
	tasks := []AnalyticsTask{
		{Type: TaskPeriodStats, UserID: task.UserID, Data: map[string]interface{}{"period": "week"}},
		{Type: TaskPeriodStats, UserID: task.UserID, Data: map[string]interface{}{"period": "month"}},
		{Type: TaskMMRTrajectory, UserID: task.UserID, Data: map[string]interface{}{"days": 30}},
		{Type: TaskRecommendations, UserID: task.UserID, Data: map[string]interface{}{}},
	}

	for _, warmupTask := range tasks {
		warmupTask.Priority = 3 // Low priority for warmup tasks
		w.pool.SubmitTask(warmupTask)
	}

	return nil
}

func (w *Worker) processCacheInvalidate(task AnalyticsTask) error {
	// Invalidate cache for user
	log.Printf("🗑️  Invalidating cache for user %d", task.UserID)

	if !w.pool.cacheService.IsEnabled() {
		return nil
	}

	// Delete user-specific cache patterns
	patterns := []string{
		fmt.Sprintf("user:%d:*", task.UserID),
		fmt.Sprintf("analytics:%d:*", task.UserID),
		fmt.Sprintf("mmr:%d:*", task.UserID),
		fmt.Sprintf("recommendations:%d", task.UserID),
		fmt.Sprintf("champion:%d:*", task.UserID),
	}

	for _, pattern := range patterns {
		w.pool.cacheService.DeletePattern(pattern)
	}

	return nil
}

// Worker utility methods

func (w *Worker) setActive(active bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.active = active
}

func (w *Worker) isActive() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.active
}

// Pool utility methods

func (wp *WorkerPool) processResults() {
	defer wp.wg.Done()

	for {
		select {
		case result := <-wp.resultQueue:
			wp.handleTaskResult(result)
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) handleTaskResult(result TaskResult) {
	wp.stats.incrementTasksProcessed()
	wp.stats.addProcessTime(result.Duration)

	if result.Success {
		wp.stats.incrementTasksSucceeded()
	} else {
		wp.stats.incrementTasksFailed()
		log.Printf("❌ Task %s failed: %v", result.TaskID, result.Error)
	}
}

func (wp *WorkerPool) updateStats() {
	defer wp.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wp.stats.updateAverageTaskTime()
		case <-wp.ctx.Done():
			return
		}
	}
}

// Stats methods

func (s *WorkerPoolStats) incrementTasksProcessed() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TasksProcessed++
}

func (s *WorkerPoolStats) incrementTasksSucceeded() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TasksSucceeded++
}

func (s *WorkerPoolStats) incrementTasksFailed() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TasksFailed++
}

func (s *WorkerPoolStats) incrementTasksInQueue() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TasksInQueue++
}

func (s *WorkerPoolStats) addProcessTime(duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.TotalProcessTime += duration
}

func (s *WorkerPoolStats) updateAverageTaskTime() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.TasksProcessed > 0 {
		s.AverageTaskTime = s.TotalProcessTime / time.Duration(s.TasksProcessed)
	}
}