# todoist-cli
A cli for interacting with todoist

## Usage
```text
todoist - Todoist CLI

USAGE:
  todoist [global flags] <command> [args]

COMMANDS:
  task    Manage tasks
  project Manage projects
  comment Manage comments
  label   Manage labels
  upload  Manage uploads
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
```
