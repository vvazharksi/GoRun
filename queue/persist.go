package queue

import (
	"encoding/json"
	"os"

	"github.com/vvazharksi/GoRun/config"
)

func LoadQueue() ([]QueueJob, error) {
	if _, err := os.Stat(config.QueueFile); os.IsNotExist(err) {
		return []QueueJob{}, nil
	}

	data, err := os.ReadFile(config.QueueFile)
	if err != nil {
		return nil, err
	}

	var jobs []QueueJob
	err = json.Unmarshal(data, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func saveQueue(jobs []QueueJob) error {
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.QueueFile, data, 0644)
}
