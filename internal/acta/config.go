// Package acta implements the local Acta journal and its domain contracts.
package acta

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config contains the user-local paths owned by Acta.
type Config struct {
	ConfigRoot            string
	FilesRoot             string
	JournalRoot           string
	PendingPath           string
	VerificationStatePath string
	RuntimeRoot           string
	SocketPath            string
}

// DefaultConfigForCommand skips service paths for the offline migrator.
func DefaultConfigForCommand(command string) (Config, error) {
	if command == "migrate" {
		return Config{}, nil
	}

	return DefaultConfig()
}

// DefaultConfig resolves Acta paths from the current process environment.
func DefaultConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve user home: %w", err)
	}

	config, err := ConfigFromHome(home)
	if err != nil {
		return Config{}, err
	}

	runtimeDirectory := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDirectory == "" {
		return Config{}, fmt.Errorf("resolve Acta runtime paths: XDG_RUNTIME_DIR is required")
	}

	config.RuntimeRoot = filepath.Join(filepath.Clean(runtimeDirectory), "acta")
	config.SocketPath = filepath.Join(config.RuntimeRoot, "acta.sock")

	return config, nil
}

// ConfigFromHome resolves Acta paths for a user home directory.
func ConfigFromHome(home string) (Config, error) {
	if !filepath.IsAbs(home) {
		return Config{}, fmt.Errorf("resolve configuration: home path must be absolute")
	}

	filesRoot := filepath.Join(filepath.Clean(home), ".acta")

	return Config{
		ConfigRoot:            filepath.Join(filepath.Clean(home), ".config", "acta"),
		FilesRoot:             filesRoot,
		JournalRoot:           filepath.Join(filesRoot, "events"),
		PendingPath:           filepath.Join(filesRoot, "pending.json"),
		VerificationStatePath: filepath.Join(filesRoot, "verification.json"),
	}, nil
}
