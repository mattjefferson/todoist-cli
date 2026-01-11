package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattjefferson/todi/internal/todi"
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

func resolveTask(ctx context.Context, client *todi.Client, identifier string, forceID bool) (todi.Task, error) {
	if forceID {
		return client.GetTask(ctx, identifier)
	}
	return client.FindTaskByContent(ctx, identifier)
}

func resolveTaskID(ctx context.Context, client *todi.Client, identifier string, forceID bool) (string, error) {
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

func resolveProjectID(ctx context.Context, client *todi.Client, projectName, projectID string) (string, error) {
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
