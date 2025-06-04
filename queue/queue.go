package queue

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// var queueFile = "queue/jobs.json"
var historyFile = "queue/jobs_history.json"
var queueLock sync.Mutex

func Enqueue(job QueueJob) error {
	queueLock.Lock()
	defer queueLock.Unlock()

	jobs, _ := LoadQueue()

	jobs = append(jobs, job)
	return saveQueue(jobs)
}

func Dequeue() (*QueueJob, error) {
	queueLock.Lock()
	defer queueLock.Unlock()

	jobs, _ := LoadQueue()
	if len(jobs) == 0 {
		return nil, errors.New("queue is empty")
	}

	job := jobs[0]
	remaining := jobs[1:]

	err := saveQueue(remaining)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func ListWaitingJobs() ([]QueueJob, error) {
	jobs, err := LoadQueue()
	if err != nil {
		return nil, err
	}

	var waiting []QueueJob
	for _, job := range jobs {
		if job.Status == "waiting" {
			waiting = append(waiting, job)
		}
	}
	return waiting, nil
}

// Usage:
// err := queue.RemoveJobByID("2506012259000-001")
//
//	if err != nil {
//		log.Printf("Failed to remove job from queue: %v\n", err)
//	}
func RemoveJobByID(jobID string) error {
	queueLock.Lock()
	defer queueLock.Unlock()

	jobs, err := LoadQueue()
	if err != nil {
		return err
	}

	var updated []QueueJob
	found := false

	for _, job := range jobs {
		if job.JobID == jobID {
			found = true
			continue // Skip this job (i.e., remove it)
		}
		updated = append(updated, job)
	}

	if !found {
		return errors.New("job not found in queue")
	}

	return saveQueue(updated)
}

func UpdateJobStatus(jobID string, status string) error {
	// 	queueLock.Lock()
	// 	defer queueLock.Unlock()

	// 	jobs, err := LoadQueue()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	updated := false
	// 	for i := range jobs {
	// 		if jobs[i].JobID == jobID {
	// 			jobs[i].Status = status
	// 			updated = true
	// 			break
	// 		}
	// 	}

	// 	if !updated {
	// 		return errors.New("job not found in queue")
	// 	}

	// 	return saveQueue(jobs)
	// }
	queueLock.Lock()
	defer queueLock.Unlock()

	jobs, err := LoadQueue()
	if err != nil {
		return err
	}

	var updatedJob *QueueJob
	updated := false
	for i := range jobs {
		if jobs[i].JobID == jobID {
			jobs[i].Status = status
			updatedJob = &jobs[i]
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("job not found in queue")
	}

	// Save updated queue
	err = saveQueue(jobs)
	if err != nil {
		return err
	}

	// Save job to history
	return appendJobHistory(*updatedJob)
}

func appendJobHistory(job QueueJob) error {
	var history []QueueJob

	data, err := os.ReadFile(historyFile)
	if err == nil {
		_ = json.Unmarshal(data, &history)
	}

	history = append(history, job)

	out, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, out, 0644)
}

func loadJobHistory() ([]QueueJob, error) {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []QueueJob{}, nil
		}
		return nil, err
	}

	var history []QueueJob
	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func GetJobByID(jobID string) (*QueueJob, error) {
	jobs, err := LoadQueue()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.JobID == jobID {
			return &job, nil
		}
	}

	return nil, errors.New("job not found in queue")
}
