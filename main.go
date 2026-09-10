// claude-statusline reads a JSON descriptor from stdin (as written by Claude
// Code) and prints an ANSI-coloured statusline to stdout. Themes, icon sets
// and segments are pluggable — see the theme, icons and statusline packages.
//
// Settings live in ~/.claude/claude-statusline.json and are edited with
// `claude-statusline config`; environment variables override the file for a
// single run.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"claude-statusline/config"
	"claude-statusline/icons"
	"claude-statusline/statusline"
	"claude-statusline/theme"
)

// Per-run overrides. Handy in a shell or a one-off settings.json command;
// the config file is the place for permanent choices.
const (
	envTheme    = "CLAUDE_STATUSLINE_THEME"
	envSegments = "CLAUDE_STATUSLINE_SEGMENTS"
	envIconFont = "CLAUDE_STATUSLINE_ICON_FONT"
)

func main() {
	if args := os.Args[1:]; len(args) > 0 {
		switch args[0] {
		case "config":
			if err := config.Edit(config.Path()); err != nil {
				fmt.Fprintf(os.Stderr, "claude-statusline: %v\n", err)
				os.Exit(1)
			}
			return
		case "help", "-h", "--help":
			usage(os.Stdout)
			return
		}
		if !strings.HasPrefix(args[0], "-") {
			fmt.Fprintf(os.Stderr, "claude-statusline: unknown command %q\n\n", args[0])
			usage(os.Stderr)
			os.Exit(2)
		}
		// Flags were replaced by the config file. A stale settings.json must
		// not break the statusline, so warn on stderr and render anyway.
		fmt.Fprintf(os.Stderr, "claude-statusline: flags are no longer supported (%s); "+
			"run `claude-statusline config` and drop them from settings.json\n", strings.Join(args, " "))
	}
	render()
}

// render builds the statusline from the config file, the environment and the
// JSON payload on stdin. Nothing but the statusline itself reaches stdout.
func render() {
	cfg, err := config.Load(config.Path())
	if err != nil {
		fmt.Fprintf(os.Stderr, "claude-statusline: %v\n", err)
	}

	palette := resolveTheme(cfg)
	iconSet := resolveIconFont(cfg)
	segments := resolveSegments(cfg)

	var in statusline.Input
	// Decode errors are deliberately ignored: the statusline must never abort,
	// it should just print whatever rows it can build from an empty Input.
	_ = json.NewDecoder(os.Stdin).Decode(&in)

	rows := make([]string, 0, len(segments))
	for _, name := range segments {
		seg, ok := statusline.Get(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "claude-statusline: unknown segment %q\n", name)
			continue
		}
		if line, ok := seg.Render(in, palette, iconSet); ok {
			rows = append(rows, line)
		}
	}

	// No trailing newline — matches the original JS statusline.
	fmt.Print(strings.Join(rows, "\n"))
}

func resolveTheme(cfg config.Config) theme.Palette {
	name := pick(os.Getenv(envTheme), cfg.Theme, config.DefaultTheme)
	if p, ok := theme.Get(name); ok {
		return p
	}
	fmt.Fprintf(os.Stderr, "claude-statusline: unknown theme %q, falling back to %q\n", name, config.DefaultTheme)
	p, _ := theme.Get(config.DefaultTheme)
	return p
}

func resolveIconFont(cfg config.Config) icons.Set {
	name := pick(os.Getenv(envIconFont), cfg.IconFont, config.DefaultIconFont)
	if s, ok := icons.Get(name); ok {
		return s
	}
	fmt.Fprintf(os.Stderr, "claude-statusline: unknown icon font %q, falling back to %q\n", name, config.DefaultIconFont)
	s, _ := icons.Get(config.DefaultIconFont)
	return s
}

// resolveSegments returns the rows to render, in order. The environment
// variable wins when set (its order is honoured as given); otherwise every
// registered segment renders unless the config file switched it off.
func resolveSegments(cfg config.Config) []string {
	if raw := os.Getenv(envSegments); raw != "" {
		var out []string
		for _, s := range strings.Split(raw, ",") {
			if s = strings.TrimSpace(s); s != "" {
				out = append(out, s)
			}
		}
		return out
	}

	all := statusline.List()
	out := make([]string, 0, len(all))
	for _, name := range all {
		if cfg.SegmentEnabled(name) {
			out = append(out, name)
		}
	}
	return out
}

// pick returns the first non-empty string from values.
func pick(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `claude-statusline — a native statusline for Claude Code

Usage:
  claude-statusline           render one statusline from the JSON payload on stdin
  claude-statusline config    edit the settings interactively (segments, theme, icons)
  claude-statusline help      print this message

Settings file (edited by `+"`config`"+`, safe to hand-edit):
  %s

Environment overrides, for a single run:
  %-29s theme name (default: %s)
  %-29s icon set name (default: %s)
  %-29s comma-separated segments to render, in order
  %-29s path to the settings file

Registered segments:  %s
Registered themes:    %s
Registered icon sets: %s
`,
		config.Display(config.Path()),
		envTheme, config.DefaultTheme,
		envIconFont, config.DefaultIconFont,
		envSegments,
		config.EnvPath,
		strings.Join(statusline.List(), ", "),
		strings.Join(theme.List(), ", "),
		strings.Join(icons.List(), ", "))
}
