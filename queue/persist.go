package queue

import (
	"encoding/json"
	"os"
)

func LoadQueue() ([]QueueJob, error) {
	if _, err := os.Stat(queueFile); os.IsNotExist(err) {
		return []QueueJob{}, nil
	}

	data, err := os.ReadFile(queueFile)
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
	return os.WriteFile(queueFile, data, 0644)
}
