package app

import (
	"context"
	"flag"
	"io"
	"strconv"
	"strings"

	"github.com/mattjefferson/todoist-cli/internal/todoist"
)

func runActivity(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printActivityUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runActivityList(ctx, state, args[1:])
	case "-h", "--help", "help":
		printActivityUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown activity command:", args[0])
		printActivityUsage(state.Err)
		return 2
	}
}

func runActivityList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todoist activity list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var limit int
	var cursor string
	var objectType string
	var objectID string
	var parentProjectID string
	var parentItemID string
	var includeParent bool
	var includeChildren bool
	var initiatorID string
	var initiatorIDNull bool
	var eventType string
	var objectEventTypes string
	var annotateNotes bool
	var annotateParents bool
	var all bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.IntVar(&limit, "limit", 30, "Max events per page (1-100)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.StringVar(&objectType, "object-type", "", "Object type filter")
	fs.StringVar(&objectID, "object-id", "", "Object ID filter")
	fs.StringVar(&parentProjectID, "parent-project-id", "", "Parent project ID filter")
	fs.StringVar(&parentItemID, "parent-item-id", "", "Parent item ID filter")
	fs.BoolVar(&includeParent, "include-parent-object", false, "Include parent object data")
	fs.BoolVar(&includeChildren, "include-child-objects", false, "Include child objects data")
	fs.StringVar(&initiatorID, "initiator-id", "", "Initiator user ID")
	fs.BoolVar(&initiatorIDNull, "initiator-id-null", false, "Only events without an initiator")
	fs.StringVar(&eventType, "event-type", "", "Event type filter")
	fs.StringVar(&objectEventTypes, "object-event-types", "", "Object event types (comma-separated)")
	fs.BoolVar(&annotateNotes, "annotate-notes", false, "Include note info in extra_data")
	fs.BoolVar(&annotateParents, "annotate-parents", false, "Include parent info in extra_data")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printActivityUsage(state.Out)
		return 0
	}
	if len(fs.Args()) > 0 {
		writeLine(state.Err, "error: unexpected arguments")
		return 2
	}
	if initiatorID != "" && initiatorIDNull {
		writeLine(state.Err, "error: cannot use --initiator-id and --initiator-id-null together")
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
	if objectType != "" {
		params["object_type"] = objectType
	}
	if objectID != "" {
		params["object_id"] = objectID
	}
	if parentProjectID != "" {
		params["parent_project_id"] = parentProjectID
	}
	if parentItemID != "" {
		params["parent_item_id"] = parentItemID
	}
	if includeParent {
		params["include_parent_object"] = "true"
	}
	if includeChildren {
		params["include_child_objects"] = "true"
	}
	if initiatorID != "" {
		params["initiator_id"] = initiatorID
	}
	if initiatorIDNull {
		params["initiator_id_null"] = "true"
	}
	if eventType != "" {
		params["event_type"] = eventType
	}
	if objectEventTypes != "" {
		normalized := normalizeCSV(objectEventTypes)
		if normalized != "" {
			params["object_event_types"] = normalized
		}
	}
	if annotateNotes {
		params["annotate_notes"] = "true"
	}
	if annotateParents {
		params["annotate_parents"] = "true"
	}

	if all {
		activities, err := client.ListActivitiesAll(ctx, params)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printActivities(state.Out, activities, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	activities, next, err := client.ListActivities(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": activities, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printActivities(state.Out, activities, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func normalizeCSV(input string) string {
	if input == "" {
		return ""
	}
	parts := strings.Split(input, ",")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		cleaned = append(cleaned, part)
	}
	return strings.Join(cleaned, ",")
}
