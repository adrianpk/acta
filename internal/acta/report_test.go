package acta

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestJournalReportGroupsCrossMonthEventsInSuppliedTimezone(t *testing.T) {
	t.Parallel()

	journal := newCrossMonthReportJournal(t)

	result, err := journal.Report(context.Background(), ReportParams{
		From:     "2026-08-31",
		To:       "2026-09-02",
		Timezone: "Europe/Warsaw",
	})
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	assertCrossMonthReport(t, result)
}

func assertCrossMonthReport(t *testing.T, result ReportResult) {
	t.Helper()

	if result.Interval.Timezone != "Europe/Warsaw" || result.ReturnedCount != 3 || result.Truncated {
		t.Fatalf("Report() result = %#v", result)
	}

	if len(result.Days) != 2 || result.Days[0].Date != "2026-08-31" || result.Days[1].Date != "2026-09-01" {
		t.Fatalf("report days = %#v", result.Days)
	}

	septemberRepositories := result.Days[1].Repositories
	if len(septemberRepositories) != 2 || septemberRepositories[0].Name != "alpha" || septemberRepositories[1].Name != "zeta" {
		t.Fatalf("September repositories = %#v", septemberRepositories)
	}
}

func TestJournalReportUsesInclusiveFromAndExclusiveTo(t *testing.T) {
	t.Parallel()

	journal := newCrossMonthReportJournal(t)

	result, err := journal.Report(context.Background(), ReportParams{
		From:     "2026-09-01",
		To:       "2026-09-02",
		Timezone: "Europe/Warsaw",
	})
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	if result.ReturnedCount != 2 || len(result.Days) != 1 || result.Days[0].Date != "2026-09-01" {
		t.Fatalf("Report() result = %#v", result)
	}
}

func TestJournalReportDefaultsToLocalTimezoneAndMaximumLimit(t *testing.T) {
	t.Parallel()

	normalized, err := normalizeReport(ReportParams{From: "2026-08-01", To: "2026-08-02"})
	if err != nil {
		t.Fatalf("normalizeReport() error = %v", err)
	}

	if normalized.location != time.Local || normalized.interval.Timezone != time.Local.String() {
		t.Fatalf("normalized timezone = %#v", normalized.interval)
	}

	if normalized.limit != maximumQueryLimit {
		t.Fatalf("normalized limit = %d, want %d", normalized.limit, maximumQueryLimit)
	}
}

func TestJournalReportValidatesIntervalTimezoneAndLimit(t *testing.T) {
	t.Parallel()

	journal := newCrossMonthReportJournal(t)
	tests := []ReportParams{
		{},
		{From: "2026-09-01", To: "2026-09-01", Timezone: "UTC"},
		{From: "2026-09-01", To: "2026-09-02", Timezone: "Missing/Timezone"},
		{From: "2026-09-01", To: "2026-09-02", Timezone: "UTC", Limit: maximumQueryLimit + 1},
	}

	for _, params := range tests {
		if _, err := journal.Report(context.Background(), params); !errors.Is(err, ErrReportInvalid) {
			t.Fatalf("Report(%#v) error = %v, want ErrReportInvalid", params, err)
		}
	}
}

func TestJournalReportTruncatesByCountAndEncodedSize(t *testing.T) {
	t.Parallel()

	journal := newCrossMonthReportJournal(t)

	limited, err := journal.Report(context.Background(), ReportParams{
		From:     "2026-08-31",
		To:       "2026-09-02",
		Timezone: "Europe/Warsaw",
		Limit:    1,
	})
	if err != nil {
		t.Fatalf("limited Report() error = %v", err)
	}

	if limited.ReturnedCount != 1 || !limited.Truncated {
		t.Fatalf("limited report = %#v", limited)
	}

	const matchCount = 10

	event := validMilestoneEvent()
	event.Git.Subject = strings.Repeat("x", 1<<20)

	events := make([]Event, matchCount)
	for index := range events {
		events[index] = event
	}

	normalized, err := normalizeReport(ReportParams{From: "2026-08-01", To: "2026-08-02", Timezone: "UTC"})
	if err != nil {
		t.Fatalf("normalize report: %v", err)
	}

	bounded := boundedReportResult(normalized, events, false)

	encoded, err := json.Marshal(SuccessResponse(bounded))
	if err != nil {
		t.Fatalf("marshal bounded report: %v", err)
	}

	if !bounded.Truncated || bounded.ReturnedCount >= matchCount || len(encoded)+1 > MaxResponseSize {
		t.Fatalf("bounded report count = %d, truncated = %v, bytes = %d", bounded.ReturnedCount, bounded.Truncated, len(encoded)+1)
	}
}

func newCrossMonthReportJournal(t *testing.T) *Journal {
	t.Helper()

	journal := NewJournal(testJournalConfig(t))
	if err := journal.Load(context.Background()); err != nil {
		t.Fatalf("load journal: %v", err)
	}

	first := distinctMilestoneEvent(t, 11)
	first.RecordedAt = time.Date(2026, 8, 31, 21, 30, 0, 0, time.UTC)
	first.Repository = reportRepository("middle")

	second := distinctMilestoneEvent(t, 12)
	second.RecordedAt = time.Date(2026, 8, 31, 22, 30, 0, 0, time.UTC)
	second.Repository = reportRepository("zeta")

	third := distinctMilestoneEvent(t, 13)
	third.RecordedAt = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	third.Repository = reportRepository("alpha")

	for _, event := range []Event{first, second, third} {
		appendJournalEvent(t, journal, event)
	}

	return journal
}

func reportRepository(name string) RepositoryReference {
	return RepositoryReference{
		Name:     name,
		Root:     "/workspace/" + name,
		Identity: "https://example.test/team/" + name + ".git",
		Remote:   "https://example.test/team/" + name + ".git",
	}
}
