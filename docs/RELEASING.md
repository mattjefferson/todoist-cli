---
summary: 'Shared release guardrails (GitHub releases + changelog hygiene)'
read_when:
  - Preparing a release or editing release notes.
---

# Shared Release Guardrails

- Title every GitHub release as `<Project> <version>` — never the version alone.
- Release body = the curated changelog bullets for that version, verbatim and in order; no extra meta text.
- Attach all shipping artifacts (zips/tarballs/checksums/dSYMs as applicable) that the downstream clients expect.
- If the repo has its own release doc, follow it; otherwise adapt this guidance to the stack and add a repo-local checklist.
- When a release publishes, verify the tag, assets, and notes on GitHub before announcing; fix mismatches immediately (retitle, re-upload assets, or retag if necessary).

# Todoist CLI Release (this repo)

## Notes
- Tag format: `v<version>` (release workflow triggers on `v*` tags).
- GoReleaser config: `.goreleaser.yml` (project name `todoist`).
- Version stamping via ldflags (`internal/app.version`, `commit`, `date`).
- GoReleaser changelog disabled; release notes must come from `CHANGELOG.md`.

## Checklist
1. Update `CHANGELOG.md` using `docs/update-changelog.md` guidance.
2. Move relevant bullets from `## Unreleased` into a new version section (e.g., `## v1.2.3`).
3. Commit changelog.
4. Tag + push:
   - `git tag v1.2.3`
   - `git push origin v1.2.3`
5. Verify GitHub Actions `release` workflow ran.
6. Verify GitHub release:
   - Title: `todoist v1.2.3` (matches GoReleaser `name_template`).
   - Body: `CHANGELOG.md` bullets for that version, verbatim + ordered.
   - Assets: OS/arch tarballs + `checksums.txt` (archives include LICENSE/README/CHANGELOG).
