package acta

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

type searchOptions struct {
	query Query
	json  bool
}

type showOptions struct {
	params ShowParams
	json   bool
}

type correctOptions struct {
	params CorrectParams
	json   bool
}

type searchFlagValues struct {
	eventTypes    string
	recordedFrom  string
	recordedTo    string
	authoredFrom  string
	authoredTo    string
	committedFrom string
	committedTo   string
}

type optionalString struct {
	target **string
}

func (app *App) runSearch(ctx context.Context, arguments []string) error {
	options, err := parseSearchOptions(arguments, app.stderr)
	if err != nil {
		return err
	}

	var result QueryResult
	if err := NewClient(app.config.SocketPath).Call(ctx, "search", options.query, &result); err != nil {
		return err
	}

	if options.json {
		return writeJSON(app.stdout, result)
	}

	for _, event := range result.Events {
		if err := writeEventLine(app.stdout, event); err != nil {
			return err
		}
	}

	if result.Truncated {
		_, err = fmt.Fprintf(app.stdout, "truncated after %d events\n", result.ReturnedCount)
	}

	return err
}

func (app *App) runShow(ctx context.Context, arguments []string) error {
	options, err := parseShowOptions(arguments, app.stderr)
	if err != nil {
		return err
	}

	var result ShowResult
	if err := NewClient(app.config.SocketPath).Call(ctx, "show", options.params, &result); err != nil {
		return err
	}

	if options.json {
		return writeJSON(app.stdout, result)
	}

	return writeEventLine(app.stdout, result.Event)
}

func (app *App) runCorrect(ctx context.Context, arguments []string) error {
	options, err := parseCorrectOptions(arguments, app.stderr)
	if err != nil {
		return err
	}

	var result CorrectResult
	if err := NewClient(app.config.SocketPath).Call(ctx, "correct", options.params, &result); err != nil {
		return err
	}

	if options.json {
		return writeJSON(app.stdout, result)
	}

	_, err = fmt.Fprintf(
		app.stdout,
		"corrected %s %s %s\n",
		result.Event.ID,
		result.Event.Repository.Name,
		abbreviatedCommit(result.Event.Git.Commit),
	)

	return err
}

func parseSearchOptions(arguments []string, errorOutput io.Writer) (searchOptions, error) {
	flags := flag.NewFlagSet("search", flag.ContinueOnError)
	flags.SetOutput(errorOutput)

	var (
		options searchOptions
		values  searchFlagValues
	)

	bindSearchFlags(flags, &options, &values)

	if err := flags.Parse(arguments); err != nil {
		return searchOptions{}, err
	}

	if len(flags.Args()) != 0 {
		return searchOptions{}, fmt.Errorf("usage: acta search [flags]")
	}

	if err := applySearchFlagValues(&options.query, values); err != nil {
		return searchOptions{}, err
	}

	return options, nil
}

func bindSearchFlags(flags *flag.FlagSet, options *searchOptions, values *searchFlagValues) {
	flags.StringVar(&options.query.Text, "query", "", "text query")
	flags.StringVar(&options.query.Repository, "repository", "", "repository name or identity")
	flags.StringVar(&options.query.Commit, "commit", "", "full or abbreviated commit")
	flags.StringVar(&options.query.Harness, "harness", "", "execution harness")
	flags.StringVar(&options.query.Model, "model", "", "execution model")
	flags.StringVar(&values.eventTypes, "type", "", "comma-separated event types")
	flags.StringVar(&values.recordedFrom, "recorded-from", "", "inclusive RFC3339 recording time")
	flags.StringVar(&values.recordedTo, "recorded-to", "", "exclusive RFC3339 recording time")
	flags.StringVar(&values.authoredFrom, "authored-from", "", "inclusive RFC3339 author time")
	flags.StringVar(&values.authoredTo, "authored-to", "", "exclusive RFC3339 author time")
	flags.StringVar(&values.committedFrom, "committed-from", "", "inclusive RFC3339 committer time")
	flags.StringVar(&values.committedTo, "committed-to", "", "exclusive RFC3339 committer time")
	flags.IntVar(&options.query.Limit, "limit", 0, "maximum returned events")
	flags.Var(&options.query.Order, "order", "ascending or descending")
	flags.BoolVar(&options.json, "json", false, "write JSON output")
}

func applySearchFlagValues(query *Query, values searchFlagValues) error {
	query.EventTypes = parseEventTypes(values.eventTypes)

	times := []struct {
		name   string
		value  string
		target **time.Time
	}{
		{name: "recorded-from", value: values.recordedFrom, target: &query.RecordedFrom},
		{name: "recorded-to", value: values.recordedTo, target: &query.RecordedTo},
		{name: "authored-from", value: values.authoredFrom, target: &query.AuthoredFrom},
		{name: "authored-to", value: values.authoredTo, target: &query.AuthoredTo},
		{name: "committed-from", value: values.committedFrom, target: &query.CommittedFrom},
		{name: "committed-to", value: values.committedTo, target: &query.CommittedTo},
	}

	for _, input := range times {
		parsed, err := parseOptionalTime(input.name, input.value)
		if err != nil {
			return err
		}

		*input.target = parsed
	}

	return nil
}

func parseShowOptions(arguments []string, errorOutput io.Writer) (showOptions, error) {
	flags := flag.NewFlagSet("show", flag.ContinueOnError)
	flags.SetOutput(errorOutput)

	var options showOptions
	flags.StringVar(&options.params.Commit, "commit", "", "full or abbreviated commit")
	flags.StringVar(&options.params.Repository, "repository", "", "repository name or identity")
	flags.BoolVar(&options.json, "json", false, "write JSON output")

	if err := flags.Parse(arguments); err != nil {
		return showOptions{}, err
	}

	if len(flags.Args()) != 0 || strings.TrimSpace(options.params.Commit) == "" {
		return showOptions{}, fmt.Errorf("usage: acta show --commit <commit> [--repository <repository>] [--json]")
	}

	return options, nil
}

func parseCorrectOptions(arguments []string, errorOutput io.Writer) (correctOptions, error) {
	flags := flag.NewFlagSet("correct", flag.ContinueOnError)
	flags.SetOutput(errorOutput)

	var options correctOptions

	note := optionalString{target: &options.params.ReplacementNote}

	flags.Var(&note, "note", "replacement note")
	flags.StringVar((*string)(&options.params.TargetID), "target", "", "current effective event ID")
	flags.StringVar(&options.params.Reason, "reason", "", "correction reason")
	flags.BoolVar(&options.json, "json", false, "write JSON output")

	if err := flags.Parse(arguments); err != nil {
		return correctOptions{}, err
	}

	if len(flags.Args()) != 0 || options.params.TargetID == "" || strings.TrimSpace(options.params.Reason) == "" {
		return correctOptions{}, fmt.Errorf("usage: acta correct --target <event-id> --reason <reason> [--note <note>] [--json]")
	}

	return options, nil
}

func (value *optionalString) String() string {
	if value == nil || value.target == nil || *value.target == nil {
		return ""
	}

	return **value.target
}

func (value *optionalString) Set(input string) error {
	*value.target = &input

	return nil
}

// Set parses a CLI order value. Query validation rejects unsupported values.
func (order *QueryOrder) Set(value string) error {
	*order = QueryOrder(value)

	return nil
}

func (order *QueryOrder) String() string {
	return string(*order)
}

func parseEventTypes(value string) []EventType {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	eventTypes := make([]EventType, 0, len(parts))

	for _, part := range parts {
		eventTypes = append(eventTypes, EventType(strings.TrimSpace(part)))
	}

	return eventTypes
}

func parseOptionalTime(name, value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("parse --%s: %w", name, err)
	}

	return &parsed, nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)

	return encoder.Encode(value)
}

func writeEventLine(writer io.Writer, event Event) error {
	_, err := fmt.Fprintf(
		writer,
		"%s %s %s %s %q\n",
		event.RecordedAt.Format(time.RFC3339),
		event.Type,
		event.Repository.Name,
		abbreviatedCommit(event.Git.Commit),
		event.Git.Subject,
	)

	return err
}
