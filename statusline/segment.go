package statusline

import (
	"sort"

	"claude-statusline/icons"
	"claude-statusline/theme"
)

// Segment renders one row of the statusline.
//
// Render returns the finished line and ok=true when the segment has something
// to show. Returning ok=false tells the pipeline to silently skip this row —
// use it when there's no data (e.g. cwd is not a git repo, model is missing).
type Segment struct {
	Name   string
	Render func(in Input, p theme.Palette, ic icons.Set) (line string, ok bool)
}

var registry = map[string]Segment{}

// Register adds s to the global segment registry. Panics on empty Name, nil
// Render, or duplicate registration — all programming errors in a segment file.
func Register(s Segment) {
	if s.Name == "" {
		panic("statusline: segment name required")
	}
	if s.Render == nil {
		panic("statusline: segment render required: " + s.Name)
	}
	if _, dup := registry[s.Name]; dup {
		panic("statusline: segment already registered: " + s.Name)
	}
	registry[s.Name] = s
}

// Get returns the segment registered under name. The boolean is false when
// no such segment exists.
func Get(name string) (Segment, bool) {
	s, ok := registry[name]
	return s, ok
}

// List returns the names of all registered segments, sorted alphabetically.
func List() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
