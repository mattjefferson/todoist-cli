package app

import (
	"context"
	"flag"
	"io"
)

func runTaskAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi task add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var description string
	var projectName string
	var projectID string
	var labels stringSlice
	var labelsCSV string
	var priority int
	var assignee string
	var due string
	var dueDate string
	var dueDatetime string
	var dueLang string
	var duration int
	var durationUnit string
	var deadlineDate string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&description, "description", "", "Task description")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.Var(&labels, "label", "Label (repeatable)")
	fs.StringVar(&labelsCSV, "labels", "", "Labels (comma-separated)")
	fs.IntVar(&priority, "priority", 0, "Priority 1-4")
	fs.StringVar(&assignee, "assignee", "", "Assignee ID")
	fs.StringVar(&due, "due", "", "Due string")
	fs.StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	fs.StringVar(&dueDatetime, "due-datetime", "", "Due datetime (RFC3339)")
	fs.StringVar(&dueLang, "due-lang", "", "Due language code")
	fs.IntVar(&duration, "duration", 0, "Duration value")
	fs.StringVar(&durationUnit, "duration-unit", "", "Duration unit (minute|day)")
	fs.StringVar(&deadlineDate, "deadline-date", "", "Deadline date (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printTaskUsage(state.Out)
		return 0
	}
	content := joinArgs(fs.Args())
	if content == "" {
		writeLine(state.Err, "error: content required")
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

	labelsAll := mergeLabels(labels, labelsCSV)
	if state.LabelCLI {
		labelsAll = appendUniqueLabel(labelsAll, cliLabel)
	}
	if err := validateDueFlags(due, dueDate, dueDatetime); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if priority != 0 && (priority < 1 || priority > 4) {
		writeLine(state.Err, "error: priority must be 1-4")
		return 2
	}

	body := map[string]any{"content": content}
	if description != "" {
		body["description"] = description
	}
	if projectIDValue != "" {
		body["project_id"] = projectIDValue
	}
	if len(labelsAll) > 0 {
		body["labels"] = labelsAll
	}
	if priority != 0 {
		body["priority"] = priority
	}
	if assignee != "" {
		body["assignee_id"] = assignee
	}
	if due != "" {
		body["due_string"] = due
	}
	if dueDate != "" {
		body["due_date"] = dueDate
	}
	if dueDatetime != "" {
		body["due_datetime"] = dueDatetime
	}
	if dueLang != "" {
		body["due_lang"] = dueLang
	}
	if duration != 0 {
		body["duration"] = duration
	}
	if durationUnit != "" {
		body["duration_unit"] = durationUnit
	}
	if deadlineDate != "" {
		body["deadline_date"] = deadlineDate
	}

	task, raw, err := client.CreateTask(ctx, body)
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
	if err := printTask(state.Out, task, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}
