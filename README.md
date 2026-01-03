# todi
CLI for Todoist.

## Quick start
- `todi auth login`
- `todi auth status`
- `todi list`

## Auth
- `todi auth login` prompts for token (TTY required).
- `TODOIST_TOKEN` overrides token stored in config.
- `todi auth logout` clears config token.

## Output modes
- Default: human-readable tables.
- `--plain`: tab-delimited output (stable for scripts).
- `--json`: structured JSON output.
- Errors go to stderr.

## Global flags
- `-h, --help` show help
- `--version` show version
- `-q, --quiet` less output
- `-v, --verbose` verbose output (logs HTTP requests)
- `--json` JSON output
- `--plain` plain output
- `--no-input` disable prompts
- `--no-color` disable color
- `--config <path>` config path override
- `--api-base <url>` API base override
- `--label-cli` add label `cli` to created tasks

## Name resolution rules
- Task, project, label, section identifiers are exact name matches unless `--id` is set.
- Section name lookups should be scoped with `--project` or `--project-id`.
- Comment list/add requires a task or project scope.

## Commands

### user
Fetch current user profile fields.

Subcommands:
- `info`
  - Output: `id`, `email`, `full_name`

Example:
- `todi user info`

### task
Manage tasks.

Task commands can be called without the `task` prefix (for example, `todi list`).

Subcommands:
- `list [project_title]`
  - Flags: `--project`, `--label`, `--limit`, `--cursor`, `--all`
- `get <task>`
  - Flags: `--id`
- `add <content>`
  - Flags: `--description`, `--project`, `--project-id`, `--label` (repeatable), `--labels`,
    `--priority`, `--assignee`, `--due`, `--due-date`, `--due-datetime`, `--due-lang`,
    `--duration`, `--duration-unit`, `--deadline-date`
- `update <task>`
  - Flags: `--id`, `--content`, `--description`, `--label` (repeatable), `--labels`, `--priority`,
    `--assignee`, `--due`, `--due-date`, `--due-datetime`, `--due-lang`, `--duration`,
    `--duration-unit`, `--deadline-date`
- `close <task>`
  - Flags: `--id`
- `reopen <task>`
  - Flags: `--id`
- `delete <task>`
  - Flags: `--id`, `--force`
- `quick <text>`
  - Flags: `--note`, `--reminder`, `--auto-reminder`, `--meta`

Examples:
- `todi list`
- `todi list "Inbox" --all`
- `todi add "Write docs" --project "Docs" --priority 2`
- `todi update "Write docs" --content "Write help"`
- `todi delete "Write docs" --force`

### project
Manage projects.

Subcommands:
- `list`
  - Flags: `--limit`, `--cursor`, `--all`
- `get <project>`
  - Flags: `--id`
- `add <name>`
  - Flags: `--parent`, `--parent-id`, `--color`, `--favorite`, `--view`
- `update <project>`
  - Flags: `--id`, `--name`, `--color`, `--favorite`, `--unfavorite`, `--view`
- `delete <project>`
  - Flags: `--id`, `--force`

Examples:
- `todi project list --all`
- `todi project add "Docs" --favorite`
- `todi project update "Docs" --view board`

### section
Manage sections.

Subcommands:
- `list`
  - Flags: `--project`, `--project-id`, `--limit`, `--cursor`, `--all`
- `get <section>`
  - Flags: `--id`, `--project`, `--project-id`
- `add <name>`
  - Flags: `--project`, `--project-id`, `--order`
- `update <section>`
  - Flags: `--id`, `--project`, `--project-id`, `--name`
- `delete <section>`
  - Flags: `--id`, `--project`, `--project-id`, `--force`

Examples:
- `todi section list --project "Docs"`
- `todi section add "Backlog" --project "Docs"`
- `todi section update "Backlog" --project "Docs" --name "Next"`

### label
Manage labels.

Subcommands:
- `list`
  - Flags: `--limit`, `--cursor`, `--all`
- `get <label>`
  - Flags: `--id`
- `add <name>`
  - Flags: `--color`, `--favorite`
- `update <label>`
  - Flags: `--id`, `--name`, `--color`, `--favorite`, `--unfavorite`
- `delete <label>`
  - Flags: `--id`, `--force`

Examples:
- `todi label list`
- `todi label add "waiting" --color blue --favorite`
- `todi label update "waiting" --color red`

### comment
Manage comments.

Subcommands:
- `list`
  - Flags: `--task`, `--task-id`, `--project`, `--project-id`, `--limit`, `--cursor`, `--all`
- `get <comment_id>`
- `add <content>`
  - Flags: `--task`, `--task-id`, `--project`, `--project-id`, `--notify` (repeatable),
    `--file`, `--file-name`
- `update <comment_id>`
  - Flags: `--content`
- `delete <comment_id>`
  - Flags: `--force`

Examples:
- `todi comment list --task "Write docs"`
- `todi comment add "LGTM" --task-id 123 --notify 456`
- `todi comment add "See file" --task "Inbox" --file ./spec.pdf`

### activity
Fetch activity logs.

Subcommands:
- `list`
  - Flags: `--limit`, `--cursor`, `--object-type`, `--object-id`, `--parent-project-id`,
    `--parent-item-id`, `--include-parent-object`, `--include-child-objects`,
    `--initiator-id`, `--initiator-id-null`, `--event-type`, `--object-event-types`,
    `--annotate-notes`, `--annotate-parents`, `--all`

Examples:
- `todi activity list`
- `todi activity list --object-type item --event-type completed`
- `todi activity list --object-id 123 --include-parent-object`

### upload
Manage uploads for comment attachments.

Subcommands:
- `add <path>`
  - Flags: `--project`, `--project-id`, `--name`
- `delete <file_url>`
  - Flags: `--file-url`, `--force`

Examples:
- `todi upload add ./spec.pdf --project "Docs"`
- `todi upload delete https://.../file.pdf`

### auth
Manage auth token.

Subcommands:
- `login` (TTY only)
- `logout`
- `status`

### config
Manage local config.

Subcommands:
- `get <key>`
- `set <key> <value>`
- `path`
- `view`

Keys:
- `token` (set via `auth login`)
- `api_base`
- `default_project`
- `default_labels`
- `label_cli`

Notes:
- Use `todi config path` to find the config file.
- `default_project` and `default_labels` are stored but not yet applied automatically.
