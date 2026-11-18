// Package glitch contains code related to representing a lint error.
package glitch

import (
	"fmt"
	"strings"
)

const (
	// DotGithubFileTypeAction represents the action file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeAction = 1
	// DotGithubFileTypeWorkflow represents the workflow file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeWorkflow = 2
)

// Glitch represents a linting error.
type Glitch struct {
	Type     int
	Name     string
	Path     string
	RuleName string
	ErrText  string
	IsError  bool
}

// ListToMarkdown takes a list of Glitch instances and generates a Markdown table from it.
func ListToMarkdown(glitches []*Glitch, limit int) string {
	var sb strings.Builder
	sb.WriteString(`|Item|Error|
|---|---|
`)

	for i, glitch := range glitches {
		if limit > 0 && i == limit {
			break
		}

		namePrefix := "a/"
		if glitch.Type == DotGithubFileTypeWorkflow {
			namePrefix = "w/"
		}

		name := glitch.Name
		if name == "" {
			name = "unknown"
		}
		name = namePrefix + name

		level := `🟠`
		if glitch.IsError {
			level = `🔴`
		}

		fmt.Fprintf(&sb, "|%s|%s %s *(%s)*|\n", name, level, glitch.ErrText, glitch.RuleName)
	}

	if len(glitches) > limit && limit > 0 {
		fmt.Fprintf(&sb, "\n...and many more (%d in total).", len(glitches))
	}

	return sb.String()
}
