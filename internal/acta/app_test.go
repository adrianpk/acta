package acta

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAppRecordEndToEndAcrossServiceRestart(t *testing.T) {
	t.Parallel()

	repositoryRoot := initializeTestRepository(t)
	config := testAppConfig(t)

	stopFirst := startTestAppService(t, config)
	firstOutput := runTestAppRecord(t, config, repositoryRoot, []string{
		"record",
		"--harness", "codex",
		"--harness-version", "1.2.3",
		"--model", "openai/example-model",
		"Reusable integration behavior",
	})

	stopFirst()

	if !strings.HasPrefix(firstOutput, "recorded ") {
		t.Fatalf("first output = %q", firstOutput)
	}

	stopSecond := startTestAppService(t, config)
	secondOutput := runTestAppRecord(t, config, repositoryRoot, []string{"record"})

	stopSecond()

	if !strings.HasPrefix(secondOutput, "already recorded ") {
		t.Fatalf("second output = %q", secondOutput)
	}

	journal := NewJournal(config)
	if err := journal.Load(context.Background()); err != nil {
		t.Fatalf("load final journal: %v", err)
	}

	if journal.Count() != 1 {
		t.Fatalf("journal count = %d, want 1", journal.Count())
	}
}

func TestAppRecordDoesNotFallbackWithoutService(t *testing.T) {
	t.Parallel()

	repositoryRoot := initializeTestRepository(t)
	config := testAppConfig(t)

	statusBefore := testGitOutput(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=normal")
	app := NewApp(config, &bytes.Buffer{}, &bytes.Buffer{})
	app.workingDirectory = func() (string, error) { return repositoryRoot, nil }

	if err := app.Run(context.Background(), []string{"record", "Significant behavior"}); err == nil {
		t.Fatal("expected unavailable service error")
	}

	if _, err := os.Stat(config.JournalRoot); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("journal root exists without service: %v", err)
	}

	statusAfter := testGitOutput(t, repositoryRoot, "status", "--porcelain=v1", "--untracked-files=normal")
	if statusAfter != statusBefore {
		t.Fatalf("record changed Git status from %q to %q", statusBefore, statusAfter)
	}
}

func TestParseRecordOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		want      recordOptions
		wantErr   bool
	}{
		{
			name:      "note only",
			arguments: []string{"Reusable behavior"},
			want:      recordOptions{note: "Reusable behavior"},
		},
		{
			name:      "execution",
			arguments: []string{"--harness", "pi", "--harness-version", "2.0", "--model", "provider/model"},
			want: recordOptions{
				harness:        "pi",
				harnessVersion: "2.0",
				model:          "provider/model",
			},
		},
		{
			name:      "too many notes",
			arguments: []string{"one", "two"},
			wantErr:   true,
		},
		{
			name:      "model without harness",
			arguments: []string{"--model", "provider/model"},
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseRecordOptions(test.arguments, &bytes.Buffer{})
			if test.wantErr && err == nil {
				t.Fatal("expected parse error")
			}

			if !test.wantErr && (err != nil || got != test.want) {
				t.Fatalf("parseRecordOptions() = %#v, %v", got, err)
			}
		})
	}
}

func TestAppRejectsUnknownCommand(t *testing.T) {
	t.Parallel()

	app := NewApp(Config{}, &bytes.Buffer{}, &bytes.Buffer{})
	if err := app.Run(context.Background(), []string{"missing"}); err == nil {
		t.Fatal("expected unknown command error")
	}
}

func testAppConfig(t *testing.T) Config {
	t.Helper()

	config := testJournalConfig(t)
	config.RuntimeRoot = filepath.Join(t.TempDir(), "runtime", "acta")
	config.SocketPath = filepath.Join(config.RuntimeRoot, "acta.sock")

	return config
}

func startTestAppService(t *testing.T, config Config) func() {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	app := NewApp(config, &bytes.Buffer{}, &bytes.Buffer{})
	go func() {
		done <- app.Run(ctx, []string{"serve"})
	}()

	waitForSocket(t, config.SocketPath)

	var once sync.Once

	stop := func() {
		once.Do(func() {
			cancel()

			if err := <-done; err != nil {
				t.Errorf("serve app: %v", err)
			}
		})
	}

	t.Cleanup(stop)

	return stop
}

func waitForSocket(t *testing.T, socketPath string) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)

	for {
		if _, err := os.Stat(socketPath); err == nil {
			return
		}

		if time.Now().After(deadline) {
			t.Fatalf("socket %q was not created", socketPath)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func runTestAppRecord(t *testing.T, config Config, repositoryRoot string, arguments []string) string {
	t.Helper()

	var output bytes.Buffer

	app := NewApp(config, &output, &bytes.Buffer{})
	app.workingDirectory = func() (string, error) { return repositoryRoot, nil }

	if err := app.Run(context.Background(), arguments); err != nil {
		t.Fatalf("record app: %v", err)
	}

	return output.String()
}
