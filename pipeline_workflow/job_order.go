package pipeline_workflow

import (
	"errors"
	"fmt"

	"github.com/vvazharksi/GoRun/utils"
)

type JobNode struct {
	Job      *Job
	Depends  []*JobNode
	Resolved bool
}

// BuildExecutionPlan organizes jobs by stages and dependencies
func BuildExecutionPlan(pipeline *Pipeline) ([][]*Job, error) {
	pipeline.ExecutionID = utils.GenerateTimestampID()

	stageMap := map[string][]*Job{}
	jobLookup := map[string]*Job{}

	// Group jobs by stage and build job name lookup
	for i := range pipeline.Jobs {
		job := &pipeline.Jobs[i]

		job.PipelineExecutionID = pipeline.ExecutionID
		job.JobInstanceID = utils.GenerateTimestampID()

		stageMap[job.Stage] = append(stageMap[job.Stage], job)
		jobLookup[job.Name] = job
	}

	var executionPlan [][]*Job

	for _, stageName := range pipeline.Stages {
		jobsInStage := stageMap[stageName]
		if jobsInStage == nil {
			continue
		}

		// Build job nodes for topological sort
		nodeMap := map[string]*JobNode{}
		for _, job := range jobsInStage {
			nodeMap[job.Name] = &JobNode{Job: job}
		}

		// Set dependencies
		for _, job := range jobsInStage {
			node := nodeMap[job.Name]
			for _, depName := range job.Dependencies {
				if depJob, ok := jobLookup[depName]; ok {
					// Allow dependencies from earlier stages
					node.Depends = append(node.Depends, &JobNode{Job: depJob})
				} else {
					return nil, fmt.Errorf("job '%s' depends on unknown job '%s'", job.Name, depName)
				}
			}
		}

		// Topological sort
		orderedJobs, err := topologicalSortStage(nodeMap)
		if err != nil {
			return nil, err
		}

		for _, level := range orderedJobs {
			executionPlan = append(executionPlan, level)
		}
	}

	return executionPlan, nil
}

// Sorts jobs topologically into levels (batches of jobs that can run in parallel)
func topologicalSortStage(nodeMap map[string]*JobNode) ([][]*Job, error) {
	inDegree := map[string]int{}
	graph := map[string][]string{}

	// Initialize
	for name := range nodeMap {
		inDegree[name] = 0
		graph[name] = []string{}
	}

	// Build graph and in-degree count
	for name, node := range nodeMap {
		for _, dep := range node.Depends {
			if _, ok := nodeMap[dep.Job.Name]; ok {
				graph[dep.Job.Name] = append(graph[dep.Job.Name], name)
				inDegree[name]++
			}
		}
	}

	var result [][]*Job

	// Kahn's algorithm level-by-level
	for {
		var level []*Job

		for name, deg := range inDegree {
			if deg == 0 {
				level = append(level, nodeMap[name].Job)
			}
		}

		if len(level) == 0 {
			break
		}

		// Remove from graph
		for _, job := range level {
			delete(inDegree, job.Name)
			for _, dependent := range graph[job.Name] {
				inDegree[dependent]--
			}
		}

		result = append(result, level)
	}

	if len(inDegree) > 0 {
		return nil, errors.New("circular dependency detected in jobs")
	}

	return result, nil
}
