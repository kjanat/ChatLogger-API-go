package jobs

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task types
const (
	TypeExportProcess = "export:process"
)

// ExportPayload contains the data needed for export processing
type ExportPayload struct {
	ExportID uint64 `json:"export_id"`
}

// Queue handles job queueing operations
type Queue struct {
	client *asynq.Client
}

// NewQueue creates a new job queue
func NewQueue(redisAddr string) *Queue {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &Queue{client: client}
}

// EnqueueExport adds an export job to the queue
func (q *Queue) EnqueueExport(exportID uint64) error {
	payload, err := json.Marshal(ExportPayload{ExportID: exportID})
	if err != nil {
		return fmt.Errorf("failed to marshal export payload: %w", err)
	}

	task := asynq.NewTask(TypeExportProcess, payload)

	// Set some processing options
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue("exports"),
		asynq.Timeout(20 * 60), // 20 minutes timeout for long exports
	}

	_, err = q.client.Enqueue(task, opts...)
	return err
}

// Close closes the queue client connection
func (q *Queue) Close() error {
	return q.client.Close()
}
