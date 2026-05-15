// Package icons defines icon sets used by statusline segments.
//
// An icon set is a Set value registered under a unique name. New sets are
// added by creating a file in this package that calls Register from an init()
// function — no central wiring needs editing. Empty fields render as no icon.
package icons

import "sort"

// Set bundles the glyphs used by built-in segments. Fields are semantic so
// segments stay decoupled from any concrete set. A "" field means "no icon" —
// segments must omit the leading separator in that case (use Prefix).
type Set struct {
	Name string

	Branch   string
	Folder   string
	Model    string
	Worktree string
	Ahead    string
	Behind   string
}

// Prefix returns "icon text" when icon is non-empty, otherwise just text.
// Lets segments support sets with missing glyphs without dangling spaces.
func Prefix(icon, text string) string {
	if icon == "" {
		return text
	}
	return icon + " " + text
}

var registry = map[string]Set{}

// Register adds s to the global icon-set registry. Panics on empty Name or
// duplicate registration — both indicate a programming error in a set file.
func Register(s Set) {
	if s.Name == "" {
		panic("icons: set name required")
	}
	if _, dup := registry[s.Name]; dup {
		panic("icons: set already registered: " + s.Name)
	}
	registry[s.Name] = s
}

// Get returns the set registered under name. The boolean is false when no
// such set exists.
func Get(name string) (Set, bool) {
	s, ok := registry[name]
	return s, ok
}

// List returns the names of all registered sets, sorted alphabetically.
func List() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
