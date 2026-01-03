package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mattjefferson/todi/internal/todi"
)

func runSection(ctx context.Context, state *state, args []string) int {
	if len(args) == 0 {
		printSectionUsage(state.Out)
		return 2
	}
	switch args[0] {
	case "list":
		return runSectionList(ctx, state, args[1:])
	case "get":
		return runSectionGet(ctx, state, args[1:])
	case "add":
		return runSectionAdd(ctx, state, args[1:])
	case "update":
		return runSectionUpdate(ctx, state, args[1:])
	case "delete":
		return runSectionDelete(ctx, state, args[1:])
	case "-h", "--help", "help":
		printSectionUsage(state.Out)
		return 0
	default:
		writeLine(state.Err, "error: unknown section command:", args[0])
		printSectionUsage(state.Err)
		return 2
	}
}

func runSectionList(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi section list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var projectName string
	var projectID string
	var limit int
	var cursor string
	var all bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.IntVar(&limit, "limit", 50, "Max sections per page (1-200)")
	fs.StringVar(&cursor, "cursor", "", "Pagination cursor")
	fs.BoolVar(&all, "all", false, "Fetch all pages")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printSectionUsage(state.Out)
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

	projectIDValue, err := resolveProjectID(ctx, client, projectName, projectID)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	params := map[string]string{}
	if projectIDValue != "" {
		params["project_id"] = projectIDValue
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if cursor != "" {
		params["cursor"] = cursor
	}

	if all {
		sections, err := client.ListSectionsAll(ctx, params)
		if err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		if err := printSections(state.Out, sections, state.Mode); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}

	sections, next, err := client.ListSections(ctx, params)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if state.Mode == modeJSON {
		payload := map[string]any{"results": sections, "next_cursor": next}
		if err := printJSON(state.Out, payload); err != nil {
			writeLine(state.Err, "error:", err)
			return 1
		}
		return 0
	}
	if err := printSections(state.Out, sections, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runSectionGet(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi section get", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var projectName string
	var projectID string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as section ID")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printSectionUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: section identifier required")
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

	section, err := resolveSection(ctx, client, identifier, forceID, projectIDValue)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := printSection(state.Out, section, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runSectionAdd(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi section add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var projectName string
	var projectID string
	var orderValue string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.StringVar(&orderValue, "order", "", "Section order")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printSectionUsage(state.Out)
		return 0
	}
	name := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if name == "" {
		writeLine(state.Err, "error: section name required")
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
	if projectIDValue == "" {
		writeLine(state.Err, "error: project required")
		return 2
	}

	body := map[string]any{"name": name, "project_id": projectIDValue}
	if orderValue != "" {
		order, err := strconv.Atoi(orderValue)
		if err != nil {
			writeLine(state.Err, "error: order must be an integer")
			return 2
		}
		body["order"] = order
	}

	section, raw, err := client.CreateSection(ctx, body)
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
	if err := printSection(state.Out, section, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runSectionUpdate(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi section update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var projectName string
	var projectID string
	var name string
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as section ID")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.StringVar(&name, "name", "", "Section name")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printSectionUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: section identifier required")
		return 2
	}
	if name == "" {
		writeLine(state.Err, "error: name required")
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

	sectionID, err := resolveSectionIDFromIdentifier(ctx, client, identifier, forceID, projectIDValue)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}

	body := map[string]any{"name": name}
	section, raw, err := client.UpdateSection(ctx, sectionID, body)
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
	if section.ID == "" {
		writeLine(state.Out, "ok")
		return 0
	}
	if err := printSection(state.Out, section, state.Mode); err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	return 0
}

func runSectionDelete(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi section delete", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var help bool
	var forceID bool
	var projectName string
	var projectID string
	var force bool
	fs.BoolVar(&help, "help", false, "Show help")
	fs.BoolVar(&help, "h", false, "Show help")
	fs.BoolVar(&forceID, "id", false, "Treat argument as section ID")
	fs.StringVar(&projectName, "project", "", "Project title (exact match)")
	fs.StringVar(&projectID, "project-id", "", "Project ID")
	fs.BoolVar(&force, "force", false, "Skip confirmation")
	if err := fs.Parse(args); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}
	if help {
		printSectionUsage(state.Out)
		return 0
	}
	identifier := strings.TrimSpace(strings.Join(fs.Args(), " "))
	if identifier == "" {
		writeLine(state.Err, "error: section identifier required")
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

	sectionID, err := resolveSectionIDFromIdentifier(ctx, client, identifier, forceID, projectIDValue)
	if err != nil {
		writeLine(state.Err, "error:", err)
		return 1
	}
	if err := confirmDelete(state, "section", identifier, force); err != nil {
		writeLine(state.Err, "error:", err)
		return 2
	}

	raw, err := client.DeleteSection(ctx, sectionID)
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
	if _, err := fmt.Fprintf(state.Out, "deleted %s\n", sectionID); err != nil {
		return 1
	}
	return 0
}

func resolveSection(ctx context.Context, client *todi.Client, identifier string, forceID bool, projectID string) (todi.Section, error) {
	if forceID {
		return client.GetSection(ctx, identifier)
	}
	return client.FindSectionByName(ctx, identifier, projectID)
}

func resolveSectionIDFromIdentifier(ctx context.Context, client *todi.Client, identifier string, forceID bool, projectID string) (string, error) {
	if forceID {
		return identifier, nil
	}
	return client.FindSectionIDByName(ctx, identifier, projectID)
}
