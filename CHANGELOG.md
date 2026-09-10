# Changelog

## 0.2.0 — interactive configuration

- New `config` subcommand: a full-screen editor (`↑`/`↓` to move, `space` to toggle, `s` to save, `q` to quit) that switches every segment on or off and picks the theme and icon set, with a live preview rendered from the current directory. Standard library only — raw mode comes from `stty`.
- Settings persist in `~/.claude/claude-statusline.json` (`{ "theme", "icon_font", "segments" }`), overridable per run by `CLAUDE_STATUSLINE_THEME`, `CLAUDE_STATUSLINE_ICON_FONT`, `CLAUDE_STATUSLINE_SEGMENTS`; the file's location is overridable by `CLAUDE_STATUSLINE_CONFIG`. Priority: environment → file → built-in default.
- Missing file, missing key, or a segment absent from `segments` all fall back to defaults; a malformed file warns on stderr and renders defaults. Saves go through a temp file + rename.
- New `help` subcommand: usage, the settings path, and every registered segment, theme and icon set — replacing the `--list-*` flags.
- **Breaking:** the `--theme`, `--segments`, `--icon-font` and `--list-*` flags are gone. A leftover flag in `settings.json` warns on stderr and still renders, so an old command line degrades instead of breaking the statusline.
- Icon sets, previously unreleased: `plain` (default, universal Unicode), `nerd` (Nerd Font PUA glyphs), `none` (text only), plus `icons.Prefix` so partial sets render without dangling separators.
- `meta` recognises the Fable model family.
- Segments carry an explicit `Order`, so rows always render git → path → meta; a segment registered without one is appended last.
- `statusline.Input` exposes named `Workspace`, `Model` and `ContextWindow` types, so callers (the config preview) can build a payload.

## 0.1.0 — initial release

- Native Go statusline for Claude Code, replacing the JS/bun script.
- Theme registry with two bundled palettes: `atom-one-dark` (default, byte-identical to the previous statusline) and `dracula`.
- Segment registry with three bundled segments: `git` (branch + worktree mark + ahead/behind), `path` (cwd with `$HOME → ~`), `meta` (model family + version + context-usage percentage chip).
- CLI flags: `--theme`, `--segments`, `--list-themes`, `--list-segments`.
- Environment variables: `CLAUDE_STATUSLINE_THEME`, `CLAUDE_STATUSLINE_SEGMENTS`.
- Unknown theme/segment values produce stderr warnings and graceful fallback; nothing ever leaks onto stdout.
- Makefile targets: `build`, `install`, `clean`, `test`, `fmt`, `vet`.
