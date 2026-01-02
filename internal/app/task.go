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

type stringSlice []string

const cliLabel = "cli"

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func runTask(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printTaskUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runTaskList(ctx, state, args[1:])
	case "get":
		return runTaskGet(ctx, state, args[1:])
	case "add":
		return runTaskAdd(ctx, state, args[1:])
	case "update":
		return runTaskUpdate(ctx, state, args[1:])
	case "close":
		return runTaskClose(ctx, state, args[1:])
	case "reopen":
		return runTaskReopen(ctx, state, args[1:])
	case "delete":
		return runTaskDelete(ctx, state, args[1:])
	case "quick":
		return runTaskQuick(ctx, state, args[1:])
	case "-h", "--help", "help":
		printTaskUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown task command:", args[0])
		printTaskUsage(state.Err)
		return 2
	}
}

func runTaskList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist task list", flag.ContinueOnError)
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

func runTaskGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist task get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as task ID")
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

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	task, err := resolveTask(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printTask(state.Out, task, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runTaskAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist task add", flag.ContinueOnError)
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

func runTaskUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist task update", flag.ContinueOnError)
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

func runTaskClose(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "close", args, false)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "closed %s\n", id); err != nil {
		return 1
	}
	return 0
}

func runTaskReopen(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "reopen", args, false)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "reopened %s\n", id); err != nil {
		return 1
	}
	return 0
}

func runTaskDelete(ctx context.Context, state *state, args []string) int {
	id, raw, code := taskAction(ctx, state, "delete", args, true)
	if code != 0 {
		return code
	}
	if state.Mode == modeJSON {
		if err := printRawJSON(state.Out, raw); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", id); err != nil {
		return 1
	}
	return 0
}

func runTaskQuick(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist task quick", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var note string
	var reminder string
	var autoReminder bool
	var meta bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&note, "note", "", "Note")
	fs.StringVar(&reminder, "reminder", "", "Reminder")
	fs.BoolVar(&autoReminder, "auto-reminder", false, "Auto reminder")
	fs.BoolVar(&meta, "meta", false, "Include metadata")

	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printTaskUsage(state.Out)
		return 0
	}
	text := joinArgs(fs.Args())
	if text == "" {
		writeLine(state.Err, "error: quick-add text required")
		return 2
	}
	if state.LabelCLI {
		text = ensureQuickAddLabel(text, cliLabel)
	}

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	body := map[string]any{"text": text}
	if note != "" {
		body["note"] = note
	}
	if reminder != "" {
		body["reminder"] = reminder
	}
	if autoReminder {
		body["auto_reminder"] = true
	}
	if meta {
		body["meta"] = true
	}

	task, raw, err := client.QuickAdd(ctx, body)
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

func taskAction(ctx context.Context, state *state, action string, args []string, destructive bool) (string, []byte, int) {
	fs := flag.NewFlagSet("todoist task "+action, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as task ID")
	if destructive {
		fs.BoolVar(&force, "force", false, "Skip confirmation")
	}
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 2
	}
	if help {
		printTaskUsage(state.Out)
		return "", nil, 0
	}
	if len(fs.Args()) == 0 {
		writeLine(state.Err, "error: task identifier required")
		return "", nil, 2
	}
	identifier := strings.Join(fs.Args(), " ")

	client, err := state.client()
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}

	id, err := resolveTaskID(ctx, client, identifier, forceID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}

	if destructive {
		if err := confirmDelete(state, "task", identifier, force); err != nil {
			writeLine(state.Err, "error:", err)
			return "", nil, 2
		}
	}

	var raw []byte
	switch action {
	case "close":
		raw, err = client.CloseTask(ctx, id)
	case "reopen":
		raw, err = client.ReopenTask(ctx, id)
	case "delete":
		raw, err = client.DeleteTask(ctx, id)
	default:
		return "", nil, 2
	}
	if err != nil {
		writeLine(state.Err, "error:", err)
		return "", nil, 1
	}
	return id, raw, 0
}

func resolveTask(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (todoist.Task, error) {
	if forceID {
		return client.GetTask(ctx, identifier)
	}
	return client.FindTaskByContent(ctx, identifier)
}

func resolveTaskID(ctx context.Context, client *todoist.Client, identifier string, forceID bool) (string, error) {
	if forceID {
		return identifier, nil
	}
	task, err := client.FindTaskByContent(ctx, identifier)
	if err != nil {
		return "", err
	}
	if task.ID == "" {
		return "", fmt.Errorf("task not found: %s", identifier)
	}
	return task.ID, nil
}

func resolveProjectID(ctx context.Context, client *todoist.Client, projectName, projectID string) (string, error) {
	if projectName != "" && projectID != "" {
		return "", fmt.Errorf("cannot use --project and --project-id together")
	}
	if projectID != "" {
		return projectID, nil
	}
	if projectName == "" {
		return "", nil
	}
	return client.FindProjectIDByName(ctx, projectName)
}

func mergeLabels(labels stringSlice, labelsCSV string) []string {
	combined := make([]string, 0, len(labels))
	combined = append(combined, labels...)
	if labelsCSV == "" {
		return combined
	}
	for _, label := range strings.Split(labelsCSV, ",") {
		label = strings.TrimSpace(label)
		if label == "" {
			continue
		}
		combined = appendUniqueLabel(combined, label)
	}
	return combined
}

func appendUniqueLabel(labels []string, label string) []string {
	for _, existing := range labels {
		if existing == label {
			return labels
		}
	}
	return append(labels, label)
}

func ensureQuickAddLabel(text, label string) string {
	needle := "#" + strings.ToLower(label)
	if strings.Contains(strings.ToLower(text), needle) {
		return text
	}
	return strings.TrimSpace(text) + " #" + label
}

func validateDueFlags(due, dueDate, dueDatetime string) error {
	count := 0
	if due != "" {
		count++
	}
	if dueDate != "" {
		count++
	}
	if dueDatetime != "" {
		count++
	}
	if count > 1 {
		return fmt.Errorf("use only one of --due, --due-date, or --due-datetime")
	}
	return nil
}
