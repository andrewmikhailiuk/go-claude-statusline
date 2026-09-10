# claude-statusline

A native Go statusline for [Claude Code](https://claude.com/claude-code) — zero dependencies, pluggable themes, pluggable segments, configured from an interactive terminal UI.

## Design goals

- **No runtime dependency.** A single self-contained binary; no node, no bun, no python. Standard library only — even the interactive config UI, which drives the terminal through `stty`.
- **Multiple themes.** Themes live in their own files and self-register on import. Adding one is a single-file change.
- **Pluggable segments.** Each row of the statusline (git, path, model+context, …) is a `Segment` registered through `init()`. New segments do not touch the core.
- **Configurable without editing JSON by hand.** `claude-statusline config` toggles every row on or off and picks the theme and icon set, with a live preview. Settings persist in `~/.claude/claude-statusline.json`.

## How Claude Code drives it

Claude Code launches the command from `~/.claude/settings.json` → `statusLine.command`, pipes a JSON descriptor to stdin, and renders whatever the program prints to stdout.

Example payload (only fields used by built-in segments are shown):

```json
{
  "workspace":      { "current_dir": "/Users/andrew/code/foo" },
  "model":          { "display_name": "Opus 5 (1M context)", "id": "claude-opus-5[1m]" },
  "context_window": { "used_percentage": 42.7 }
}
```

Output is plain text with ANSI true-color escapes; rows are separated by `\n`, no trailing newline.

## Install

```sh
make build                       # produces ./claude-statusline in the project dir
make install                     # copies it to ~/.claude/claude-statusline
make install PREFIX=/usr/local   # override the install location
```

Default install target is `$(PREFIX)/claude-statusline`, where `PREFIX` defaults to `$HOME/.claude`.

Then wire it into Claude Code by editing `~/.claude/settings.json`:

```json
"statusLine": {
  "type": "command",
  "command": "~/.claude/claude-statusline"
}
```

No arguments — everything is configured through `claude-statusline config`. Restart Claude Code (or just open a new session) and the new statusline takes effect on the next render.

### Make the command global

Claude Code calls the binary by the path in `settings.json`, so `~/.claude` never has to be on your `PATH` — but `claude-statusline config` is easier to reach when it is. Symlink the installed binary into a directory that already is:

```sh
ln -sfn ~/.claude/claude-statusline ~/.local/bin/claude-statusline   # or /usr/local/bin
claude-statusline config
```

One copy of the binary, two names for it: every later `make install` overwrites `~/.claude/claude-statusline` and the symlink follows. Remove it with `rm ~/.local/bin/claude-statusline`.

Alternatives, if a symlink is not to your taste:

```sh
make install PREFIX="$HOME/.local/bin"   # a second copy on PATH — remember to install to both
export PATH="$HOME/.claude:$PATH"        # in ~/.zshrc, if you want ~/.claude itself on PATH
```

Without any of this, run it by path: `~/.claude/claude-statusline config`.

## Commands

| Command                    | What it does                                                        |
|----------------------------|---------------------------------------------------------------------|
| `claude-statusline`        | Render one statusline from the JSON payload on stdin (what Claude Code calls) |
| `claude-statusline config` | Open the interactive settings editor                                |
| `claude-statusline help`   | Print usage, the settings path, and every registered segment / theme / icon set |

## Interactive configuration

```sh
~/.claude/claude-statusline config
```

A full-screen editor listing every registered segment, theme and icon set, with a preview rendered from the current directory:

```
claude-statusline config
~/.claude/claude-statusline.json

Segments
  [x] git
  [x] path
  [ ] meta

Theme
  (•) atom-one-dark
  ( ) dracula

Icons
  ( ) nerd
  ( ) none
  (•) plain

Preview
  ⎇ main
  📁 ~/Projects/paparoot/claude_statusline

↑/↓ move · space toggle · s save · q quit
```

| Key            | Action                                                                    |
|----------------|---------------------------------------------------------------------------|
| `↑` / `↓`, `k` / `j` | Move the cursor                                                     |
| `space`, `enter` | Toggle a segment on/off, or select the theme / icon set under the cursor |
| `s`            | Save to the settings file (the editor stays open)                          |
| `q`, `esc`     | Quit; with unsaved changes it asks once, press `q` again to discard        |
| `ctrl-c`       | Quit immediately, discarding unsaved changes                               |

The editor needs a real terminal — running it with stdin redirected exits with `` `config` needs an interactive terminal ``.

## Configuration file

Path: `~/.claude/claude-statusline.json`, or wherever `CLAUDE_STATUSLINE_CONFIG` points.

```json
{
  "theme": "dracula",
  "icon_font": "nerd",
  "segments": {
    "git": true,
    "path": true,
    "meta": false
  }
}
```

- Safe to hand-edit; `config` rewrites it in full (via a temp file + rename, so an interrupted save can't truncate a working config).
- Every field is optional. A missing file, a missing key, or a segment absent from `segments` falls back to the default — so a config written today keeps working when new segments are added (they render until switched off).
- A malformed file produces a stderr warning and the defaults; the statusline itself never breaks.
- Segment order is fixed by each segment's `Order` (git → path → meta) and is not part of the config.

### Environment overrides

Priority: **environment variable → config file → built-in default**. Env vars are meant for a single run (a shell session, a one-off `settings.json` command), not as the permanent home for settings.

| Variable                      | Default         | Purpose                                              |
|-------------------------------|-----------------|------------------------------------------------------|
| `CLAUDE_STATUSLINE_THEME`     | `atom-one-dark` | Pick a registered theme                              |
| `CLAUDE_STATUSLINE_ICON_FONT` | `plain`         | Pick a registered icon set                           |
| `CLAUDE_STATUSLINE_SEGMENTS`  | all enabled     | Comma-separated segments to render, in the order given — ignores the `segments` toggles |
| `CLAUDE_STATUSLINE_CONFIG`    | `~/.claude/claude-statusline.json` | Path to the settings file          |

Unknown theme / icon set → warning to stderr, fall back to the default. Unknown segment → warning to stderr, skip it. Errors **never** go to stdout — they would corrupt the statusline.

### Examples

```sh
# Persist a setup once
~/.claude/claude-statusline config

# Preview a theme for a single render, without touching the config
export CLAUDE_STATUSLINE_THEME=dracula

# Show only the model row for one render
echo '{...}' | ~/.claude/claude-statusline   # with CLAUDE_STATUSLINE_SEGMENTS=meta

# Keep a separate config for experiments
CLAUDE_STATUSLINE_CONFIG=/tmp/try.json ~/.claude/claude-statusline config
```

> **Migrating from the flag-based versions.** `--theme`, `--segments`, `--icon-font` and the `--list-*` flags are gone; `config` and `help` replace them. A leftover flag in `settings.json` is not fatal — the binary warns on stderr and renders with the configured settings — but drop it and run `config` once to keep the old look (`--icon-font nerd` becomes `"icon_font": "nerd"`).

## Running it by hand

You normally **don't run the binary yourself** — Claude Code spawns it on every statusline refresh, writes the JSON descriptor to its stdin, and reads ANSI text from its stdout.

To render manually (for testing, or to see a config change), pipe a payload:

```sh
echo '{"workspace":{"current_dir":"/Users/andrew"},"model":{"display_name":"Opus 5 (1M context)","id":"claude-opus-5[1m]"},"context_window":{"used_percentage":42.7}}' \
  | ~/.claude/claude-statusline
```

Running with no stdin (`~/.claude/claude-statusline < /dev/null`) is safe — every segment will return "no data" and the binary prints nothing.

## Bundled themes

| Name              | Source                                         |
|-------------------|------------------------------------------------|
| `atom-one-dark`   | [Atom One Dark](https://github.com/atom/atom/tree/master/packages/one-dark-syntax) — default, byte-identical to the previous JS statusline |
| `dracula`         | [Dracula](https://draculatheme.com)            |

## Bundled segments

Each segment is one row, rendered in `Order`. Any of them can be switched off in `config`.

| Name   | Order | What it shows                                                                                        |
|--------|-------|------------------------------------------------------------------------------------------------------|
| `git`  | 10    | Current branch with a worktree marker and `↑N` / `↓N` ahead/behind vs upstream. Skipped outside a repo. |
| `path` | 20    | Current working directory with `$HOME` collapsed to `~`. Always renders.                              |
| `meta` | 30    | Model family + version (e.g. `Opus 5`) and a context-usage percentage chip with colour-coded thresholds. |

## Bundled icon sets

The icon set decides which glyph each segment uses for its prefix and markers. Pick one whose glyphs your terminal font can actually render — that's why `plain` is the default.

| Name    | Requires       | Notes                                                                              |
|---------|----------------|------------------------------------------------------------------------------------|
| `plain` | nothing        | **Default.** Branch `⎇`, folder `📁`, model `🧊`, plus `⇟ ↑ ↓` markers. Universal Unicode — works with any modern terminal font. |
| `nerd`  | a Nerd Font    | Uses Private Use Area glyphs from [Nerd Fonts](https://www.nerdfonts.com) (`nf-dev-git_branch`, `nf-fa-folder`, `nf-fa-cube`, `nf-oct-file_submodule`, `nf-fa-arrow_up/down`). Renders as tofu without a Nerd Font. |
| `none`  | nothing        | All glyphs empty — text-only output. Useful for minimalist setups or terminals with poor emoji rendering. |

## Extending

Anything you register shows up in `claude-statusline config` and in `claude-statusline help` — no central wiring to edit.

### Add a theme

Create `theme/<name>.go`:

```go
package theme

func init() {
    Register(Palette{
        Name:   "solarized-dark",
        Red:    NewColor(220,  50,  47),
        Green:  NewColor(133, 153,   0),
        Yellow: NewColor(181, 137,   0),
        Blue:   NewColor( 38, 139, 210),
        Purple: NewColor(108, 113, 196),
        Cyan:   NewColor( 42, 161, 152),
        White:  NewColor(238, 232, 213),
        Gutter: NewColor(101, 123, 131),
        Bg2:    NewColor(  7,  54,  66),
        Bg3:    NewColor(  0,  43,  54),
        ChipBg: NewColor(  0,  43,  54),
        ChipFg: FgEscape(  0,  43,  54),
    })
}
```

That's it — `make build`, then pick it in `config`.

### Add an icon set

Create `icons/<name>.go`:

```go
package icons

func init() {
    Register(Set{
        Name:     "powerline-extra",
        Branch:   "", // your font's branch glyph
        Folder:   "",
        Model:    "",
        Worktree: "",
        Ahead:    "",
        Behind:   "",
    })
}
```

Any field left empty renders as no icon — segments use `icons.Prefix` to drop the leading space when a glyph is missing, so partial sets work too.

### Add a segment

Create `statusline/<name>.go`:

```go
package statusline

import (
    "time"

    "claude-statusline/icons"
    "claude-statusline/theme"
)

func init() {
    Register(Segment{Name: "time", Render: renderTime})
}

func renderTime(in Input, p theme.Palette, ic icons.Set) (string, bool) {
    return ChipOutline("⏱ "+time.Now().Format("15:04"), p.Cyan, p.Bg2), true
}
```

The new row renders on the next build (segments are enabled unless the config file says otherwise) and appears in the `config` editor.

- `Order` places the row: lower comes first. Leaving it 0 — as above — appends the segment after every ordered one, so a new file never has to renumber the built-ins.
- A segment that has no data should return `(_, false)` — the pipeline silently skips it.
- To prefix output with a configurable icon, call `icons.Prefix(ic.Foo, text)` — it omits the separator when the glyph is empty.

## Build targets

| Target           | What it does                                                    |
|------------------|-----------------------------------------------------------------|
| `make build`     | Build `./claude-statusline` (`-ldflags="-s -w"` strips symbols) |
| `make install`   | Build, then `install -m 755` to `$(PREFIX)/claude-statusline`   |
| `make clean`     | Remove the local `./claude-statusline` binary                   |
| `make test`      | Run `go test ./...` (no tests bundled yet — placeholder)        |
| `make fmt`       | Run `gofmt -w .` over the project                               |
| `make vet`       | Run `go vet ./...`                                              |

Override the install destination with `PREFIX`:

```sh
make install PREFIX=/usr/local       # → /usr/local/claude-statusline
make install PREFIX="$HOME/bin"      # → ~/bin/claude-statusline
```

## Project layout

```
.
├── main.go               # subcommand dispatch, settings resolution, render pipeline
├── config/               # persisted settings + interactive editor
│   ├── config.go         # Config type, Path / Load / Save
│   ├── ui.go             # full-screen editor with live preview
│   └── term.go           # raw mode via stty, key decoding, ANSI constants
├── statusline/           # core types + built-in segments
│   ├── input.go          # JSON contract with Claude Code
│   ├── segment.go        # Segment struct + Register / Get / List
│   ├── chip.go           # Chip / ChipOutline helpers
│   ├── git.go            # git segment
│   ├── path.go           # path segment
│   └── meta.go           # model + context percentage segment
├── icons/                # icon sets (glyph profiles)
│   ├── icons.go          # Set type + Register / Get / List / Prefix
│   ├── plain.go          # universal-Unicode set (default)
│   ├── nerd.go           # Nerd Font PUA set
│   └── none.go           # empty set — text only
└── theme/                # palettes
    ├── theme.go          # Color / Palette types + Register / Get / List
    ├── atom_one_dark.go
    └── dracula.go
```

All packages use only the Go standard library.

## Verification

```sh
make vet                    # go vet ./...
make build                  # compile
./claude-statusline help    # usage + every registered segment / theme / icon set

# Render the default layout from a sample payload
echo '{"workspace":{"current_dir":"/Users/andrew"},"model":{"display_name":"Opus 5 (1M context)","id":"claude-opus-5[1m]"},"context_window":{"used_percentage":42.7}}' \
  | ./claude-statusline

# Try settings without touching the real config
CLAUDE_STATUSLINE_CONFIG=/tmp/try.json ./claude-statusline config

# Drive it from inside an actual git repo to see the git row
cd ~/some/repo && echo "$(cat <<'EOF'
{"workspace":{"current_dir":"PWD"},"model":{"display_name":"Opus 5 (1M context)","id":"claude-opus-5[1m]"},"context_window":{"used_percentage":42.7}}
EOF
)" | sed "s|PWD|$PWD|" | ~/projects/paparoot/claude_statusline/claude-statusline
```

## License

Personal project — no license declared.
