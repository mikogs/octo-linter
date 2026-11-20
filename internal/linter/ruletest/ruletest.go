// Package ruletest contains helper functions for testing rules.
package ruletest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"octo-linter/internal/action"
	"octo-linter/internal/dotgithub"
	"octo-linter/internal/linter/glitch"
	"octo-linter/internal/linter/rule"
	"octo-linter/internal/workflow"
)

//nolint:gochecknoglobals
var (
	testDotGithub     *dotgithub.DotGithub
	testDotGithubOnce sync.Once
)

var errTimeout = errors.New("timeout")

// GetDotGithub returns DitGithub with test rules.
func GetDotGithub() *dotgithub.DotGithub {
	testDotGithubOnce.Do(func() {
		testDotGithub = &dotgithub.DotGithub{
			Actions:         map[string]*action.Action{},
			ExternalActions: map[string]*action.Action{},
			Workflows:       map[string]*workflow.Workflow{},
			Vars:            map[string]bool{},
			Secrets:         map[string]bool{},
		}
		logger := slog.New(slog.DiscardHandler)
		slog.SetDefault(logger)

		err := testDotGithub.ReadDir(
			context.Background(),
			"../../../../tests/rules",
			map[string]string{
				"fake-username/fake-external-action@v1": "./../../../../tests/fake-external-action",
			},
			map[string][]*regexp.Regexp{
				"fake-username/fake-external-action@v1": {
					regexp.MustCompile(`^dynamic-output-[0-9]+$`),
				},
				"aws-actions/amazon-ecr-login@v2": {
					regexp.MustCompile(
						`^docker_(username|password)_[0-9]+_dkr_ecr_[a-z1-6_]+_amazonaws_com$`,
					),
				},
			},
		)
		if err != nil {
			panic(fmt.Sprintf("failed to load test data: %v", err))
		}
	})

	return testDotGithub
}

// Lint runs a rule with a specific configuration on a specified file and returns all lint errors and a boolean
// indicating whether it is compliant or not.
func Lint(
	timeout int,
	rule rule.Rule,
	conf any,
	file dotgithub.File,
	dotGithub *dotgithub.DotGithub,
) (bool, []string, error) {
	ruleErrors := []string{}

	type lintResult struct {
		compliant bool
		err       error
	}

	chResult := make(chan lintResult, 1)
	chErrors := make(chan glitch.Glitch)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				close(chErrors)

				chResult <- lintResult{compliant: false, err: fmt.Errorf("rule panicked: %v", r)}
			}
		}()

		compliant, err := rule.Lint(conf, file, dotGithub, chErrors)
		close(chErrors)

		chResult <- lintResult{compliant: compliant, err: err}
	}()

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return false, ruleErrors, errTimeout
		case glitchInstance, more := <-chErrors:
			if more {
				ruleError := fmt.Sprintf(
					"%s %s: %s",
					glitchInstance.Path,
					glitchInstance.RuleName,
					glitchInstance.ErrText,
				)
				ruleErrors = append(ruleErrors, ruleError)
			} else {
				res := <-chResult

				return res.compliant, ruleErrors, res.err
			}
		}
	}
}

// Action runs a test function on a specific action in DotGithub.
func Action(
	dotGithub *dotgithub.DotGithub,
	actionToTest string,
	testFunc func(file dotgithub.File, name string),
) {
	actionFile, ok := dotGithub.Actions[actionToTest]
	if !ok {
		return
	}

	log.Printf("running test on action %s...", actionToTest)
	testFunc(actionFile, actionToTest)
}

// Workflow runs a test function on a specific workflow from DotGithub.
func Workflow(
	dotGithub *dotgithub.DotGithub,
	workflowToTest string,
	testFunc func(file dotgithub.File, name string),
) {
	workflowFile, ok := dotGithub.Workflows[workflowToTest]
	if !ok {
		return
	}

	log.Printf("running test on workflow %s...", workflowToTest)
	testFunc(workflowFile, workflowToTest)
}
