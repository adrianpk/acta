package acta

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

// App owns process-level command dispatch and dependency assembly.
type App struct {
	config           Config
	stdout           io.Writer
	stderr           io.Writer
	workingDirectory func() (string, error)
	inspector        *GitInspector
}

// NewApp constructs the executable application without performing IO.
func NewApp(config Config, stdout, stderr io.Writer) *App {
	return &App{
		config:           config,
		stdout:           stdout,
		stderr:           stderr,
		workingDirectory: os.Getwd,
		inspector:        NewGitInspector(),
	}
}

// Run executes one CLI command.
func (app *App) Run(ctx context.Context, arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("usage: acta <record|serve>")
	}

	switch arguments[0] {
	case "record":
		return app.runRecord(ctx, arguments[1:])
	case "serve":
		return app.runServe(ctx, arguments[1:])
	default:
		return fmt.Errorf("unknown command %q", arguments[0])
	}
}

func (app *App) runServe(ctx context.Context, arguments []string) error {
	if len(arguments) != 0 {
		return fmt.Errorf("usage: acta serve")
	}

	journal := NewJournal(app.config)
	if err := journal.Load(ctx); err != nil {
		return fmt.Errorf("load Acta journal: %w", err)
	}

	listener, err := ListenServiceSocket(app.config.SocketPath)
	if err != nil {
		return err
	}

	server := NewServer(NewService(NewRecorder(journal)))
	serveErr := server.Serve(ctx, listener)

	closeErr := listener.Close()
	if errors.Is(closeErr, net.ErrClosed) {
		closeErr = nil
	}

	if err := errors.Join(serveErr, closeErr); err != nil {
		return fmt.Errorf("serve Acta requests: %w", err)
	}

	return nil
}
