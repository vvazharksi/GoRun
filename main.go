package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vvazharksi/GoRun/config"
	pipeline_workflow "github.com/vvazharksi/GoRun/pipeline_workflow"
	"github.com/vvazharksi/GoRun/queue"
)

func main() {
	pipeline, err := pipeline_workflow.ReadPipelineYaml(config.PipelineDefinitionPath)
	if err != nil {
		log.Fatalf("Failed to read pipeline YAML: %v", err)
	}

	fmt.Println("Stages:", pipeline.Stages)

	for _, job := range pipeline.Jobs {
		fmt.Printf("\n---\nJob: %s (Stage: %s)\n", job.Name, job.Stage)
		fmt.Println("Image:", job.Image)
		fmt.Println("Vars:", job.Vars)
		fmt.Println("Dependencies:", job.Dependencies)

		if job.Action != nil {
			fmt.Println("Action:", job.Action)

			if _, ok := job.Action["script"]; ok {
				scriptSteps := pipeline_workflow.GetScriptSteps(job.Action)
				for i, cmd := range scriptSteps {
					fmt.Printf("Script Step %d:\n%s\n", i+1, cmd)
				}
			}
		}
	}

	executionPlan, err := pipeline_workflow.BuildExecutionPlan(pipeline)
	if err != nil {
		log.Fatalf("Failed to build execution plan: %v", err)
	}

	fmt.Printf("\nPipeline Execution ID: %s\n", pipeline.ExecutionID)

	// Print execution plan level by level
	for i, level := range executionPlan {
		fmt.Printf("\n--- Execution Level %d ---\n", i+1)
		for _, job := range level {
			fmt.Printf("Job: %s (Stage: %s, JobID: %s, PipelineID: %s)\n", job.Name, job.Stage, job.JobInstanceID, job.PipelineExecutionID)
		}
	}

	for _, level := range executionPlan {
		for _, job := range level {
			qJob := queue.QueueJob{
				JobID:       job.JobInstanceID,
				PipelineID:  job.PipelineExecutionID,
				ProjectName: pipeline.ProjectName,
				Name:        job.Name,
				Stage:       job.Stage,
				Status:      "waiting",
				EnqueuedAt:  time.Now(),
				TriggeredBy: job.TriggeredBy,
				Image:       job.Image,
				Vars:        job.Vars,
				Action:      job.Action,
			}
			_ = queue.Enqueue(qJob)
		}
	}

	pipeline_workflow.StartJobWorker()

}

// Roadmap:
//
// save job execution logs stdout
// create artifacts and share them with jobs defined under `dependencies`
// create concept of global variables
//
// create logic of project to have the pipeline under, create concept of project variables ,keep logs under project, project to have ID, use the project ID for the runner path alongside job id
// allow project to be fetched from github for the pipeline code
// use DB instead of files
// create cli UI
//
// introduce kubelet or minikube for runners orchastration
// create runner_config and setup the runner
// connect to ec2 for job execution
//
// add basic concept of users
// add permissions
//
// tf plugin as ci definition and run
