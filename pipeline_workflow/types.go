package pipeline_workflow

import "time"

// PipelineStatus represents the state of the entire pipeline execution.
type PipelineStatus string

// JobStatus represents the current state of an individual job.
type JobStatus string

const (
	// Pipeline statuses
	PipelinePending  PipelineStatus = "pending"
	PipelineRunning  PipelineStatus = "running"
	PipelineSuccess  PipelineStatus = "success"
	PipelineFailed   PipelineStatus = "failed"
	PipelineCanceled PipelineStatus = "canceled"

	// Job statuses
	JobWaiting JobStatus = "waiting"
	JobRunning JobStatus = "running"
	JobSuccess JobStatus = "success"
	JobFailed  JobStatus = "failed"
	JobSkipped JobStatus = "skipped"
)

// Pipeline defines a full pipeline structure parsed from YAML with internal metadata.
type Pipeline struct {
	Stages      []string `yaml:"stages"`
	Jobs        []Job    `yaml:"jobs"`
	ExecutionID string
	Status      PipelineStatus
	StartedAt   time.Time
	FinishedAt  time.Time

	ProjectName string
	TriggeredBy string
}

// Job represents a single unit of execution in a pipeline.
type Job struct {
	Name         string                 `yaml:"name"`
	Stage        string                 `yaml:"stage"`
	Dependencies []string               `yaml:"dependencies,omitempty"`
	Image        string                 `yaml:"image"`
	Vars         map[string]interface{} `yaml:"vars"`
	Action       map[string]interface{} `yaml:"action"`

	PipelineExecutionID string    `yaml:"-"` // generated at runtime
	JobInstanceID       string    `yaml:"-"` // unique per execution
	Status              JobStatus `yaml:"-"`
	StartedAt           time.Time `yaml:"-"`
	FinishedAt          time.Time `yaml:"-"`
	ErrorMsg            string    `yaml:"-"`
	RetryCount          int       `yaml:"-"`

	TriggeredBy string `yaml:"-"`
}
