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

func runProject(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printProjectUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runProjectList(ctx, state, args[1:])
	case "get":
		return runProjectGet(ctx, state, args[1:])
	case "add":
		return runProjectAdd(ctx, state, args[1:])
	case "update":
		return runProjectUpdate(ctx, state, args[1:])
	case "delete":
		return runProjectDelete(ctx, state, args[1:])
	case "-h", "--help", "help":
		printProjectUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown project command:", args[0])
		printProjectUsage(state.Err)
		return 2
	}
}

func runProjectList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist project list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var limit int
	var cursor string
	var all bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.IntVar(&limit, "limit", 50, "Max projects per page (1-200)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printProjectUsage(state.Out)
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
		projects, err := client.ListProjectsAll(ctx)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printProjects(state.Out, projects, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	projects, next, err := client.ListProjects(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": projects, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printProjects(state.Out, projects, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runProjectGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist project get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as project ID")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printProjectUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: project identifier required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	project, err := resolveProject(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printProject(state.Out, project, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runProjectAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist project add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var parentName string
	var parentID string
	var color string
	var favorite bool
	var viewStyle string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&parentName, "parent", "", "Parent project title (exact match)")
	fs.StringVar(&parentID, "parent-id", "", "Parent project ID")
	fs.StringVar(&color, "color", "", "Project color")
	fs.BoolVar(&favorite, "favorite", false, "Favorite project")
	fs.StringVar(&viewStyle, "view", "", "View style (list|board)")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printProjectUsage(state.Out)
		return 0
	}
	name := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if name == "" {
		writeLine(state.Err, "error: project name required")
		return 2
	}
	if viewStyle != "" {
		viewStyle = strings.ToLower(viewStyle)
	}
	if viewStyle != "" && !validViewStyle(viewStyle) {
		writeLine(state.Err, "error: view must be list or board")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	parentIDValue, err := resolveParentProjectID(ctx, client, parentName, parentID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	body := map[string]any{"name": name}
	if parentIDValue != "" {
		body["parent_id"] = parentIDValue
	}
	if color != "" {
		body["color"] = color
	}
	if favorite {
		body["is_favorite"] = true
	}
	if viewStyle != "" {
		body["view_style"] = viewStyle
	}

	project, raw, err := client.CreateProject(ctx, body)
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
	if err := printProject(state.Out, project, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runProjectUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist project update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var name string
	var color string
	var favorite bool
	var unfavorite bool
	var viewStyle string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as project ID")
	fs.StringVar(&name, "name", "", "Project name")
	fs.StringVar(&color, "color", "", "Project color")
	fs.BoolVar(&favorite, "favorite", false, "Favorite project")
	fs.BoolVar(&unfavorite, "unfavorite", false, "Remove favorite")
	fs.StringVar(&viewStyle, "view", "", "View style (list|board)")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printProjectUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: project identifier required")
		return 2
	}
	if viewStyle != "" {
		viewStyle = strings.ToLower(viewStyle)
	}
	if favorite && unfavorite {
		writeLine(state.Err, "error: cannot use --favorite and --unfavorite together")
		return 2
	}
	if viewStyle != "" && !validViewStyle(viewStyle) {
		writeLine(state.Err, "error: view must be list or board")
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
	if viewStyle != "" {
		body["view_style"] = viewStyle
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

	projectID, err := resolveProjectIDFromIdentifier(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	project, raw, err := client.UpdateProject(ctx, projectID, body)
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
	if project.ID == "" {
		writeLine(state.Out, "ok")
		return 0
	}
	if err := printProject(state.Out, project, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runProjectDelete(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist project delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as project ID")
	fs.BoolVar(&force, "force", false, "Skip confirmation")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printProjectUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: project identifier required")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	projectID, err := resolveProjectIDFromIdentifier(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := confirmDelete(state, "project", identifier, force); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	raw, err := client.DeleteProject(ctx, projectID)
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
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", projectID); err != nil {
		return 1
	}
	return 0
}

func resolveProject(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (todoist.Project, error) {
	if forceID {
		return client.GetProject(ctx, identifier)
	}
	return client.FindProjectByName(ctx, identifier)
}

func resolveProjectIDFromIdentifier(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (string, error) {
	if forceID {
		return identifier, nil
	}
	return client.FindProjectIDByName(ctx, identifier)
}

func resolveParentProjectID(ctx context.Context, client *todoist.Client, parentName, parentID string) (string, error) {
	if parentName != "" && parentID != "" {
		return "", fmt.Errorf("cannot use --parent and --parent-id together")
	}
	if parentID != "" {
		return parentID, nil
	}
	if parentName == "" {
		return "", nil
	}
	return client.FindProjectIDByName(ctx, parentName)
}

func validViewStyle(value string) bool {
	switch value {
	case "list", "board":
		return true
	default:
		return false
	}
}
