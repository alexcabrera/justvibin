---
id: jus-3ibt
status: closed
deps: [jus-w9pp, jus-njz0]
links: []
created: 2026-02-17T15:13:46Z
type: task
priority: 3
assignee: Alex Cabrera
parent: jus-fxaf
tags: [cli, fang, cobra, completions]
---
# Add shell completion generation

Verify and document shell completion generation provided by Fang/Cobra.

## Background

justvibin is a CLI for scaffolding web projects. Cobra provides built-in shell completion generation, and Fang adds a 'completion' command automatically.

## Current State

Fang automatically adds:
- justvibin completion bash
- justvibin completion zsh
- justvibin completion fish
- justvibin completion powershell

These should work out of the box after Fang integration.

## Implementation

1. Verify completion command exists:
   - Run justvibin completion --help
   - Run justvibin completion bash
   - Run justvibin completion zsh
   - Run justvibin completion fish

2. Test completions actually work:
   - Source bash completion and test
   - Source zsh completion and test
   - Verify flag and subcommand completion

3. Add documentation to help text or README:
   - Installation instructions for each shell
   - Where to put completion scripts

4. Optional: Customize completion descriptions if needed

## Completion Installation Instructions (for docs)

### Bash
```bash
# Add to ~/.bashrc
source <(justvibin completion bash)

# Or save to file
justvibin completion bash > /usr/local/etc/bash_completion.d/justvibin
```

### Zsh
```zsh
# Add to ~/.zshrc (before compinit)
source <(justvibin completion zsh)

# Or save to fpath
justvibin completion zsh > "${fpath[1]}/_justvibin"
```

### Fish
```fish
justvibin completion fish | source

# Or save permanently
justvibin completion fish > ~/.config/fish/completions/justvibin.fish
```

## Acceptance Criteria

- [ ] justvibin completion bash generates valid script
- [ ] justvibin completion zsh generates valid script
- [ ] justvibin completion fish generates valid script
- [ ] Completions work for commands (new, templates, install, proxy, setup)
- [ ] Completions work for flags (--template, --json, etc.)
- [ ] Documentation added to help or README

