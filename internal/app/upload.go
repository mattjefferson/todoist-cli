package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"
)

func runUpload(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printUploadUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "add":
		return runUploadAdd(ctx, state, args[1:])
	case "delete":
		return runUploadDelete(ctx, state, args[1:])
	case "-h", "--help", "help":
		printUploadUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown upload command:", args[0])
		printUploadUsage(state.Err)
		return 2
	}
}

func runUploadAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist upload add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var projectName string
	var projectID string
	var name string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.StringVar(&name, "name", "", "Override file name")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printUploadUsage(state.Out)
		return 0
	}
	path := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if path == "" {
		writeLine(state.Err, "error: file path required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	projectIDValue, err := resolveProjectID(ctx, client, projectName, projectID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	upload, raw, err := client.UploadFile(ctx, path, name, projectIDValue)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printUpload(state.Out, upload, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runUploadDelete(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist upload delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var fileURL string
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&fileURL, "file-url", "", "File URL")
	fs.BoolVar(&force, "force", false, "Skip confirmation")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printUploadUsage(state.Out)
		return 0
	}
	if len(fs.Args()) > 0 {
		if fileURL != "" {
			writeLine(state.Err, "error: file url specified twice")
			return 2
		}
		fileURL = strings.TrimSpace(strings.Join(fs.Args(), " "))
	}
	if fileURL == "" {
		writeLine(state.Err, "error: file url required")
		return 2
	}

	if err := confirmDelete(state, "upload", fileURL, force); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	raw, err := client.DeleteUpload(ctx, fileURL)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", fileURL); err != nil {
		return 1
	}
	return 0
}
