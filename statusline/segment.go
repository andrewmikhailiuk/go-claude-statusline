package statusline

import (
	"math"
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
	Name string

	// Order places the row in the statusline: lower comes first. Leave it 0
	// (the zero value) and the segment is appended after every ordered one,
	// so a new segment file never has to renumber the built-ins.
	Order int

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

// List returns the names of all registered segments in render order: by
// Order, then alphabetically among segments that share one.
func List() []string {
	segs := make([]Segment, 0, len(registry))
	for _, s := range registry {
		segs = append(segs, s)
	}
	sort.Slice(segs, func(i, j int) bool {
		oi, oj := orderKey(segs[i].Order), orderKey(segs[j].Order)
		if oi != oj {
			return oi < oj
		}
		return segs[i].Name < segs[j].Name
	})

	names := make([]string, 0, len(segs))
	for _, s := range segs {
		names = append(names, s.Name)
	}
	return names
}

// orderKey maps the zero Order to "after everything ordered", so segments
// registered without an explicit Order land at the end of the statusline.
func orderKey(order int) int {
	if order == 0 {
		return math.MaxInt
	}
	return order
}
