# Changelog

## 0.1.0 — initial release

- Native Go statusline for Claude Code, replacing the JS/bun script.
- Theme registry with two bundled palettes: `atom-one-dark` (default, byte-identical to the previous statusline) and `dracula`.
- Segment registry with three bundled segments: `git` (branch + worktree mark + ahead/behind), `path` (cwd with `$HOME → ~`), `meta` (model family + version + context-usage percentage chip).
- CLI flags: `--theme`, `--segments`, `--list-themes`, `--list-segments`.
- Environment variables: `CLAUDE_STATUSLINE_THEME`, `CLAUDE_STATUSLINE_SEGMENTS`.
- Unknown theme/segment values produce stderr warnings and graceful fallback; nothing ever leaks onto stdout.
- Makefile targets: `build`, `install`, `clean`, `test`, `fmt`, `vet`.
