package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mattjefferson/todoist-cli/internal/todoist"
)

func runLabel(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printLabelUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runLabelList(ctx, state, args[1:])
	case "get":
		return runLabelGet(ctx, state, args[1:])
	case "add":
		return runLabelAdd(ctx, state, args[1:])
	case "update":
		return runLabelUpdate(ctx, state, args[1:])
	case "delete":
		return runLabelDelete(ctx, state, args[1:])
	case "-h", "--help", "help":
		printLabelUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown label command:", args[0])
		printLabelUsage(state.Err)
		return 2
	}
}

func runLabelList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist label list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var limit int
	var cursor string
	var all bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.IntVar(&limit, "limit", 50, "Max labels per page (1-200)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printLabelUsage(state.Out)
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

	params := map[string]string{}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}

	if all {
		labels, err := client.ListLabelsAll(ctx, params)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printLabels(state.Out, labels, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	labels, next, err := client.ListLabels(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": labels, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printLabels(state.Out, labels, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runLabelGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist label get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as label ID")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printLabelUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: label identifier required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	label, err := resolveLabel(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printLabel(state.Out, label, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runLabelAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist label add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var color string
	var favorite bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&color, "color", "", "Label color")
	fs.BoolVar(&favorite, "favorite", false, "Favorite label")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printLabelUsage(state.Out)
		return 0
	}
	name := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if name == "" {
		writeLine(state.Err, "error: label name required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	body := map[string]any{"name": name}
	if color != "" {
		body["color"] = color
	}
	if favorite {
		body["is_favorite"] = true
	}

	label, raw, err := client.CreateLabel(ctx, body)
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
	if err := printLabel(state.Out, label, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runLabelUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist label update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var name string
	var color string
	var favorite bool
	var unfavorite bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as label ID")
	fs.StringVar(&name, "name", "", "Label name")
	fs.StringVar(&color, "color", "", "Label color")
	fs.BoolVar(&favorite, "favorite", false, "Favorite label")
	fs.BoolVar(&unfavorite, "unfavorite", false, "Remove favorite")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printLabelUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: label identifier required")
		return 2
	}
	if favorite && unfavorite {
		writeLine(state.Err, "error: cannot use --favorite and --unfavorite together")
		return 2
	}

	body := map[string]any{}
	if name != "" {
		body["name"] = name
	}
	if color != "" {
		body["color"] = color
	}
	if favorite {
		body["is_favorite"] = true
	}
	if unfavorite {
		body["is_favorite"] = false
	}
	if len(body) == 0 {
		writeLine(state.Err, "error: no updates specified")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	labelID, err := resolveLabelIDFromIdentifier(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	label, raw, err := client.UpdateLabel(ctx, labelID, body)
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
	if label.ID == "" {
		writeLine(state.Out, "ok")
		return 0
	}
	if err := printLabel(state.Out, label, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runLabelDelete(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist label delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as label ID")
	fs.BoolVar(&force, "force", false, "Skip confirmation")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printLabelUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: label identifier required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	labelID, err := resolveLabelIDFromIdentifier(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := confirmDelete(state, "label", identifier, force); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	raw, err := client.DeleteLabel(ctx, labelID)
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
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", labelID); err != nil {
		return 1
	}
	return 0
}

func resolveLabel(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (todoist.Label, error) {
	if forceID {
		return client.GetLabel(ctx, identifier)
	}
	return client.FindLabelByName(ctx, identifier)
}

func resolveLabelIDFromIdentifier(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (string, error) {
	if forceID {
		return identifier, nil
	}
	return client.FindLabelIDByName(ctx, identifier)
}
