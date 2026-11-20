package linter

import "path/filepath"

// Paths contains paths to exclude and include.
type Paths struct {
	NoChecking []string `yaml:"no_checking,omitempty"`
	Checking   []string `yaml:"checking,omitempty"`
}

// Check checks if a specified path should be included for linting.
func (p *Paths) Check(path string) bool {
	if p == nil {
		return true
	}

	if len(p.NoChecking) > 0 {
		for _, pattern := range p.NoChecking {
			match, _ := filepath.Match(pattern, path)
			if match {
				for _, allowPattern := range p.Checking {
					allowMatch, _ := filepath.Match(allowPattern, path)
					if allowMatch {
						return true
					}
				}

				return false
			}
		}
	}

	if len(p.Checking) > 0 {
		for _, allowPattern := range p.Checking {
			allowMatch, _ := filepath.Match(allowPattern, path)
			if allowMatch {
				return true
			}
		}

		return false
	}

	return true
}
