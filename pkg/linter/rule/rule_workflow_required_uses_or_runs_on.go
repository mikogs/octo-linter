package rule

import (
	"fmt"

	"gopkg.pl/mikogs/octo-linter/pkg/dotgithub"
	"gopkg.pl/mikogs/octo-linter/pkg/workflow"
)

type RuleWorkflowRequiredUsesOrRunsOn struct {
	Value      bool
	ConfigName string
	LogLevel   int
	IsError    bool
}

func (r RuleWorkflowRequiredUsesOrRunsOn) Validate() error {
	return nil
}

func (r RuleWorkflowRequiredUsesOrRunsOn) Lint(f dotgithub.File, d *dotgithub.DotGithub, chWarnings chan<- string, chErrors chan<- string) (compliant bool, err error) {
	compliant = true
	if f.GetType() != DotGithubFileTypeWorkflow {
		return
	}
	w := f.(*workflow.Workflow)

	if !r.Value || w.Jobs == nil || len(w.Jobs) == 0 {
		return
	}

	for jobName, job := range w.Jobs {
		if job.RunsOn == nil && job.Uses == "" {
			compliant = false
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' job '%s' should have either 'uses' or 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
		}

		runsOnStr, ok := job.RunsOn.(string)
		if ok {
			if job.Uses == "" && runsOnStr == "" {
				compliant = false
				printErrOrWarn(r.ConfigName, r.IsError, r.LogLevel, fmt.Sprintf("workflow '%s' job '%s' should have either 'uses' or 'runs-on' field", w.FileName, jobName), chWarnings, chErrors)
			}
		}
	}

	return
}

func (r RuleWorkflowRequiredUsesOrRunsOn) GetConfigName() string {
	return r.ConfigName
}
