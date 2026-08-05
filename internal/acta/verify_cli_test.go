package acta

import (
	"bytes"
	"testing"
	"time"
)

func TestParseVerifyOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		want      verifyOptions
		wantErr   bool
	}{
		{name: "active"},
		{name: "full JSON", arguments: []string{"--full", "--json"}, want: verifyOptions{full: true, json: true}},
		{name: "scheduled", arguments: []string{"--scheduled"}, want: verifyOptions{scheduled: true}},
		{name: "conflicting scopes", arguments: []string{"--full", "--scheduled"}, wantErr: true},
		{name: "positional", arguments: []string{"extra"}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseVerifyOptions(test.arguments, &bytes.Buffer{})
			if test.wantErr && err == nil {
				t.Fatal("expected parse error")
			}

			if !test.wantErr && (err != nil || got != test.want) {
				t.Fatalf("parseVerifyOptions() = %#v, %v", got, err)
			}
		})
	}
}

func TestScheduledVerificationUsesFullScopeOnFirstLocalDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		now  time.Time
		want bool
	}{
		{name: "first day", now: time.Date(2026, 9, 1, 3, 0, 0, 0, time.Local), want: true},
		{name: "later day", now: time.Date(2026, 9, 2, 3, 0, 0, 0, time.Local)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := scheduledVerificationIsFull(test.now); got != test.want {
				t.Fatalf("scheduledVerificationIsFull() = %v, want %v", got, test.want)
			}
		})
	}
}
