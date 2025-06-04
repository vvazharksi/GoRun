package pipeline_workflow

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type DockerRunner struct {
	Logger *log.Logger
}

func NewDockerRunner(logger *log.Logger) *DockerRunner {
	return &DockerRunner{Logger: logger}
}

/// Via bash -c argument
// func (dr *DockerRunner) Run(job Job) error {
// 	containerName := fmt.Sprintf("job-%s", job.JobInstanceID)

// 	var envVars []string
// 	for k, v := range job.Vars {
// 		envVars = append(envVars, "-e", fmt.Sprintf("%s=%v", k, v))
// 	}

// 	script := GetScriptSteps(job.Action)
// 	if len(script) == 0 {
// 		return fmt.Errorf("no script steps to run")
// 	}
// 	joinedScript := strings.Join(script, " && ")

// 	args := append([]string{
// 		"run", "--rm", "--name", containerName,
// 	}, envVars...)
// 	args = append(args, job.Image, "bash", "-c", joinedScript)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// 	defer cancel()

// 	cmd := exec.CommandContext(ctx, "docker", args...)
// 	var outBuf, errBuf bytes.Buffer
// 	cmd.Stdout = &outBuf
// 	cmd.Stderr = &errBuf

// 	dr.Logger.Printf("Starting job %s in Docker container: %s", job.Name, containerName)

// 	err := cmd.Run()

// 	dr.Logger.Printf("Output for job %s:\n%s", job.Name, outBuf.String())
// 	if err != nil {
// 		dr.Logger.Printf("Error output for job %s:\n%s", job.Name, errBuf.String())
// 		return fmt.Errorf("job %s failed: %v", job.Name, err)
// 	}

// 	return nil
// }

// / Via file-mounted scripts
func (dr *DockerRunner) Run(job Job) error {
	containerName := fmt.Sprintf("job-%s", job.JobInstanceID)

	// Create working directory
	baseWorkDir := "/tmp/ci/run"
	jobDir := filepath.Join(baseWorkDir, job.JobInstanceID)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return fmt.Errorf("failed to create job dir: %w", err)
	}

	// Prepare script
	scriptSteps := GetScriptSteps(job.Action)
	if len(scriptSteps) == 0 {
		return fmt.Errorf("no script steps to run")
	}
	scriptContent := strings.Join(scriptSteps, "\n")
	scriptPath := filepath.Join(jobDir, "run.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to write run.sh: %w", err)
	}

	// Prepare .env file
	var envBuilder strings.Builder
	for k, v := range job.Vars {
		envBuilder.WriteString(fmt.Sprintf("%s=%v\n", k, v))
	}
	envPath := filepath.Join(jobDir, ".env")
	if err := os.WriteFile(envPath, []byte(envBuilder.String()), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	// Docker command
	args := []string{
		"run", "--rm", "--name", containerName,
		"--env-file", envPath,
		"-v", fmt.Sprintf("%s:/ci", jobDir),
		job.Image,
		"bash", "/ci/run.sh",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	dr.Logger.Printf("Starting job %s in Docker container: %s", job.Name, containerName)

	err := cmd.Run()

	dr.Logger.Printf("Output for job %s:\n%s", job.Name, outBuf.String())
	if err != nil {
		dr.Logger.Printf("Error output for job %s:\n%s", job.Name, errBuf.String())
		return fmt.Errorf("job %s failed: %v", job.Name, err)
	}

	return nil
}
