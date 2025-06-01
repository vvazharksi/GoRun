package pipeline_workflow

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ReadPipelineYaml reads and parses the YAML file into a Pipeline struct
func ReadPipelineYaml(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pipeline Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

// getScriptSteps extracts script steps from the action block
func GetScriptSteps(action map[string]interface{}) []string {
	scriptRaw, ok := action["script"]
	if !ok {
		return nil
	}

	if rawList, ok := scriptRaw.([]interface{}); ok {
		var steps []string
		for _, entry := range rawList {
			if str, ok := entry.(string); ok {
				steps = append(steps, str)
			}
		}
		return steps
	}

	return nil
}
