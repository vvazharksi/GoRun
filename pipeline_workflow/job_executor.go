package pipeline_workflow

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vvazharksi/GoRun/queue"
)

// func ExecuteJob(job Job) error {
// 	logger := log.New(os.Stdout, "[DockerRunner] ", log.LstdFlags)
// 	runner := NewDockerRunner(logger)

// 	return runner.Run(job)
// }

func ExecuteJob(qj *queue.QueueJob) error {
	logger := log.New(os.Stdout, "[DockerRunner] ", log.LstdFlags)
	runner := NewDockerRunner(logger)

	job := Job{
		JobInstanceID:       qj.JobID,
		PipelineExecutionID: qj.PipelineID,
		Name:                qj.Name,
		Stage:               qj.Stage,
		Image:               qj.Image,
		Vars:                qj.Vars,
		Action:              qj.Action,
		TriggeredBy:         qj.TriggeredBy,
	}

	return runner.Run(job)
}

func StartJobWorker() {
	fmt.Println("Job worker started...")
	for {
		waitingJobs, err := queue.ListWaitingJobs()
		if err != nil {
			fmt.Println("Failed to get jobs from queue:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(waitingJobs) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		// Take the first job
		jobMeta := waitingJobs[0]

		// Find the actual job struct (from pipeline or in-memory map)
		// Ideally, you'd keep a map[pipelineID][jobID] -> Job struct
		// For now: assume you can load it from queue or store it globally

		// This is placeholder:
		job, err := queue.GetJobByID(jobMeta.JobID)
		if err != nil {
			fmt.Println("Could not find job:", err)
			continue
		}

		fmt.Printf("Executing job: %s\n", job.Name)
		err = ExecuteJob(job)
		if err != nil {
			fmt.Printf("Job %s failed: %v\n", job.Name, err)
			_ = queue.UpdateJobStatus(jobMeta.JobID, "failed")
		} else {
			fmt.Printf("Job %s completed successfully.\n", job.Name)
			_ = queue.UpdateJobStatus(jobMeta.JobID, "done")
		}

		// Remove from queue
		_ = queue.RemoveJobByID(jobMeta.JobID)
	}
}
