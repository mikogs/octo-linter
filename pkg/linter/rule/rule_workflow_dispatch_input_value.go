package rule

import (
	"fmt"
	"regexp"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowDispatchInputValue struct {
	Value      map[string]string
	ConfigName string
	LogLevel   int
	IsError    map[string]bool
}

func (r RuleWorkflowDispatchInputValue) Validate() error {
	if len(r.Value) > 0 {
		for k, v := range r.Value {
			if k != "name" {
				return fmt.Errorf("%s can only contain 'name' key", r.ConfigName)
			}
			if v != "lowercase-hyphens" {
				return fmt.Errorf("%s supports 'lowercase-hyphens' or empty value only", r.ConfigName)
			}
		}
	}
	return nil
}

func (r RuleWorkflowDispatchInputValue) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if len(r.Value) == 0 {
		return
	}
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if w.On == nil || w.On.WorkflowDispatch == nil || w.On.WorkflowDispatch.Inputs == nil || len(w.On.WorkflowDispatch.Inputs) == 0 {
		return
	}

	re := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
	for inputName := range w.On.WorkflowDispatch.Inputs {
		for k, v := range r.Value {
			if k == "name" && v != "" {
				m := re.MatchString(inputName)
				if !m {
					compliant = false
					printErrOrWarn(r.ConfigName, r.IsError[k], r.LogLevel, fmt.Sprintf("workflow '%s' call input '%s' %s must be lower-case and hyphens only", w.FileName, inputName, v), chWarnings, chErrors)
				}
			}
		}
	}

	return
}

func (r RuleWorkflowDispatchInputValue) GetConfigName() string {
	return r.ConfigName
}
