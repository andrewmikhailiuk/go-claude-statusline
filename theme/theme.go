// Package theme defines color palettes used by the statusline renderer.
//
// A theme is a Palette value registered under a unique name. New themes are
// added by creating a file in this package that calls Register from an init()
// function — no central wiring needs editing.
package theme

import (
	"fmt"
	"sort"
)

// Reset is the ANSI escape that clears all attributes.
const Reset = "\x1b[0m"

// Color bundles the foreground and background ANSI escapes for a single RGB
// triple, so callers can pick whichever side they need without recomputing.
type Color struct {
	Fg string
	Bg string
}

// FgEscape returns the ANSI true-color foreground escape for the given RGB.
func FgEscape(r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// BgEscape returns the ANSI true-color background escape for the given RGB.
func BgEscape(r, g, b int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// NewColor builds a Color with both Fg and Bg escapes from a single RGB triple.
func NewColor(r, g, b int) Color {
	return Color{Fg: FgEscape(r, g, b), Bg: BgEscape(r, g, b)}
}

// Palette is the full set of colors a theme provides. Field names are
// semantic (Red, Green, …, Bg2, Bg3, ChipBg) so segments stay decoupled from
// any concrete palette.
type Palette struct {
	Name string

	Red    Color
	Green  Color
	Yellow Color
	Blue   Color
	Purple Color
	Cyan   Color
	White  Color
	Gutter Color

	Bg2    Color // raised background — used as outline-chip backdrop
	Bg3    Color // deeper background — used behind the git chip
	ChipBg Color // backdrop for the path chip

	ChipFg string // dark text used on top of bright solid chips
}

var registry = map[string]Palette{}

// Register adds p to the global theme registry. It panics if Name is empty or
// already taken — both indicate a programming error in a theme file.
func Register(p Palette) {
	if p.Name == "" {
		panic("theme: palette name required")
	}
	if _, dup := registry[p.Name]; dup {
		panic("theme: palette already registered: " + p.Name)
	}
	registry[p.Name] = p
}

// Get returns the palette registered under name. The boolean is false when no
// such theme exists.
func Get(name string) (Palette, bool) {
	p, ok := registry[name]
	return p, ok
}

// List returns the names of all registered palettes, sorted alphabetically.
func List() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
