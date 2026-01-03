package app

import (
	"fmt"
	"io"
)

func printUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi - CLI

USAGE:
  todi [global flags] <task-command> [args]
  todi [global flags] <command> [args]

COMMANDS:
  task    Manage tasks
  project Manage projects
  comment Manage comments
  activity Manage activity log
  label   Manage labels
  upload  Manage uploads
  section Manage sections
  user    Manage user info
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

OUTPUT MODES:
  default           Human-friendly tables
  --plain           Tab-delimited output for scripts
  --json            Structured JSON output

AUTH:
  todi auth login            Save token to config
  TODOIST_TOKEN                 Overrides token in config
  todi config path           Print config file path

NOTES:
  Task commands can be used without the "task" prefix (e.g. todi list, todi add).
  Task identifiers accept exact task titles unless --id is set.
  Project identifiers accept exact project titles unless --id is set.
  Label identifiers accept exact label names unless --id is set.
  Section identifiers accept exact section names unless --id is set.
  Destructive commands require TTY confirmation or --force.
`); err != nil {
		return
	}
}

func printTaskUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi task - task commands

USAGE:
  todi list [project_title]
  todi get <task>
  todi add <content>
  todi update <task>
  todi close <task>
  todi reopen <task>
  todi delete <task>
  todi quick <text>

FLAGS (list):
  --project <title>        Project title (exact match)
  --label <name>           Filter by label
  --limit <n>              Max tasks per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (get/close/reopen/delete):
  --id                     Treat argument as task ID

FLAGS (delete):
  --force                  Skip confirmation

FLAGS (add):
  --description <text>     Task description
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --label <name>           Label (repeatable)
  --labels <a,b>           Labels (comma-separated)
  --priority <1-4>         Task priority
  --assignee <id>          Assignee ID
  --due <text>             Due string
  --due-date <YYYY-MM-DD>  Due date
  --due-datetime <RFC3339> Due datetime
  --due-lang <code>        Due language
  --duration <n>           Duration value
  --duration-unit <unit>   Duration unit (minute|day)
  --deadline-date <date>   Deadline date (YYYY-MM-DD)

FLAGS (update):
  --id                     Treat argument as task ID
  --content <text>         New content
  --description <text>     New description
  --label <name>           Label (repeatable)
  --labels <a,b>           Labels (comma-separated)
  --priority <1-4>         Task priority
  --assignee <id>          Assignee ID
  --due <text>             Due string
  --due-date <YYYY-MM-DD>  Due date
  --due-datetime <RFC3339> Due datetime
  --due-lang <code>        Due language
  --duration <n>           Duration value
  --duration-unit <unit>   Duration unit (minute|day)
  --deadline-date <date>   Deadline date (YYYY-MM-DD)

FLAGS (quick):
  --note <text>            Add note
  --reminder <text>        Reminder
  --auto-reminder          Auto reminder
  --meta                   Include metadata

EXAMPLES:
  todi list
  todi list "Inbox" --all
  todi add "Write docs" --project "Docs"
  todi update "Write docs" --content "Write help" --priority 2
  todi close "Write docs"

NOTES:
  Task commands can also be called with the "task" prefix.
  <task> accepts exact task title unless --id is set.
`); err != nil {
		return
	}
}

func printAuthUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi auth - auth commands

USAGE:
  todi auth login
  todi auth logout
  todi auth status

NOTES:
  login prompts for a token (TTY required).
  status reports token source (TODOIST_TOKEN or config).
`); err != nil {
		return
	}
}

func printCommentUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi comment - comment commands

USAGE:
  todi comment list --task <title>
  todi comment list --task-id <id>
  todi comment list --project <title>
  todi comment list --project-id <id>
  todi comment get <comment_id>
  todi comment add <content>
  todi comment update <comment_id> --content <text>
  todi comment delete <comment_id>

FLAGS (list):
  --task <title>           Task title (exact match)
  --task-id <id>           Task ID
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --limit <n>              Max comments per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (add):
  --task <title>           Task title (exact match)
  --task-id <id>           Task ID
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --notify <uid>           UID to notify (repeatable)
  --file <path>            Upload file attachment
  --file-name <name>       Override upload file name

FLAGS (update):
  --content <text>         New content

FLAGS (delete):
  --force                  Skip confirmation

EXAMPLES:
  todi comment list --task "Write docs"
  todi comment add "LGTM" --task-id 123 --notify 456
  todi comment add "See file" --task "Inbox" --file ./spec.pdf
`); err != nil {
		return
	}
}

func printActivityUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi activity - activity log commands

USAGE:
  todi activity list

FLAGS (list):
  --limit <n>              Max events per page (1-100)
  --cursor <cursor>        Pagination cursor
  --object-type <type>     Object type filter
  --object-id <id>         Object ID filter
  --parent-project-id <id> Parent project ID filter
  --parent-item-id <id>    Parent item ID filter
  --include-parent-object  Include parent object data
  --include-child-objects  Include child objects data
  --initiator-id <id>      Initiator user ID
  --initiator-id-null      Only events without an initiator
  --event-type <type>      Event type filter
  --object-event-types <a,b>  Object event types (comma-separated)
  --annotate-notes         Include note info in extra_data
  --annotate-parents       Include parent info in extra_data
  --all                    Fetch all pages

EXAMPLES:
  todi activity list
  todi activity list --object-type item --event-type completed
  todi activity list --object-id 123 --include-parent-object
`); err != nil {
		return
	}
}

func printLabelUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi label - label commands

USAGE:
  todi label list
  todi label get <label>
  todi label add <name>
  todi label update <label>
  todi label delete <label>

FLAGS (list):
  --limit <n>              Max labels per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (get/update/delete):
  --id                     Treat argument as label ID

FLAGS (add):
  --color <name>           Label color
  --favorite               Mark as favorite

FLAGS (update):
  --name <name>            New name
  --color <name>           New color
  --favorite               Mark as favorite
  --unfavorite             Remove favorite

FLAGS (delete):
  --force                  Skip confirmation

EXAMPLES:
  todi label list
  todi label add "waiting" --color blue --favorite
  todi label update "waiting" --color red

NOTES:
  <label> accepts exact label name unless --id is set.
`); err != nil {
		return
	}
}

func printUploadUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi upload - upload commands

USAGE:
  todi upload add <path>
  todi upload delete <file_url>

FLAGS (add):
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --name <name>            Override file name

FLAGS (delete):
  --file-url <url>         File URL
  --force                  Skip confirmation

EXAMPLES:
  todi upload add ./spec.pdf --project "Docs"
  todi upload delete https://.../file.pdf
`); err != nil {
		return
	}
}

func printSectionUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi section - section commands

USAGE:
  todi section list
  todi section get <section>
  todi section add <name>
  todi section update <section> --name <name>
  todi section delete <section>

FLAGS (list):
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --limit <n>              Max sections per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (get/update/delete):
  --id                     Treat argument as section ID
  --project <title>        Project title (exact match) for name lookup
  --project-id <id>        Project ID for name lookup

FLAGS (add):
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --order <n>              Section order

FLAGS (delete):
  --force                  Skip confirmation

EXAMPLES:
  todi section list --project "Docs"
  todi section add "Backlog" --project "Docs"
  todi section update "Backlog" --project "Docs" --name "Next"

NOTES:
  <section> accepts exact section name unless --id is set.
  Use --project or --project-id to scope name lookups.
`); err != nil {
		return
	}
}

func printProjectUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi project - project commands

USAGE:
  todi project list
  todi project get <project>
  todi project add <name>
  todi project update <project>
  todi project archive <project>
  todi project unarchive <project>
  todi project delete <project>

FLAGS (list):
  --limit <n>              Max projects per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (get/update/archive/unarchive/delete):
  --id                     Treat argument as project ID

FLAGS (add):
  --parent <title>         Parent project title (exact match)
  --parent-id <id>         Parent project ID
  --color <name>           Project color
  --favorite               Mark as favorite
  --view <style>           View style (list|board)

FLAGS (update):
  --name <name>            New name
  --color <name>           New color
  --favorite               Mark as favorite
  --unfavorite             Remove favorite
  --view <style>           View style (list|board)

FLAGS (delete):
  --force                  Skip confirmation

EXAMPLES:
  todi project list --all
  todi project add "Docs" --favorite
  todi project update "Docs" --view board
  todi project archive "Docs"

NOTES:
  <project> accepts exact project title unless --id is set.
`); err != nil {
		return
	}
}

func printUserUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi user - user commands

USAGE:
  todi user info

OUTPUT:
  id, email, full_name
`); err != nil {
		return
	}
}

func printConfigUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todi config - config commands

USAGE:
  todi config get <key>
  todi config set <key> <value>
  todi config path
  todi config view

KEYS:
  token              Stored auth token (set via auth login)
  api_base           API base URL override
  default_project    Stored default project name
  default_labels     Stored default labels (comma-separated)
  label_cli          Add label 'cli' to created tasks

NOTES:
  token cannot be set via config set.
`); err != nil {
		return
	}
}
