package app

import (
	"context"
	"flag"
	"io"
)

func runUser(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printUserUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "info", "get":
		return runUserInfo(ctx, state, args[1:])
	case "-h", "--help", "help":
		printUserUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown user command:", args[0])
		printUserUsage(state.Err)
		return 2
	}
}

func runUserInfo(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist user info", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printUserUsage(state.Out)
		return 0
	}
	if len(fs.Args()) > 0 {
		writeLine(state.Err, "error: unexpected arguments")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	user, err := client.GetUserInfo(ctx)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printUser(state.Out, user, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}
