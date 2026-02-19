---
id: jus-a97f
status: closed
deps: [jus-bx3x]
links: []
created: 2026-02-12T20:34:47Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-swo4
tags: [templates, documentation]
---
# Update django template AGENTS.md for justvibin

Update the Django template's AGENTS.md to document justvibin commands.

## Changes

Add section for justvibin development workflow:

```markdown
## Development with justvibin

If using justvibin CLI:

\`\`\`bash
# Start development server
justvibin start

# Open in browser
justvibin open

# Stop server
justvibin stop

# Share publicly
justvibin tunnel
\`\`\`

Your project is accessible at https://<project-name>.localhost

## Standalone Development

Without justvibin:

\`\`\`bash
./start.sh --dev
\`\`\`
```

## Also Update

- Development Commands section
- Any references to ./start.sh --dev

## Acceptance Criteria

- [ ] justvibin commands documented
- [ ] Standalone commands still documented
- [ ] URL format explained

