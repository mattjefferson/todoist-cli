package app

import (
	"context"
	"flag"
	"io"
	"strings"
)

func runTaskUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi task update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var content string
	var description string
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
	fs.BoolVar(&forceID, "id", false, "Treat argument as task ID")
	fs.StringVar(&content, "content", "", "Task content")
	fs.StringVar(&description, "description", "", "Task description")
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
	if len(fs.Args()) == 0 {
		writeLine(state.Err, "error: task identifier required")
		return 2
	}
	identifier := strings.Join(fs.Args(), " ")

	labelsAll := mergeLabels(labels, labelsCSV)
	if err := validateDueFlags(due, dueDate, dueDatetime); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if priority != 0 && (priority < 1 || priority > 4) {
		writeLine(state.Err, "error: priority must be 1-4")
		return 2
	}

	body := map[string]any{}
	if content != "" {
		body["content"] = content
	}
	if description != "" {
		body["description"] = description
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
	if len(body) == 0 {
		writeLine(state.Err, "error: no updates specified")
		return 2
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	id, err := resolveTaskID(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	task, raw, err := client.UpdateTask(ctx, id, body)
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
	if task.ID == "" {
		writeLine(state.Out, "ok")
		return 0
	}
	if err := printTask(state.Out, task, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}
