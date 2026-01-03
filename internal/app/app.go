package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattjefferson/todi/internal/config"
	"github.com/mattjefferson/todi/internal/todi"
)

const defaultAPIBase = "https://api.todoist.com"

// Run executes the CLI entrypoint with the provided args.
func Run(args []string) int {
	ctx := context.Background()
	out := os.Stdout
	errOut := os.Stderr

	globals, rest, code := parseGlobal(args, errOut)
	if code != 0 {
		return code
	}
	if globals.ShowVersion {
		if _, err := fmt.Fprintln(out, VersionString()); err != nil {
			return 1
		}
		return 0
	}
	if globals.ShowHelp || len(rest) == 0 {
		printUsage(out)
		if len(rest) == 0 {
			return 2
		}
		return 0
	}

	configPath := globals.ConfigPath
	if configPath == "" {
		path, err := config.DefaultPath()
		if err != nil {
			if _, writeErr := fmt.Fprintln(errOut, "error:", err); writeErr != nil {
				return 1
			}
			return 1
		}
		configPath = path
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		if _, writeErr := fmt.Fprintln(errOut, "error:", err); writeErr != nil {
			return 1
		}
		return 1
	}

	if globals.APIBase != "" {
		cfg.APIBase = globals.APIBase
	}
	if cfg.APIBase == "" {
		cfg.APIBase = envOrDefault("TODOIST_API_BASE", defaultAPIBase)
	}

	mode, err := parseOutputMode(globals.JSON, globals.Plain)
	if err != nil {
		if _, writeErr := fmt.Fprintln(errOut, "error:", err); writeErr != nil {
			return 2
		}
		return 2
	}

	state := &state{
		Out:        out,
		Err:        errOut,
		Mode:       mode,
		NoInput:    globals.NoInput,
		Quiet:      globals.Quiet,
		Verbose:    globals.Verbose,
		Config:     cfg,
		ConfigPath: configPath,
		LabelCLI:   globals.LabelCLI || cfg.LabelCLI,
	}

	switch rest[0] {
	case "task":
		return runTask(ctx, state, rest[1:])
	case "project":
		return runProject(ctx, state, rest[1:])
	case "comment":
		return runComment(ctx, state, rest[1:])
	case "label":
		return runLabel(ctx, state, rest[1:])
	case "activity":
		return runActivity(ctx, state, rest[1:])
	case "upload":
		return runUpload(ctx, state, rest[1:])
	case "section":
		return runSection(ctx, state, rest[1:])
	case "user":
		return runUser(ctx, state, rest[1:])
	case "auth":
		return runAuth(ctx, state, rest[1:])
	case "config":
		return runConfig(ctx, state, rest[1:])
	case "help", "-h", "--help":
		printUsage(out)
		return 0
	default:
		if isTaskSubcommand(rest[0]) {
			return runTask(ctx, state, rest)
		}
		if _, writeErr := fmt.Fprintln(errOut, "error: unknown command:", rest[0]); writeErr != nil {
			return 2
		}
		printUsage(errOut)
		return 2
	}
}

type globalFlags struct {
	ShowHelp    bool
	ShowVersion bool
	JSON        bool
	Plain       bool
	Quiet       bool
	Verbose     bool
	NoInput     bool
	NoColor     bool
	ConfigPath  string
	APIBase     string
	LabelCLI    bool
}

type state struct {
	Out        io.Writer
	Err        io.Writer
	Mode       outputMode
	NoInput    bool
	Quiet      bool
	Verbose    bool
	Config     *config.Config
	ConfigPath string
	LabelCLI   bool
}

func (s *state) client() (*todi.Client, error) {
	token := firstNonEmpty(os.Getenv("TODOIST_TOKEN"), s.Config.Token)
	if token == "" {
		return nil, errors.New("missing Todoist token: run 'todi auth login' or set TODOIST_TOKEN")
	}
	client := todi.NewClient(s.Config.APIBase, token, s.Verbose)
	return client, nil
}

func parseGlobal(args []string, errOut io.Writer) (globalFlags, []string, int) {
	fs := flag.NewFlagSet("todi", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var flags globalFlags
	fs.BoolVar(&flags.ShowHelp, "help", false, "Show help")
	fs.BoolVar(&flags.ShowHelp, "h", false, "Show help")
	fs.BoolVar(&flags.ShowVersion, "version", false, "Show version")
	fs.BoolVar(&flags.JSON, "json", false, "JSON output")
	fs.BoolVar(&flags.Plain, "plain", false, "Plain output")
	fs.BoolVar(&flags.Quiet, "quiet", false, "Quiet output")
	fs.BoolVar(&flags.Quiet, "q", false, "Quiet output")
	fs.BoolVar(&flags.Verbose, "verbose", false, "Verbose output")
	fs.BoolVar(&flags.Verbose, "v", false, "Verbose output")
	fs.BoolVar(&flags.NoInput, "no-input", false, "Disable prompts")
	fs.BoolVar(&flags.NoColor, "no-color", false, "Disable color")
	fs.StringVar(&flags.ConfigPath, "config", "", "Config path")
	fs.StringVar(&flags.APIBase, "api-base", "", "API base URL")
	fs.BoolVar(&flags.LabelCLI, "label-cli", false, "Add label 'cli' to created tasks")

	if err := fs.Parse(args); err != nil {
		if _, writeErr := fmt.Fprintln(errOut, "error:", err); writeErr != nil {
			return globalFlags{}, nil, 2
		}
		printUsage(errOut)
		return globalFlags{}, nil, 2
	}
	return flags, fs.Args(), 0
}

func envOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func parseOutputMode(json, plain bool) (outputMode, error) {
	if json && plain {
		return modeHuman, fmt.Errorf("cannot use --json and --plain together")
	}
	if json {
		return modeJSON, nil
	}
	if plain {
		return modePlain, nil
	}
	return modeHuman, nil
}

func joinArgs(args []string) string {
	return strings.TrimSpace(strings.Join(args, " "))
}

func isTaskSubcommand(arg string) bool {
	switch arg {
	case "list", "get", "add", "update", "close", "reopen", "delete", "quick":
		return true
	default:
		return false
	}
}

func isTTY(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
