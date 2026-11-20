package workflow

import (
	"octo-linter/internal/step"
)

// Job represents a job in a GitHub Actions workflow parsed from YAML.
type Job struct {
	Name   string            `yaml:"name"`
	Uses   string            `yaml:"uses"`
	RunsOn any               `yaml:"runs-on"`
	Steps  []*step.Step      `yaml:"steps"`
	Env    map[string]string `yaml:"env"`
	Needs  any               `yaml:"needs,omitempty"`
}

// SetParentType sets parent type for all the steps.
func (wj *Job) SetParentType(t string) {
	for _, s := range wj.Steps {
		if s == nil {
			continue
		}

		s.ParentType = t
	}
}
