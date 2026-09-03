package acta

import (
	"context"
	"flag"
	"fmt"
	"io"
)

type recordOptions struct {
	summary        string
	harness        string
	harnessVersion string
	model          string
}

func (app *App) runRecord(ctx context.Context, arguments []string) error {
	options, err := parseRecordOptions(arguments, app.stderr)
	if err != nil {
		return err
	}

	directory, err := app.workingDirectory()
	if err != nil {
		return fmt.Errorf("resolve current directory: %w", err)
	}

	evidence, err := app.inspector.Inspect(ctx, directory)
	if err != nil {
		return err
	}

	params := commitRecordParams(evidence, options)

	var result RecordResult
	if err := NewClient(app.config.SocketPath).Call(ctx, "record", params, &result); err != nil {
		return err
	}

	status := "recorded"
	if !result.Recorded {
		status = "already recorded"
	}

	_, err = fmt.Fprintf(app.stdout, "%s %s %s %s %s\n", status, result.ID, result.Type, result.ScopeName, result.Identifier)

	return err
}

func parseRecordOptions(arguments []string, errorOutput io.Writer) (recordOptions, error) {
	if len(arguments) == 0 || arguments[0] != string(ActivityTypeCommit) {
		return recordOptions{}, fmt.Errorf("usage: acta record commit [flags] [summary]")
	}

	flags := flag.NewFlagSet("record commit", flag.ContinueOnError)
	flags.SetOutput(errorOutput)

	var options recordOptions
	flags.StringVar(&options.harness, "harness", "", "agent harness name")
	flags.StringVar(&options.harnessVersion, "harness-version", "", "agent harness version")
	flags.StringVar(&options.model, "model", "", "provider-qualified model identifier")

	if err := flags.Parse(arguments[1:]); err != nil {
		return recordOptions{}, err
	}

	positional := flags.Args()
	if len(positional) > 1 {
		return recordOptions{}, fmt.Errorf("usage: acta record commit [flags] [summary]")
	}

	if len(positional) == 1 {
		options.summary = positional[0]
	}

	if options.harness == "" && (options.harnessVersion != "" || options.model != "") {
		return recordOptions{}, fmt.Errorf("record harness is required with harness version or model")
	}

	return options, nil
}

func commitRecordParams(evidence RepositoryEvidence, options recordOptions) RecordParams {
	summary := options.summary
	if summary == "" {
		summary = evidence.Git.Subject
	}

	properties := &ActivityProperties{
		ObjectID: evidence.Git.Commit,
		Subject:  evidence.Git.Subject,
	}
	if !evidence.Git.AuthoredAt.Equal(evidence.Git.CommittedAt) {
		authoredAt := evidence.Git.AuthoredAt
		properties.AuthoredAt = &authoredAt
	}

	return RecordParams{
		Type:       ActivityTypeCommit,
		OccurredAt: evidence.Git.CommittedAt,
		Summary:    summary,
		Scope: &ActivityScope{
			Type:     ScopeTypeRepository,
			Name:     evidence.Repository.Name,
			Identity: evidence.Repository.Identity,
		},
		Properties: properties,
		Provenance: options.provenance(),
	}
}

func (options recordOptions) provenance() *ActivityProvenance {
	if options.harness == "" {
		return nil
	}

	return &ActivityProvenance{
		Type:           ProvenanceTypeAgent,
		Harness:        options.harness,
		HarnessVersion: options.harnessVersion,
		Model:          options.model,
	}
}

func abbreviatedCommit(commit string) string {
	const abbreviationLength = 12

	if len(commit) <= abbreviationLength {
		return commit
	}

	return commit[:abbreviationLength]
}
