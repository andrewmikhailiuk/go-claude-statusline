# claude-statusline

A native Go statusline for [Claude Code](https://claude.com/claude-code) — zero runtime dependencies, pluggable themes, pluggable segments, configurable via flags or environment variables.

## Design goals

- **No runtime dependency.** A single self-contained binary; no node, no bun, no python.
- **Multiple themes.** Themes live in their own files and self-register on import. Adding one is a single-file change.
- **Pluggable segments.** Each row of the statusline (git, path, model+context, …) is a `Segment` registered through `init()`. New segments do not touch the core.
- **Configurable.** Theme and segment order picked via CLI flags or env vars. Reasonable defaults match the original JS statusline.

## How Claude Code drives it

Claude Code launches the command from `~/.claude/settings.json` → `statusLine.command`, pipes a JSON descriptor to stdin, and renders whatever the program prints to stdout.

Example payload (only fields used by built-in segments are shown):

```json
{
  "workspace":      { "current_dir": "/Users/andrew/code/foo" },
  "model":          { "display_name": "Opus 4.7 (1M context)", "id": "claude-opus-4-7[1m]" },
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

Restart Claude Code (or just open a new session) — the new statusline takes effect on the next render.

## Running

You normally **don't run the binary yourself** — Claude Code spawns it on every statusline refresh, writes the JSON descriptor to its stdin, and reads ANSI text from its stdout.

To run it manually (for testing or to preview a theme), pipe a payload:

```sh
echo '{"workspace":{"current_dir":"/Users/andrew"},"model":{"display_name":"Opus 4.7 (1M context)","id":"claude-opus-4-7[1m]"},"context_window":{"used_percentage":42.7}}' \
  | ~/.claude/claude-statusline

# Switch theme for one render
echo '{...}' | ~/.claude/claude-statusline --theme dracula

# Inspect registered options
~/.claude/claude-statusline --list-themes
~/.claude/claude-statusline --list-segments
~/.claude/claude-statusline --help
```

Running with no stdin (`~/.claude/claude-statusline < /dev/null`) is safe — every segment will return "no data" and the binary prints nothing.

## Configuration

Priority: **flag → environment variable → built-in default**.

| Flag                | Environment variable             | Default           | Purpose                                              |
|---------------------|----------------------------------|-------------------|------------------------------------------------------|
| `--theme NAME`      | `CLAUDE_STATUSLINE_THEME`        | `atom-one-dark`   | Pick a registered theme                              |
| `--segments LIST`   | `CLAUDE_STATUSLINE_SEGMENTS`     | `git,path,meta`   | Comma-separated segments in render order             |
| `--icon-font NAME`  | `CLAUDE_STATUSLINE_ICON_FONT`    | `plain`           | Pick a registered icon set                           |
| `--list-themes`     | —                                | —                 | Print registered themes and exit                     |
| `--list-segments`   | —                                | —                 | Print registered segments and exit                   |
| `--list-icon-fonts` | —                                | —                 | Print registered icon sets and exit                  |
| `-h`, `--help`      | —                                | —                 | Print usage and exit (provided by Go's `flag` package) |

Unknown theme / icon set → warning to stderr, fall back to the default. Unknown segment → warning to stderr, skip it. Errors **never** go to stdout — they would corrupt the statusline.

### Examples

```sh
# Switch theme via flag (per-statusline override)
~/.claude/claude-statusline --theme dracula

# Switch theme via env (persistent for the user's shell)
export CLAUDE_STATUSLINE_THEME=dracula

# Show only path and model — drop the git row
~/.claude/claude-statusline --segments path,meta

# Switch icon set (requires a matching font installed in the terminal)
~/.claude/claude-statusline --icon-font nerd
~/.claude/claude-statusline --icon-font none   # text only, no glyphs

# In settings.json: combine theme + segments + icons
"command": "~/.claude/claude-statusline --theme dracula --segments path,meta --icon-font nerd"

# Pipe a hand-crafted payload (useful while developing)
echo '{"workspace":{"current_dir":"/Users/andrew"},"model":{"display_name":"Opus 4.7 (1M context)","id":"claude-opus-4-7[1m]"},"context_window":{"used_percentage":42.7}}' \
  | ./claude-statusline
```

## Bundled themes

| Name              | Source                                         |
|-------------------|------------------------------------------------|
| `atom-one-dark`   | [Atom One Dark](https://github.com/atom/atom/tree/master/packages/one-dark-syntax) — default, byte-identical to the previous JS statusline |
| `dracula`         | [Dracula](https://draculatheme.com)            |

## Bundled segments

| Name   | What it shows                                                                                               |
|--------|-------------------------------------------------------------------------------------------------------------|
| `git`  | Current branch with a worktree marker and `↑N` / `↓N` ahead/behind vs upstream. Skipped outside a repo.     |
| `path` | Current working directory with `$HOME` collapsed to `~`. Always renders.                                    |
| `meta` | Model family + version (e.g. `Opus 4.7`) and a context-usage percentage chip with colour-coded thresholds.  |

## Bundled icon sets

The icon set decides which glyph each segment uses for its prefix and markers. Pick one whose glyphs your terminal font can actually render — that's why `plain` is the default.

| Name    | Requires       | Notes                                                                              |
|---------|----------------|------------------------------------------------------------------------------------|
| `plain` | nothing        | **Default.** Branch `⎇`, folder `📁`, model `🧊`, plus `⇟ ↑ ↓` markers. Universal Unicode — works with any modern terminal font. |
| `nerd`  | a Nerd Font    | Uses Private Use Area glyphs from [Nerd Fonts](https://www.nerdfonts.com) (`nf-dev-git_branch`, `nf-fa-folder`, `nf-fa-cube`, `nf-oct-file_submodule`, `nf-fa-arrow_up/down`). Renders as tofu without a Nerd Font. |
| `none`  | nothing        | All glyphs empty — text-only output. Useful for minimalist setups or terminals with poor emoji rendering. |

## Extending

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

That's it — `make build` and `./claude-statusline --theme solarized-dark` work immediately.

### Add an icon set

Create `icons/<name>.go`:

```go
package icons

func init() {
    Register(Set{
        Name:     "powerline-extra",
        Branch:   "", // your font's branch glyph
        Folder:   "",
        Model:    "",
        Worktree: "",
        Ahead:    "",
        Behind:   "",
    })
}
```

Then enable it: `--icon-font powerline-extra`. Any field left empty renders as no icon — segments use `icons.Prefix` to drop the leading space when a glyph is missing, so partial sets work too. The set will appear in the list returned by `--list-icon-fonts`.

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

Then enable it: `--segments git,path,meta,time` (or via the env var). The segment will appear in the list returned by `--list-segments`.

A segment that has no data should return `(_, false)` — the pipeline silently skips it. To prefix output with a configurable icon, call `icons.Prefix(ic.Foo, text)` — it omits the separator when the glyph is empty.

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
├── main.go               # flag/env parsing, stdin decoding, pipeline
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
./claude-statusline --list-themes
./claude-statusline --list-segments

# Render the default layout from a sample payload
echo '{"workspace":{"current_dir":"/Users/andrew"},"model":{"display_name":"Opus 4.7 (1M context)","id":"claude-opus-4-7[1m]"},"context_window":{"used_percentage":42.7}}' \
  | ./claude-statusline

# Drive it from inside an actual git repo to see the git row
cd ~/some/repo && echo "$(cat <<'EOF'
{"workspace":{"current_dir":"PWD"},"model":{"display_name":"Opus 4.7 (1M context)","id":"claude-opus-4-7[1m]"},"context_window":{"used_percentage":42.7}}
EOF
)" | sed "s|PWD|$PWD|" | ~/projects/paparoot/claude_statusline/claude-statusline
```

## License

Personal project — no license declared.
