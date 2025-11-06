// ===== ตัวอย่าง: Background Job Processing System =====
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ===== Part 1: Job Definition =====

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusRetrying  JobStatus = "retrying"
)

type Job struct {
	ID          string
	Type        string
	Payload     map[string]interface{}
	Status      JobStatus
	Priority    int
	MaxRetries  int
	RetryCount  int
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Error       string
	Result      interface{}
}

// ===== Part 2: Job Handler =====

type JobHandler func(ctx context.Context, job *Job) error

// ===== Part 3: Job Queue =====

type JobQueue struct {
	jobs     []*Job
	mu       sync.RWMutex
	notEmpty chan struct{}
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs:     make([]*Job, 0),
		notEmpty: make(chan struct{}, 1),
	}
}

func (jq *JobQueue) Enqueue(job *Job) {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	// Insert job based on priority (higher priority first)
	inserted := false
	for i, existingJob := range jq.jobs {
		if job.Priority > existingJob.Priority {
			jq.jobs = append(jq.jobs[:i], append([]*Job{job}, jq.jobs[i:]...)...)
			inserted = true
			break
		}
	}

	if !inserted {
		jq.jobs = append(jq.jobs, job)
	}

	// Notify waiting workers
	select {
	case jq.notEmpty <- struct{}{}:
	default:
	}
}

func (jq *JobQueue) Dequeue() *Job {
	jq.mu.Lock()
	defer jq.mu.Unlock()

	if len(jq.jobs) == 0 {
		return nil
	}

	job := jq.jobs[0]
	jq.jobs = jq.jobs[1:]

	return job
}

func (jq *JobQueue) Size() int {
	jq.mu.RLock()
	defer jq.mu.RUnlock()

	return len(jq.jobs)
}

func (jq *JobQueue) WaitForJob(ctx context.Context) (*Job, error) {
	for {
		job := jq.Dequeue()
		if job != nil {
			return job, nil
		}

		select {
		case <-jq.notEmpty:
			// Queue has jobs, try dequeuing again
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// ===== Part 4: Job Worker =====

type Worker struct {
	ID       int
	queue    *JobQueue
	handlers map[string]JobHandler
	stats    *WorkerStats
}

type WorkerStats struct {
	JobsProcessed int
	JobsFailed    int
	mu            sync.RWMutex
}

func NewWorker(id int, queue *JobQueue) *Worker {
	return &Worker{
		ID:       id,
		queue:    queue,
		handlers: make(map[string]JobHandler),
		stats:    &WorkerStats{},
	}
}

func (w *Worker) RegisterHandler(jobType string, handler JobHandler) {
	w.handlers[jobType] = handler
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("🚀 Worker %d started\n", w.ID)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("⏹️  Worker %d stopped\n", w.ID)
			return
		default:
			job, err := w.queue.WaitForJob(ctx)
			if err != nil {
				return
			}

			if job != nil {
				w.processJob(ctx, job)
			}
		}
	}
}

func (w *Worker) processJob(ctx context.Context, job *Job) {
	handler, exists := w.handlers[job.Type]
	if !exists {
		job.Status = JobStatusFailed
		job.Error = fmt.Sprintf("no handler for job type: %s", job.Type)
		fmt.Printf("❌ Worker %d: No handler for job %s (type: %s)\n", w.ID, job.ID, job.Type)
		return
	}

	fmt.Printf("⚙️  Worker %d: Processing job %s (type: %s, priority: %d)\n",
		w.ID, job.ID, job.Type, job.Priority)

	job.Status = JobStatusRunning
	now := time.Now()
	job.StartedAt = &now

	err := handler(ctx, job)

	if err != nil {
		job.RetryCount++

		if job.RetryCount < job.MaxRetries {
			job.Status = JobStatusRetrying
			fmt.Printf("⚠️  Worker %d: Job %s failed (retry %d/%d): %v\n",
				w.ID, job.ID, job.RetryCount, job.MaxRetries, err)

			// Requeue with exponential backoff
			go func() {
				backoff := time.Duration(job.RetryCount) * time.Second
				time.Sleep(backoff)
				w.queue.Enqueue(job)
			}()

			w.stats.mu.Lock()
			w.stats.JobsFailed++
			w.stats.mu.Unlock()
		} else {
			job.Status = JobStatusFailed
			job.Error = err.Error()
			completed := time.Now()
			job.CompletedAt = &completed

			fmt.Printf("❌ Worker %d: Job %s failed permanently: %v\n", w.ID, job.ID, err)

			w.stats.mu.Lock()
			w.stats.JobsFailed++
			w.stats.mu.Unlock()
		}
	} else {
		job.Status = JobStatusCompleted
		completed := time.Now()
		job.CompletedAt = &completed

		duration := completed.Sub(*job.StartedAt)
		fmt.Printf("✅ Worker %d: Job %s completed in %v\n", w.ID, job.ID, duration)

		w.stats.mu.Lock()
		w.stats.JobsProcessed++
		w.stats.mu.Unlock()
	}
}

func (w *Worker) GetStats() (int, int) {
	w.stats.mu.RLock()
	defer w.stats.mu.RUnlock()

	return w.stats.JobsProcessed, w.stats.JobsFailed
}

// ===== Part 5: Job Scheduler =====

type ScheduledJob struct {
	Job      *Job
	Schedule time.Duration
	NextRun  time.Time
}

type Scheduler struct {
	queue *JobQueue
	jobs  []*ScheduledJob
	mu    sync.RWMutex
}

func NewScheduler(queue *JobQueue) *Scheduler {
	return &Scheduler{
		queue: queue,
		jobs:  make([]*ScheduledJob, 0),
	}
}

func (s *Scheduler) Schedule(job *Job, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scheduledJob := &ScheduledJob{
		Job:      job,
		Schedule: interval,
		NextRun:  time.Now().Add(interval),
	}

	s.jobs = append(s.jobs, scheduledJob)
	fmt.Printf("📅 Scheduled job %s to run every %v\n", job.ID, interval)
}

func (s *Scheduler) Start(ctx context.Context) {
	fmt.Println("📅 Scheduler started")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("📅 Scheduler stopped")
			return
		case <-ticker.C:
			s.checkAndRunJobs()
		}
	}
}

func (s *Scheduler) checkAndRunJobs() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	for _, scheduledJob := range s.jobs {
		if now.After(scheduledJob.NextRun) {
			// Create a new job instance (copy)
			newJob := &Job{
				ID:         fmt.Sprintf("%s-%d", scheduledJob.Job.ID, time.Now().Unix()),
				Type:       scheduledJob.Job.Type,
				Payload:    scheduledJob.Job.Payload,
				Status:     JobStatusPending,
				Priority:   scheduledJob.Job.Priority,
				MaxRetries: scheduledJob.Job.MaxRetries,
				CreatedAt:  time.Now(),
			}

			s.queue.Enqueue(newJob)

			scheduledJob.NextRun = now.Add(scheduledJob.Schedule)
			fmt.Printf("📅 Enqueued scheduled job: %s\n", newJob.ID)
		}
	}
}

// ===== Part 6: Job Manager =====

type JobManager struct {
	queue     *JobQueue
	workers   []*Worker
	scheduler *Scheduler
	jobs      map[string]*Job
	mu        sync.RWMutex
}

func NewJobManager(numWorkers int) *JobManager {
	queue := NewJobQueue()
	scheduler := NewScheduler(queue)

	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i+1, queue)
	}

	return &JobManager{
		queue:     queue,
		workers:   workers,
		scheduler: scheduler,
		jobs:      make(map[string]*Job),
	}
}

func (jm *JobManager) RegisterHandler(jobType string, handler JobHandler) {
	for _, worker := range jm.workers {
		worker.RegisterHandler(jobType, handler)
	}
}

func (jm *JobManager) Start(ctx context.Context) {
	// Start scheduler
	go jm.scheduler.Start(ctx)

	// Start workers
	for _, worker := range jm.workers {
		go worker.Start(ctx)
	}
}

func (jm *JobManager) Submit(job *Job) {
	jm.mu.Lock()
	jm.jobs[job.ID] = job
	jm.mu.Unlock()

	jm.queue.Enqueue(job)
	fmt.Printf("📥 Job %s submitted (type: %s, priority: %d)\n", job.ID, job.Type, job.Priority)
}

func (jm *JobManager) Schedule(job *Job, interval time.Duration) {
	jm.scheduler.Schedule(job, interval)
}

func (jm *JobManager) GetJobStatus(jobID string) (*Job, bool) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	job, exists := jm.jobs[jobID]
	return job, exists
}

func (jm *JobManager) GetStats() map[int]struct{ Processed, Failed int } {
	stats := make(map[int]struct{ Processed, Failed int })

	for _, worker := range jm.workers {
		processed, failed := worker.GetStats()
		stats[worker.ID] = struct{ Processed, Failed int }{processed, failed}
	}

	return stats
}

// ===== Example Job Handlers =====

func SendEmailHandler(ctx context.Context, job *Job) error {
	to := job.Payload["to"].(string)
	subject := job.Payload["subject"].(string)

	// Simulate email sending
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	// Simulate occasional failures
	if rand.Float64() < 0.2 {
		return errors.New("email service unavailable")
	}

	job.Result = fmt.Sprintf("Email sent to %s: %s", to, subject)
	return nil
}

func ProcessImageHandler(ctx context.Context, job *Job) error {
	imageURL := job.Payload["image_url"].(string)

	// Simulate image processing
	time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)

	job.Result = fmt.Sprintf("Processed image: %s", imageURL)
	return nil
}

func GenerateReportHandler(ctx context.Context, job *Job) error {
	reportType := job.Payload["type"].(string)

	// Simulate report generation
	time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)

	job.Result = fmt.Sprintf("Generated %s report", reportType)
	return nil
}

func CleanupHandler(ctx context.Context, job *Job) error {
	// Simulate cleanup task
	time.Sleep(100 * time.Millisecond)

	job.Result = "Cleanup completed"
	return nil
}

// ===== Main Demo =====

func main() {
	fmt.Println("===== Background Job Processing System =====\n")
	rand.Seed(time.Now().UnixNano())

	// Create job manager with 3 workers
	manager := NewJobManager(3)

	// Register job handlers
	manager.RegisterHandler("send_email", SendEmailHandler)
	manager.RegisterHandler("process_image", ProcessImageHandler)
	manager.RegisterHandler("generate_report", GenerateReportHandler)
	manager.RegisterHandler("cleanup", CleanupHandler)

	// Start the system
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	manager.Start(ctx)

	time.Sleep(500 * time.Millisecond) // Wait for workers to start

	// ===== Demo 1: Submit Jobs with Different Priorities =====
	fmt.Println("\n--- Demo 1: Job Submission (Different Priorities) ---\n")

	jobs := []*Job{
		{
			ID:       "job-1",
			Type:     "send_email",
			Payload:  map[string]interface{}{"to": "user1@example.com", "subject": "Welcome!"},
			Status:   JobStatusPending,
			Priority: 5,
			MaxRetries: 3,
			CreatedAt: time.Now(),
		},
		{
			ID:       "job-2",
			Type:     "process_image",
			Payload:  map[string]interface{}{"image_url": "https://example.com/image1.jpg"},
			Status:   JobStatusPending,
			Priority: 10, // High priority
			MaxRetries: 2,
			CreatedAt: time.Now(),
		},
		{
			ID:       "job-3",
			Type:     "generate_report",
			Payload:  map[string]interface{}{"type": "monthly"},
			Status:   JobStatusPending,
			Priority: 3,
			MaxRetries: 1,
			CreatedAt: time.Now(),
		},
		{
			ID:       "job-4",
			Type:     "send_email",
			Payload:  map[string]interface{}{"to": "user2@example.com", "subject": "Update"},
			Status:   JobStatusPending,
			Priority: 7,
			MaxRetries: 3,
			CreatedAt: time.Now(),
		},
	}

	for _, job := range jobs {
		manager.Submit(job)
		time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(3 * time.Second)

	// ===== Demo 2: Scheduled Jobs =====
	fmt.Println("\n--- Demo 2: Scheduled Jobs ---\n")

	cleanupJob := &Job{
		ID:       "cleanup",
		Type:     "cleanup",
		Payload:  map[string]interface{}{"task": "old_files"},
		Status:   JobStatusPending,
		Priority: 5,
		MaxRetries: 1,
		CreatedAt: time.Now(),
	}

	manager.Schedule(cleanupJob, 3*time.Second)

	time.Sleep(8 * time.Second)

	// ===== Demo 3: Check Job Status =====
	fmt.Println("\n--- Demo 3: Job Status ---\n")

	for _, job := range jobs {
		if currentJob, exists := manager.GetJobStatus(job.ID); exists {
			fmt.Printf("Job %s: Status=%s", currentJob.ID, currentJob.Status)
			if currentJob.Error != "" {
				fmt.Printf(", Error=%s", currentJob.Error)
			}
			if currentJob.Result != nil {
				fmt.Printf(", Result=%v", currentJob.Result)
			}
			fmt.Println()
		}
	}

	// ===== Demo 4: Worker Statistics =====
	fmt.Println("\n--- Demo 4: Worker Statistics ---\n")

	stats := manager.GetStats()
	for workerID, stat := range stats {
		fmt.Printf("Worker %d: Processed=%d, Failed=%d\n",
			workerID, stat.Processed, stat.Failed)
	}

	// Wait for remaining jobs
	time.Sleep(2 * time.Second)

	// ===== Summary =====
	fmt.Println("\n" + repeatStr("=", 60))
	fmt.Println("Background Job Processing Summary:")
	fmt.Println("  ✓ Job Queue: Priority-based job queuing")
	fmt.Println("  ✓ Workers: Concurrent job processing")
	fmt.Println("  ✓ Retry: Automatic retry with exponential backoff")
	fmt.Println("  ✓ Scheduler: Periodic job execution")
	fmt.Println("  ✓ Monitoring: Job status tracking")
	fmt.Println("  ✓ Scalability: Easy to add more workers")
	fmt.Println(repeatStr("=", 60))
}

func repeatStr(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// วิธีรัน:
// go run 09_job_queue.go
//
// Background Job Processing Concepts:
//
// 1. Job Queue:
//    - Priority-based ordering
//    - Thread-safe operations
//    - Non-blocking enqueue
//    - Blocking dequeue with timeout
//
// 2. Worker Pool:
//    - Fixed number of workers
//    - Concurrent job processing
//    - Graceful shutdown
//    - Statistics tracking
//
// 3. Job Retry:
//    - Automatic retry on failure
//    - Exponential backoff
//    - Max retry limit
//    - Status tracking
//
// 4. Job Scheduler:
//    - Periodic job execution
//    - Cron-like scheduling
//    - Job cloning for each run
//
// 5. Job Status:
//    - Pending → Running → Completed/Failed
//    - Retrying for temporary failures
//    - Timestamps for tracking
//
// Benefits:
// - Async processing
// - Better user experience
// - Resource optimization
// - Failure handling
// - Load balancing
// - Scalability
//
// Use Cases:
// - Email sending
// - Image processing
// - Report generation
// - Data import/export
// - Cleanup tasks
// - Notifications
// - Webhooks
// - PDF generation
//
// Best Practices:
// - Set appropriate priorities
// - Implement idempotent handlers
// - Use timeouts
// - Monitor queue depth
// - Implement dead letter queue
// - Log job execution
// - Track metrics
// - Use database for persistence
//
// Production Considerations:
// - Persistent queue (Redis, PostgreSQL)
// - Distributed workers
// - Job deduplication
// - Rate limiting
// - Circuit breaker integration
// - Monitoring and alerting
// - Job cancellation
// - Graceful shutdown
// - Resource limits
// - Priority inversion prevention
//
// Popular Go Libraries:
// - Asynq (Redis-backed)
// - Machinery
// - gocraft/work
// - RabbitMQ with workers
