---
id: jus-hrzu
status: closed
deps: [jus-lz4t, jus-a97f]
links: []
created: 2026-02-12T20:35:26Z
type: task
priority: 2
assignee: Alex Cabrera
parent: jus-fgp8
tags: [testing, django]
---
# Test django-hypermedia with justvibin

Test the Django template end-to-end with justvibin.

## Test Steps

```bash
# 1. Install template
justvibin install https://github.com/alexcabrera/justvibin-with-django-and-hypermedia.git
justvibin templates  # Verify listed

# 2. Create project
cd /tmp
justvibin new testdjango --template django-hypermedia
# Verify: setup.sh runs, prompts work

# 3. Start server
cd testdjango
justvibin start
# Verify: https://testdjango.localhost loads
# Verify: Homepage shows
# Verify: HTMX demo works

# 4. Test commands
justvibin stop
justvibin list  # Shows stopped
justvibin start --bg
justvibin list  # Shows running
justvibin open  # Browser opens

# 5. Test tunnel
justvibin tunnel
# Verify: Public URL generated

# 6. Cleanup
justvibin stop
justvibin remove testdjango --files
```

## Acceptance Criteria

- [ ] Template installs correctly
- [ ] Project creates with setup
- [ ] Server starts and site loads
- [ ] HTTPS works
- [ ] Stop/start cycle works
- [ ] Tunnel works
- [ ] Cleanup works

