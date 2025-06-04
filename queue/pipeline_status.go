package queue

import (
	"encoding/json"
	"os"
	"time"
)

type JobStatus struct {
	Status string `json:"status"`
}

type PipelineStatus struct {
	PipelineID string               `json:"pipeline_id"`
	Status     string               `json:"status"`
	Jobs       map[string]JobStatus `json:"jobs"`
	StartedAt  string               `json:"started_at"`
	FinishedAt string               `json:"finished_at"`
}

//	func SavePipelineStatus(status PipelineStatus) error {
//		data, err := json.MarshalIndent(status, "", "  ")
//		if err != nil {
//			return err
//		}
//		return os.WriteFile(fmt.Sprintf("queue/pipeline_%s_status.json", status.PipelineID), data, 0644)
//	}

// func SavePipelineStatus(status PipelineStatus) error {
// 	// fileName := "queue/pipeline_" + status.PipelineID + "_status.json"
// 	fileName := "queue/pipeline_history.json"

// 	data, err := json.MarshalIndent(status, "", "  ")
// 	if err != nil {
// 		return err
// 	}

//		return os.WriteFile(fileName, data, 0644)
//	}
func SavePipelineStatus(status PipelineStatus) error {
	fileName := "queue/pipeline_history.json"

	var statuses []PipelineStatus

	// Read existing file data
	fileData, err := os.ReadFile(fileName)
	if err == nil {
		// If file exists, unmarshal existing data
		err = json.Unmarshal(fileData, &statuses)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		// Error reading the file other than not exist
		return err
	}

	// Append new status
	statuses = append(statuses, status)

	// Marshal back to JSON
	data, err := json.MarshalIndent(statuses, "", "  ")
	if err != nil {
		return err
	}

	// Write entire updated JSON array back to file
	return os.WriteFile(fileName, data, 0644)
}

// func GeneratePipelineStatus(pipelineID string, startedAt time.Time) (PipelineStatus, error) {
// 	jobs, err := LoadQueue()
// 	if err != nil {
// 		return PipelineStatus{}, err
// 	}

// 	status := "success"
// 	jobMap := make(map[string]JobStatus)

// 	for _, job := range jobs {
// 		if job.PipelineID != pipelineID {
// 			continue
// 		}
// 		jobMap[job.Name] = JobStatus{Status: job.Status}
// 		if job.Status == "failed" {
// 			status = "failed"
// 		}
// 	}

//		return PipelineStatus{
//			PipelineID: pipelineID,
//			Status:     status,
//			Jobs:       jobMap,
//			StartedAt:  startedAt.Format(time.RFC3339),
//			FinishedAt: time.Now().Format(time.RFC3339),
//		}, nil
//	}
func GeneratePipelineStatus(pipelineID string, startedAt time.Time) (PipelineStatus, error) {
	history, err := loadJobHistory()
	if err != nil {
		return PipelineStatus{}, err
	}

	status := "success"
	jobStatuses := make(map[string]JobStatus)
	var finishedAt time.Time

	for _, job := range history {
		if job.PipelineID != pipelineID {
			continue
		}
		jobStatuses[job.Name] = JobStatus{Status: job.Status}
		if job.Status == "failed" {
			status = "failed"
		}
		if job.EnqueuedAt.After(finishedAt) {
			finishedAt = job.EnqueuedAt
		}
	}

	return PipelineStatus{
		PipelineID: pipelineID,
		Status:     status,
		Jobs:       jobStatuses,
		StartedAt:  startedAt.Format(time.RFC3339),
		FinishedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func AreAllJobsCompleted(pipelineID string) (bool, string, error) {
	jobs, err := LoadQueue()
	if err != nil {
		return false, "", err
	}

	var hasFailed bool
	var hasRunningOrWaiting bool

	for _, job := range jobs {
		if job.PipelineID != pipelineID {
			continue
		}
		switch job.Status {
		case "failed":
			hasFailed = true
		case "waiting", "running":
			hasRunningOrWaiting = true
		}
	}

	if hasRunningOrWaiting {
		return false, "", nil
	}
	if hasFailed {
		return true, "failed", nil
	}
	return true, "success", nil
}
