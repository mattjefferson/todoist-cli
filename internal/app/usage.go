package app

import (
	"fmt"
	"io"
)

func printUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist - Todoist CLI

USAGE:
  todoist [global flags] <command> [args]

COMMANDS:
  task    Manage tasks
  project Manage projects
  comment Manage comments
  label   Manage labels
  auth    Manage auth token
  config  Manage config

GLOBAL FLAGS:
  -h, --help        Show help
  --version         Show version
  -q, --quiet        Less output
  -v, --verbose      Verbose output
  --json            JSON output
  --plain           Plain output
  --no-input        Disable prompts
  --no-color        Disable color
  --config <path>   Config path override
  --api-base <url>  API base (default https://api.todoist.com)
  --label-cli       Add label 'cli' to created tasks

NOTES:
  Task identifiers accept exact task titles unless --id is set.
  Project identifiers accept exact project titles unless --id is set.
  Label identifiers accept exact label names unless --id is set.
`); err != nil {
		return
	}
}

func printTaskUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist task - task commands

USAGE:
  todoist task list [project_title]
  todoist task get <task>
  todoist task add <content>
  todoist task update <task>
  todoist task close <task>
  todoist task reopen <task>
  todoist task delete <task>
  todoist task quick <text>

NOTES:
  <task> accepts exact task title unless --id is set.
`); err != nil {
		return
	}
}

func printAuthUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist auth - auth commands

USAGE:
  todoist auth login
  todoist auth logout
  todoist auth status
`); err != nil {
		return
	}
}

func printCommentUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist comment - comment commands

USAGE:
  todoist comment list --task <title>
  todoist comment list --task-id <id>
  todoist comment list --project <title>
  todoist comment list --project-id <id>
  todoist comment get <comment_id>
  todoist comment add <content>
  todoist comment update <comment_id> --content <text>
  todoist comment delete <comment_id>
`); err != nil {
		return
	}
}

func printLabelUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist label - label commands

USAGE:
  todoist label list
  todoist label get <label>
  todoist label add <name>
  todoist label update <label>
  todoist label delete <label>

NOTES:
  <label> accepts exact label name unless --id is set.
`); err != nil {
		return
	}
}

func printProjectUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist project - project commands

USAGE:
  todoist project list
  todoist project get <project>
  todoist project add <name>
  todoist project update <project>
  todoist project delete <project>

NOTES:
  <project> accepts exact project title unless --id is set.
`); err != nil {
		return
	}
}

func printConfigUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist config - config commands

USAGE:
  todoist config get <key>
  todoist config set <key> <value>
  todoist config path
  todoist config view

KEYS:
  token
  api_base
  default_project
  default_labels
  label_cli
`); err != nil {
		return
	}
}
