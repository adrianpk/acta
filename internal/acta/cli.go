package acta

import (
	"context"
	"flag"
	"fmt"
	"io"
)

type recordOptions struct {
	note           string
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

	params := RecordParams{
		Repository: evidence.Repository,
		Git:        evidence.Git,
		Note:       options.note,
		Execution:  options.execution(),
	}

	var result RecordResult
	if err := NewClient(app.config.SocketPath).Call(ctx, "record", params, &result); err != nil {
		return err
	}

	status := "recorded"
	if !result.Recorded {
		status = "already recorded"
	}

	_, err = fmt.Fprintf(app.stdout, "%s %s %s %s\n", status, result.ID, result.Repository, abbreviatedCommit(result.Commit))

	return err
}

func parseRecordOptions(arguments []string, errorOutput io.Writer) (recordOptions, error) {
	flags := flag.NewFlagSet("record", flag.ContinueOnError)
	flags.SetOutput(errorOutput)

	var options recordOptions
	flags.StringVar(&options.harness, "harness", "", "agent harness name")
	flags.StringVar(&options.harnessVersion, "harness-version", "", "agent harness version")
	flags.StringVar(&options.model, "model", "", "provider-qualified model identifier")

	if err := flags.Parse(arguments); err != nil {
		return recordOptions{}, err
	}

	positional := flags.Args()
	if len(positional) > 1 {
		return recordOptions{}, fmt.Errorf("usage: acta record [flags] [note]")
	}

	if len(positional) == 1 {
		options.note = positional[0]
	}

	if options.harness == "" && (options.harnessVersion != "" || options.model != "") {
		return recordOptions{}, fmt.Errorf("record harness is required with harness version or model")
	}

	return options, nil
}

func (options recordOptions) execution() *ExecutionContext {
	if options.harness == "" {
		return nil
	}

	return &ExecutionContext{
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
