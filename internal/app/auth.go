package app

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runAuth(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printAuthUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "login":
		return runAuthLogin(ctx, state, args[1:])
	case "logout":
		return runAuthLogout(state)
	case "status":
		return runAuthStatus(state)
	case "-h", "--help", "help":
		printAuthUsage(state.Out)
		return 0
	default:
		if _, err := fmt.Fprintln(state.Err, "error: unknown auth command:", args[0]); err != nil {
			return 2
		}
		printAuthUsage(state.Err)
		return 2
	}
}

func runAuthLogin(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi auth login", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	if err := fs.Parse(args); err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 2
		}
		return 2
	}
	if help {
		printAuthUsage(state.Out)
		return 0
	}
	if state.NoInput || !isTTY(os.Stdin) {
		if _, err := fmt.Fprintln(state.Err, "error: login requires TTY (disable --no-input)"); err != nil {
			return 2
		}
		return 2
	}
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fprint(state.Err, "Todoist token: "); err != nil {
		return 1
	}
	token, err := reader.ReadString('\n')
	if err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	token = strings.TrimSpace(token)
	if token == "" {
		if _, err := fmt.Fprintln(state.Err, "error: token required"); err != nil {
			return 2
		}
		return 2
	}

	state.Config.Token = token
	if err := state.Config.Save(state.ConfigPath); err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	if _, err := fmt.Fprintln(state.Out, "token saved"); err != nil {
		return 1
	}
	_ = ctx
	return 0
}

func runAuthLogout(state *state) int {
	state.Config.Token = ""
	if err := state.Config.Save(state.ConfigPath); err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	if _, err := fmt.Fprintln(state.Out, "token cleared"); err != nil {
		return 1
	}
	return 0
}

func runAuthStatus(state *state) int {
	envToken := os.Getenv("TODOIST_TOKEN")
	if envToken != "" {
		if _, err := fmt.Fprintln(state.Out, "token set (TODOIST_TOKEN)"); err != nil {
			return 1
		}
		return 0
	}
	if state.Config.Token != "" {
		if _, err := fmt.Fprintln(state.Out, "token set (config)"); err != nil {
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintln(state.Out, "token missing"); err != nil {
		return 1
	}
	return 3
}
