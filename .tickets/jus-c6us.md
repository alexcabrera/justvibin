---
id: jus-c6us
status: open
deps: [jus-f3wf]
links: []
created: 2026-02-12T20:35:12Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-zf02
tags: [templates, documentation]
---
# Update hypertext template AGENTS.md

Update the hypertext template's AGENTS.md to document justvibin commands.

## Changes

Update Development section:

```markdown
## Development

### With justvibin (recommended)

\`\`\`bash
justvibin start    # Start server
justvibin open     # Open in browser
justvibin stop     # Stop server
justvibin tunnel   # Share publicly
\`\`\`

Your site is at https://<project-name>.localhost

### Without justvibin

\`\`\`bash
python3 -m http.server 8000
# Visit http://localhost:8000
\`\`\`
```

## Acceptance Criteria

- [ ] justvibin commands documented
- [ ] Standalone fallback documented
- [ ] Old serve.sh references removed

