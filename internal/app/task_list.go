package app

import (
	"context"
	"flag"
	"io"
	"strconv"
	"strings"
)

func runTaskList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi task list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var projectName string
	var limit int
	var cursor string
	var all bool
	var label string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.IntVar(&limit, "limit", 50, "Max tasks per page (1-200)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	fs.StringVar(&label, "label", "", "Label name")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printTaskUsage(state.Out)
		return 0
	}
	if len(fs.Args()) > 0 {
		if projectName != "" {
			writeLine(state.Err, "error: project specified twice")
			return 2
		}
		projectName = strings.Join(fs.Args(), " ")
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	params := map[string]string{}
	if label != "" {
		params["label"] = label
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	if projectName != "" {
		projectID, err := client.FindProjectIDByName(ctx, projectName)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		params["project_id"] = projectID
	}

	if all {
		tasks, err := client.ListTasksAll(ctx, params)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printTasks(state.Out, tasks, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	tasks, next, err := client.ListTasks(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": tasks, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printTasks(state.Out, tasks, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}
