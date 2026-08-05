// Package acta implements the local Acta journal and its domain contracts.
package acta

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config contains the user-local paths owned by Acta.
type Config struct {
	ConfigRoot  string
	FilesRoot   string
	JournalRoot string
	PendingPath string
}

// DefaultConfig resolves Acta paths from the current process environment.
func DefaultConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve user home: %w", err)
	}

	return ConfigFromHome(home)
}

// ConfigFromHome resolves Acta paths for a user home directory.
func ConfigFromHome(home string) (Config, error) {
	if !filepath.IsAbs(home) {
		return Config{}, fmt.Errorf("resolve configuration: home path must be absolute")
	}

	filesRoot := filepath.Join(filepath.Clean(home), ".acta")

	return Config{
		ConfigRoot:  filepath.Join(filepath.Clean(home), ".config", "acta"),
		FilesRoot:   filesRoot,
		JournalRoot: filepath.Join(filesRoot, "events"),
		PendingPath: filepath.Join(filesRoot, "pending.json"),
	}, nil
}
