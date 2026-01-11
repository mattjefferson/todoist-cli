package app

import (
	"context"
	"flag"
	"io"
	"strings"
)

func runTaskGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi task get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as task ID")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printTaskUsage(state.Out)
		return 0
	}
	if len(fs.Args()) == 0 {
		writeLine(state.Err, "error: task identifier required")
		return 2
	}
	identifier := strings.Join(fs.Args(), " ")

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	task, err := resolveTask(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printTask(state.Out, task, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}
