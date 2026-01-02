package app

import (
	"fmt"
	"io"
)

func printUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist - Todoist CLI

USAGE:
  todoist [global flags] <task-command> [args]
  todoist [global flags] <command> [args]

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
  todoist auth login            Save token to config
  TODOIST_TOKEN                 Overrides token in config
  todoist config path           Print config file path

NOTES:
  Task commands can be used without the "task" prefix (e.g. todoist list, todoist add).
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
	if _, err := fmt.Fprint(out, `todoist task - task commands

USAGE:
  todoist list [project_title]
  todoist get <task>
  todoist add <content>
  todoist update <task>
  todoist close <task>
  todoist reopen <task>
  todoist delete <task>
  todoist quick <text>

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
  todoist list
  todoist list "Inbox" --all
  todoist add "Write docs" --project "Docs"
  todoist update "Write docs" --content "Write help" --priority 2
  todoist close "Write docs"

NOTES:
  Task commands can also be called with the "task" prefix.
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

NOTES:
  login prompts for a token (TTY required).
  status reports token source (TODOIST_TOKEN or config).
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
  todoist comment list --task "Write docs"
  todoist comment add "LGTM" --task-id 123 --notify 456
  todoist comment add "See file" --task "Inbox" --file ./spec.pdf
`); err != nil {
		return
	}
}

func printActivityUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist activity - activity log commands

USAGE:
  todoist activity list

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
  todoist activity list
  todoist activity list --object-type item --event-type completed
  todoist activity list --object-id 123 --include-parent-object
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
  todoist label list
  todoist label add "waiting" --color blue --favorite
  todoist label update "waiting" --color red

NOTES:
  <label> accepts exact label name unless --id is set.
`); err != nil {
		return
	}
}

func printUploadUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist upload - upload commands

USAGE:
  todoist upload add <path>
  todoist upload delete <file_url>

FLAGS (add):
  --project <title>        Project title (exact match)
  --project-id <id>        Project ID
  --name <name>            Override file name

FLAGS (delete):
  --file-url <url>         File URL
  --force                  Skip confirmation

EXAMPLES:
  todoist upload add ./spec.pdf --project "Docs"
  todoist upload delete https://.../file.pdf
`); err != nil {
		return
	}
}

func printSectionUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist section - section commands

USAGE:
  todoist section list
  todoist section get <section>
  todoist section add <name>
  todoist section update <section> --name <name>
  todoist section delete <section>

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
  todoist section list --project "Docs"
  todoist section add "Backlog" --project "Docs"
  todoist section update "Backlog" --project "Docs" --name "Next"

NOTES:
  <section> accepts exact section name unless --id is set.
  Use --project or --project-id to scope name lookups.
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

FLAGS (list):
  --limit <n>              Max projects per page (1-200)
  --cursor <cursor>        Pagination cursor
  --all                    Fetch all pages

FLAGS (get/update/delete):
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
  todoist project list --all
  todoist project add "Docs" --favorite
  todoist project update "Docs" --view board

NOTES:
  <project> accepts exact project title unless --id is set.
`); err != nil {
		return
	}
}

func printUserUsage(out io.Writer) {
	if _, err := fmt.Fprint(out, `todoist user - user commands

USAGE:
  todoist user info

OUTPUT:
  id, email, full_name
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
