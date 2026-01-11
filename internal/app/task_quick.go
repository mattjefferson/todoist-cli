package app

import (
	"context"
	"flag"
	"io"
)

func runTaskQuick(ctx context.Context, state *state, args []string) int {
	fs := flag.NewFlagSet("todi task quick", flag.ContinueOnError)
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
