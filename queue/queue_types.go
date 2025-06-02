// package queue

// import "time"

// type QueueJob struct {
// 	JobID       string    `json:"job_id"`
// 	PipelineID  string    `json:"pipeline_id"`
// 	ProjectName string    `json:"project_name"`
// 	Stage       string    `json:"stage"`
// 	Name        string    `json:"name"`
// 	Status      string    `json:"status"` // waiting, running, done, failed
// 	EnqueuedAt  time.Time `json:"enqueued_at"`
// 	TriggeredBy string    `json:"triggered_by"`
// }

package queue

import "time"

type QueueJob struct {
	JobID       string    `json:"job_id"`
	PipelineID  string    `json:"pipeline_id"`
	ProjectName string    `json:"project_name"`
	Stage       string    `json:"stage"`
	Name        string    `json:"name"`
	Status      string    `json:"status"` // waiting, running, done, failed
	EnqueuedAt  time.Time `json:"enqueued_at"`
	TriggeredBy string    `json:"triggered_by"`

	Image  string                 `json:"image"`
	Vars   map[string]interface{} `json:"vars"`
	Action map[string]interface{} `json:"action"`
}
