package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func runConfig(_ context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printConfigUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "get":
		return runConfigGet(state, args[1:])
	case "set":
		return runConfigSet(state, args[1:])
	case "path":
		if _, err := fmt.Fprintln(state.Out, state.ConfigPath); err != nil {
			return 1
		}
		return 0
	case "view":
		return runConfigView(state)
	case "-h", "--help", "help":
		printConfigUsage(state.Out)
		return 0
	default:
		if _, err := fmt.Fprintln(state.Err, "error: unknown config command:", args[0]); err != nil {
			return 2
		}
		printConfigUsage(state.Err)
		return 2
	}
}

func runConfigGet(state *state, args []string) int {
	fs := flag.NewFlagSet("todi config get", flag.ContinueOnError)
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
		printConfigUsage(state.Out)
		return 0
	}
	if len(fs.Args()) == 0 {
		if _, err := fmt.Fprintln(state.Err, "error: key required"); err != nil {
			return 2
		}
		return 2
	}
	key := strings.ToLower(fs.Args()[0])
	switch key {
	case "token":
		if _, err := fmt.Fprintln(state.Out, state.Config.Token); err != nil {
			return 1
		}
	case "api_base":
		if _, err := fmt.Fprintln(state.Out, state.Config.APIBase); err != nil {
			return 1
		}
	case "default_project":
		if _, err := fmt.Fprintln(state.Out, state.Config.Project); err != nil {
			return 1
		}
	case "default_labels":
		if _, err := fmt.Fprintln(state.Out, state.Config.Labels); err != nil {
			return 1
		}
	case "label_cli":
		if state.Config.LabelCLI {
			if _, err := fmt.Fprintln(state.Out, "true"); err != nil {
				return 1
			}
		} else {
			if _, err := fmt.Fprintln(state.Out, "false"); err != nil {
				return 1
			}
		}
	default:
		if _, err := fmt.Fprintln(state.Err, "error: unknown key:", key); err != nil {
			return 2
		}
		return 2
	}
	return 0
}

func runConfigSet(state *state, args []string) int {
	fs := flag.NewFlagSet("todi config set", flag.ContinueOnError)
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
		printConfigUsage(state.Out)
		return 0
	}
	if len(fs.Args()) < 2 {
		if _, err := fmt.Fprintln(state.Err, "error: key and value required"); err != nil {
			return 2
		}
		return 2
	}
	key := strings.ToLower(fs.Args()[0])
	value := strings.Join(fs.Args()[1:], " ")
	if key == "token" {
		if _, err := fmt.Fprintln(state.Err, "error: set token via 'todi auth login'"); err != nil {
			return 2
		}
		return 2
	}
	switch key {
	case "api_base":
		state.Config.APIBase = value
	case "default_project":
		state.Config.Project = value
	case "default_labels":
		state.Config.Labels = value
	case "label_cli":
		parsed, err := parseBool(value)
		if err != nil {
			if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
				return 2
			}
			return 2
		}
		state.Config.LabelCLI = parsed
	default:
		if _, err := fmt.Fprintln(state.Err, "error: unknown key:", key); err != nil {
			return 2
		}
		return 2
	}
	if err := state.Config.Save(state.ConfigPath); err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	if _, err := fmt.Fprintln(state.Out, "saved"); err != nil {
		return 1
	}
	return 0
}

func runConfigView(state *state) int {
	data, err := os.ReadFile(state.ConfigPath)
	if err != nil {
		if _, writeErr := fmt.Fprintln(state.Err, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	if _, err := fmt.Fprintln(state.Out, string(data)); err != nil {
		return 1
	}
	return 0
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "t", "1", "yes", "y", "on":
		return true, nil
	case "false", "f", "0", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean: %s", value)
	}
}
