package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
)

func runTaskClose(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "close", args, false)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "closed %s\n", id); err != nil {
		return 1
	}
	return 0
}

func runTaskReopen(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "reopen", args, false)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "reopened %s\n", id); err != nil {
		return 1
	}
	return 0
}

func runTaskDelete(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "delete", args, true)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", id); err != nil {
		return 1
	}
	return 0
}

func taskAction(ctx context.Context, state *state, action string, args []string, destructive bool) (string, []byte, int) {
	fs := flag.NewFlagSet("todi task "+action, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as task ID")
	if destructive {
		fs.BoolVar(&force, "force", false, "Skip confirmation")
	}
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 2
	}
	if help {
		printTaskUsage(state.Out)
		return "", nil, 0
	}
	if len(fs.Args()) == 0 {
		writeLine(state.Err, "error: task identifier required")
		return "", nil, 2
	}
	identifier := strings.Join(fs.Args(), " ")

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}

	id, err := resolveTaskID(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}

	if destructive {
		if err := confirmDelete(state, "task", identifier, force); err != nil {
			writeLine(state.Err, "error:", err)
			return "", nil, 2
		}
	}

	var raw []byte
	switch action {
	case "close":
		raw, err = client.CloseTask(ctx, id)
	case "reopen":
		raw, err = client.ReopenTask(ctx, id)
	case "delete":
		raw, err = client.DeleteTask(ctx, id)
	default:
		return "", nil, 2
	}
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}
	return id, raw, 0
}
