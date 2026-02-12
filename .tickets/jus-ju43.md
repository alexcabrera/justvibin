---
id: jus-ju43
status: open
deps: [jus-1p4n]
links: []
created: 2026-02-12T20:00:50Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-k78p
tags: [github, deprecation]
---
# Archive srv repository

Archive the srv repository and add deprecation notice.

## Repository

alexcabrera/srv

## Steps

### 1. Update README

Replace entire README with deprecation notice:

```markdown
# ⚠️ This project has been merged into justvibin

**srv** functionality is now part of [justvibin](https://github.com/alexcabrera/justvibin).

## Migration

\`\`\`bash
# Uninstall srv (if installed via Homebrew)
brew uninstall srv

# Install justvibin
brew tap alexcabrera/tap
brew install justvibin

# Run setup (migrates existing projects)
justvibin setup
\`\`\`

## What Changed

- \`srv new\` → \`justvibin new\`
- \`srv start\` → \`justvibin start\`
- \`srv stop\` → \`justvibin stop\`
- \`srv list\` → \`justvibin list\`
- \`srv tunnel\` → \`justvibin tunnel\`
- \`srv proxy\` → \`justvibin proxy\`

justvibin includes all srv features plus:
- Template plugin system
- Framework-aware serving
- Better scaffolding

## Questions?

Open an issue on the [justvibin repository](https://github.com/alexcabrera/justvibin/issues).
\`\`\`

### 2. Update Repository Description

Set description to: "⚠️ DEPRECATED - merged into justvibin"

### 3. Archive Repository

Via GitHub UI or API:
- Settings → General → Archive this repository

### 4. Update Homebrew Tap

Remove or deprecate srv formula (separate ticket)

## Acceptance Criteria

- [ ] README replaced with deprecation notice
- [ ] Repository description updated
- [ ] Repository archived on GitHub
- [ ] Links to justvibin work

