# Changelog

## Unreleased

## 0.2.0 - 2026-01-02
- Added project commands (list/get/add/update/delete) with paging and favorites.
- Added section commands (list/get/add/update/delete) scoped by project.
- Added label commands (list/get/add/update/delete) with favorites and colors.
- Added comment commands (list/get/add/update/delete) with file attachments.
- Added upload commands for comment attachments (add/delete).
- Added MIT license.

## 0.1.0 - 2026-01-02
- Added initial Todoist CLI with task commands: list/get/add/update/close/reopen/delete/quick.
- Added auth commands (login/logout/status) and config commands (get/set/view/path).
- Added output modes (JSON/plain), quiet/verbose flags, no-input/no-color options, and optional cli label.
- Added version stamping support for `todoist --version` via Go ldflags.
