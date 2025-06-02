package pipeline_workflow

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type DockerRunner struct {
	Logger *log.Logger
}

func NewDockerRunner(logger *log.Logger) *DockerRunner {
	return &DockerRunner{Logger: logger}
}

func (dr *DockerRunner) Run(job Job) error {
	containerName := fmt.Sprintf("job-%s", job.JobInstanceID)

	var envVars []string
	for k, v := range job.Vars {
		envVars = append(envVars, "-e", fmt.Sprintf("%s=%v", k, v))
	}

	script := GetScriptSteps(job.Action)
	if len(script) == 0 {
		return fmt.Errorf("no script steps to run")
	}
	joinedScript := strings.Join(script, " && ")

	args := append([]string{
		"run", "--rm", "--name", containerName,
	}, envVars...)
	args = append(args, job.Image, "bash", "-c", joinedScript)

	// cmd := exec.Command("docker", args...)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	// defer cancel()
	// cmd = cmd.WithContext(ctx)

	dr.Logger.Printf("Starting job %s in Docker container: %s", job.Name, containerName)

	err := cmd.Run()

	dr.Logger.Printf("Output for job %s:\n%s", job.Name, outBuf.String())
	if err != nil {
		dr.Logger.Printf("Error output for job %s:\n%s", job.Name, errBuf.String())
		return fmt.Errorf("job %s failed: %v", job.Name, err)
	}

	return nil
}
