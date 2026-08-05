package acta

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestParseSearchOptions(t *testing.T) {
	t.Parallel()

	options, err := parseSearchOptions([]string{
		"--query", "runtime",
		"--repository", "example",
		"--type", "correction,milestone",
		"--recorded-from", "2026-08-01T00:00:00+02:00",
		"--limit", "12",
		"--order", "ascending",
		"--json",
	}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("parseSearchOptions() error = %v", err)
	}

	wantTime := time.Date(2026, 8, 1, 0, 0, 0, 0, time.FixedZone("offset", 2*60*60))

	valid := options.json &&
		options.query.Text == "runtime" &&
		options.query.Repository == "example" &&
		len(options.query.EventTypes) == 2 &&
		options.query.RecordedFrom.Equal(wantTime) &&
		options.query.Limit == 12 &&
		options.query.Order == QueryAscending
	if !valid {
		t.Fatalf("parseSearchOptions() = %#v", options)
	}
}

func TestParseRetrievalOptionsRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		parse func() error
	}{
		{
			name: "search timestamp",
			parse: func() error {
				_, err := parseSearchOptions([]string{"--recorded-from", "today"}, &bytes.Buffer{})

				return err
			},
		},
		{
			name: "show commit",
			parse: func() error {
				_, err := parseShowOptions(nil, &bytes.Buffer{})

				return err
			},
		},
		{
			name: "correction reason",
			parse: func() error {
				_, err := parseCorrectOptions([]string{"--target", "01K1Z8V7BMMN6W7Y6R8QJHA9SS"}, &bytes.Buffer{})

				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if err := test.parse(); err == nil {
				t.Fatal("expected parse error")
			}
		})
	}
}

func TestParseCorrectOptionsDistinguishesAbsentAndEmptyNote(t *testing.T) {
	t.Parallel()

	base := []string{
		"--target", "01K1Z8V7BMMN6W7Y6R8QJHA9SS",
		"--reason", "Clarify behavior",
	}

	absent, err := parseCorrectOptions(base, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("parse absent note: %v", err)
	}

	empty, err := parseCorrectOptions(append(base, "--note", ""), &bytes.Buffer{})
	if err != nil {
		t.Fatalf("parse empty note: %v", err)
	}

	if absent.params.ReplacementNote != nil || empty.params.ReplacementNote == nil || *empty.params.ReplacementNote != "" {
		t.Fatalf("replacement notes = %#v, %#v", absent.params.ReplacementNote, empty.params.ReplacementNote)
	}
}

func TestAppRetrievalAndCorrectionEndToEnd(t *testing.T) {
	t.Parallel()

	repositoryRoot := initializeTestRepository(t)
	config := testAppConfig(t)

	stop := startTestAppService(t, config)
	defer stop()

	runTestAppRecord(t, config, repositoryRoot, []string{"record", "Reusable integration behavior"})

	var initial QueryResult
	if err := NewClient(config.SocketPath).Call(context.Background(), "search", Query{Limit: 10}, &initial); err != nil {
		t.Fatalf("search recorded event: %v", err)
	}

	if len(initial.Events) != 1 {
		t.Fatalf("initial events = %#v", initial.Events)
	}

	event := initial.Events[0]
	showOutput := runTestAppCommand(t, config, []string{"show", "--commit", event.Git.Commit[:12], "--json"})

	var shown ShowResult
	decodeCLIResult(t, showOutput, &shown)
	assertShownEvent(t, shown, event.ID)

	searchOutput := runTestAppCommand(t, config, []string{"search", "--query", "integration", "--limit", "1", "--json"})

	var searched QueryResult
	decodeCLIResult(t, searchOutput, &searched)
	assertSearchedCount(t, searched, 1)

	correctOutput := runTestAppCommand(t, config, []string{
		"correct",
		"--target", string(event.ID),
		"--reason", "Clarify integration behavior",
		"--note", "Corrected integration behavior",
		"--json",
	})

	var corrected CorrectResult
	decodeCLIResult(t, correctOutput, &corrected)
	assertCorrectedEvent(t, corrected, event.ID, "Corrected integration behavior")

	if !strings.Contains(runTestAppCommand(t, config, []string{"search"}), "correction") {
		t.Fatal("human search output does not identify the effective correction")
	}
}

func decodeCLIResult(t *testing.T, output string, target any) {
	t.Helper()

	if err := json.Unmarshal([]byte(output), target); err != nil {
		t.Fatalf("decode CLI output %q: %v", output, err)
	}
}

func assertShownEvent(t *testing.T, result ShowResult, wantID EventID) {
	t.Helper()

	if result.Event.ID != wantID {
		t.Fatalf("shown event = %#v", result.Event)
	}
}

func assertSearchedCount(t *testing.T, result QueryResult, want int) {
	t.Helper()

	if result.ReturnedCount != want {
		t.Fatalf("searched result = %#v", result)
	}
}

func assertCorrectedEvent(t *testing.T, result CorrectResult, targetID EventID, note string) {
	t.Helper()

	if result.Event.TargetID != targetID || result.Event.Note != note {
		t.Fatalf("corrected event = %#v", result.Event)
	}
}

func runTestAppCommand(t *testing.T, config Config, arguments []string) string {
	t.Helper()

	var output bytes.Buffer

	app := NewApp(config, &output, &bytes.Buffer{})

	if err := app.Run(context.Background(), arguments); err != nil {
		t.Fatalf("run %q: %v", arguments[0], err)
	}

	return output.String()
}
