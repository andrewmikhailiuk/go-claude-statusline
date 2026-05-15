// claude-statusline reads a JSON descriptor from stdin (as written by Claude
// Code) and prints an ANSI-coloured statusline to stdout. Themes, icon sets
// and segments are pluggable — see the theme, icons and statusline packages.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"claude-statusline/icons"
	"claude-statusline/statusline"
	"claude-statusline/theme"
)

const (
	defaultTheme    = "atom-one-dark"
	defaultSegments = "git,path,meta"
	defaultIconFont = "plain"

	envTheme    = "CLAUDE_STATUSLINE_THEME"
	envSegments = "CLAUDE_STATUSLINE_SEGMENTS"
	envIconFont = "CLAUDE_STATUSLINE_ICON_FONT"
)

func main() {
	themeFlag := flag.String("theme", "",
		"color theme name (env: "+envTheme+", default: "+defaultTheme+")")
	segmentsFlag := flag.String("segments", "",
		"comma-separated segments in order (env: "+envSegments+", default: "+defaultSegments+")")
	iconFontFlag := flag.String("icon-font", "",
		"icon set name (env: "+envIconFont+", default: "+defaultIconFont+")")
	listThemes := flag.Bool("list-themes", false, "list available themes and exit")
	listSegments := flag.Bool("list-segments", false, "list available segments and exit")
	listIconFonts := flag.Bool("list-icon-fonts", false, "list available icon sets and exit")
	flag.Parse()

	if *listThemes {
		for _, n := range theme.List() {
			fmt.Println(n)
		}
		return
	}
	if *listSegments {
		for _, n := range statusline.List() {
			fmt.Println(n)
		}
		return
	}
	if *listIconFonts {
		for _, n := range icons.List() {
			fmt.Println(n)
		}
		return
	}

	palette := resolveTheme(*themeFlag)
	iconSet := resolveIconFont(*iconFontFlag)
	segments := resolveSegments(*segmentsFlag)

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

func resolveTheme(flagValue string) theme.Palette {
	name := pick(flagValue, os.Getenv(envTheme), defaultTheme)
	if p, ok := theme.Get(name); ok {
		return p
	}
	fmt.Fprintf(os.Stderr, "claude-statusline: unknown theme %q, falling back to %q\n", name, defaultTheme)
	p, _ := theme.Get(defaultTheme)
	return p
}

func resolveIconFont(flagValue string) icons.Set {
	name := pick(flagValue, os.Getenv(envIconFont), defaultIconFont)
	if s, ok := icons.Get(name); ok {
		return s
	}
	fmt.Fprintf(os.Stderr, "claude-statusline: unknown icon font %q, falling back to %q\n", name, defaultIconFont)
	s, _ := icons.Get(defaultIconFont)
	return s
}

func resolveSegments(flagValue string) []string {
	raw := pick(flagValue, os.Getenv(envSegments), defaultSegments)
	parts := strings.Split(raw, ",")
	out := parts[:0]
	for _, s := range parts {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
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
