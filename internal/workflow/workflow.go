// Package workflow contains code related to GitHub Actions workflow.
package workflow

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	// DotGithubFileTypeWorkflow represents the workflow file type. Used in a bitmask and must be a power of 2.
	DotGithubFileTypeWorkflow = 2
)

// Workflow represents a GitHub Actions' workflow parsed from a YAML file.
type Workflow struct {
	Path        string
	Raw         []byte
	FileName    string
	DisplayName string
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Env         map[string]string `yaml:"env"`
	Jobs        map[string]*Job   `yaml:"jobs"`
	On          *On               `yaml:"on"`
}

// Unmarshal parses YAML from a file in struct's Path or from struct's Raw field.
func (w *Workflow) Unmarshal(fromRaw bool) error {
	if !fromRaw {
		b, err := os.ReadFile(w.Path)
		if err != nil {
			return fmt.Errorf("cannot read file %s: %w", w.Path, err)
		}
		w.Raw = b
	}

	w.FileName = filepath.Base(w.Path)
	ext := filepath.Ext(w.FileName)
	if ext == ".yaml" || ext == ".yml" {
		w.DisplayName = strings.TrimSuffix(w.FileName, ext)
	} else {
		w.DisplayName = w.FileName
	}

	slog.Debug(
		"reading workflow file",
		slog.String("path", w.Path),
	)

	err := yaml.Unmarshal(w.Raw, &w)
	if err != nil {
		return fmt.Errorf("cannot unmarshal file %s: %w", w.Path, err)
	}

	if w.Jobs != nil {
		for _, j := range w.Jobs {
			j.SetParentType("workflow")
		}
	}

	return nil
}

// GetType returns the int value representing the workflow file type. See dotgithub.File interface.
func (w *Workflow) GetType() int {
	return DotGithubFileTypeWorkflow
}
