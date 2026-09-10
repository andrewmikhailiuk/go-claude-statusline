package config

import (
	"bufio"
	"fmt"
	"os"

	"claude-statusline/icons"
	"claude-statusline/statusline"
	"claude-statusline/theme"
)

// Edit runs the interactive settings editor against path: one full-screen
// list of every registered segment, theme and icon set, with a live preview
// of the resulting statusline. It returns when the user quits; changes reach
// disk only when they press "s".
func Edit(path string) error {
	cfg, err := Load(path)
	if err != nil {
		// A malformed file must not lock the user out of the editor: start
		// from defaults and let a save overwrite it.
		fmt.Fprintf(os.Stderr, "claude-statusline: %v\nstarting from defaults\n", err)
		cfg = Config{}
	}

	restore, err := rawMode()
	if err != nil {
		return err
	}
	out := bufio.NewWriter(os.Stdout)
	in := bufio.NewReader(os.Stdin)
	e := newEditor(path, cfg)

	out.WriteString(enterAltScreen)
	defer func() {
		out.WriteString(leaveAltScreen)
		_ = out.Flush()
		restore()
		// The alternate screen is gone by now, so leave a one-line trace of
		// what happened on the real one.
		fmt.Println(e.outcome())
	}()

	for {
		e.draw(out)
		if err := out.Flush(); err != nil {
			return err
		}

		k := readKey(in)
		if k != keyQuit {
			// Any other key cancels a pending quit confirmation.
			e.confirmQuit = false
		}
		switch k {
		case keyUp:
			e.move(-1)
		case keyDown:
			e.move(1)
		case keyToggle:
			e.toggle()
		case keySave:
			e.save()
		case keyAbort:
			return nil
		case keyQuit:
			if e.dirty && !e.confirmQuit {
				e.confirmQuit = true
				e.status = "unsaved changes — q again to discard, s to save"
				continue
			}
			return nil
		}
	}
}

// itemKind tells the editor what a row does when toggled: segments are
// independent checkboxes, themes and icon sets are single-choice.
type itemKind int

const (
	kindSegment itemKind = iota
	kindTheme
	kindIcon
)

func (k itemKind) title() string {
	switch k {
	case kindSegment:
		return "Segments"
	case kindTheme:
		return "Theme"
	default:
		return "Icons"
	}
}

type item struct {
	kind itemKind
	name string
}

type editor struct {
	path   string
	cfg    Config
	items  []item
	cursor int

	dirty       bool
	confirmQuit bool
	status      string
	saved       bool
}

// newEditor materialises every registered name as a row and fills in the
// values the config file left implicit, so what the editor shows is exactly
// what a save will write.
func newEditor(path string, cfg Config) *editor {
	segments := statusline.List()

	if cfg.Segments == nil {
		cfg.Segments = make(map[string]bool, len(segments))
	}
	for _, n := range segments {
		if _, ok := cfg.Segments[n]; !ok {
			cfg.Segments[n] = true
		}
	}
	if _, ok := theme.Get(cfg.Theme); !ok {
		cfg.Theme = DefaultTheme
	}
	if _, ok := icons.Get(cfg.IconFont); !ok {
		cfg.IconFont = DefaultIconFont
	}

	e := &editor{path: path, cfg: cfg}
	for _, n := range segments {
		e.items = append(e.items, item{kindSegment, n})
	}
	for _, n := range theme.List() {
		e.items = append(e.items, item{kindTheme, n})
	}
	for _, n := range icons.List() {
		e.items = append(e.items, item{kindIcon, n})
	}
	return e
}

func (e *editor) move(delta int) {
	e.cursor += delta
	switch {
	case e.cursor < 0:
		e.cursor = 0
	case e.cursor >= len(e.items):
		e.cursor = len(e.items) - 1
	}
}

func (e *editor) toggle() {
	if len(e.items) == 0 {
		return
	}
	it := e.items[e.cursor]
	switch it.kind {
	case kindSegment:
		e.cfg.Segments[it.name] = !e.cfg.Segments[it.name]
	case kindTheme:
		e.cfg.Theme = it.name
	case kindIcon:
		e.cfg.IconFont = it.name
	}
	e.dirty = true
	e.status = ""
}

func (e *editor) save() {
	if err := Save(e.path, e.cfg); err != nil {
		e.status = "save failed: " + err.Error()
		return
	}
	e.dirty = false
	e.saved = true
	e.status = "saved " + Display(e.path)
}

// outcome is the single line printed once the editor has torn down.
func (e *editor) outcome() string {
	switch {
	case e.dirty:
		return "discarded unsaved changes"
	case e.saved:
		return "saved " + Display(e.path)
	default:
		return "no changes"
	}
}

func (e *editor) draw(w *bufio.Writer) {
	w.WriteString(clearScreen)
	line(w, sgrBold+"claude-statusline config"+sgrReset)
	line(w, sgrDim+Display(e.path)+sgrReset)

	last := itemKind(-1)
	for i, it := range e.items {
		if it.kind != last {
			line(w, "")
			line(w, sgrDim+it.kind.title()+sgrReset)
			last = it.kind
		}
		line(w, e.itemLine(i, it))
	}

	line(w, "")
	line(w, sgrDim+"Preview"+sgrReset)
	for _, row := range e.preview() {
		line(w, "  "+row)
	}

	line(w, "")
	line(w, e.status)
	line(w, sgrDim+"↑/↓ move · space toggle · s save · q quit"+sgrReset)
}

func (e *editor) itemLine(i int, it item) string {
	mark := "[ ]"
	switch it.kind {
	case kindSegment:
		if e.cfg.Segments[it.name] {
			mark = "[x]"
		}
	case kindTheme:
		mark = radio(e.cfg.Theme == it.name)
	case kindIcon:
		mark = radio(e.cfg.IconFont == it.name)
	}

	text := fmt.Sprintf("%s %s", mark, it.name)
	if i == e.cursor {
		return sgrReverse + " " + text + " " + sgrReset
	}
	return "  " + text
}

func radio(on bool) string {
	if on {
		return "(•)"
	}
	return "( )"
}

// preview renders the enabled segments exactly as the statusline will, using
// the theme and icon set currently selected in the editor.
func (e *editor) preview() []string {
	p, _ := theme.Get(e.cfg.Theme)
	ic, _ := icons.Get(e.cfg.IconFont)
	in := previewInput()

	var rows []string
	for _, name := range statusline.List() {
		if !e.cfg.SegmentEnabled(name) {
			continue
		}
		seg, ok := statusline.Get(name)
		if !ok {
			continue
		}
		if row, ok := seg.Render(in, p, ic); ok {
			rows = append(rows, row)
		}
	}
	if len(rows) == 0 {
		return []string{sgrDim + "(every row disabled — the statusline prints nothing)" + sgrReset}
	}
	return rows
}

// previewInput fakes the payload Claude Code would send. The real cwd keeps
// the git and path rows honest; model and percentage are stand-ins.
func previewInput() statusline.Input {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}
	var in statusline.Input
	in.Workspace.CurrentDir = cwd
	in.Model = statusline.Model{DisplayName: "Opus 5 (1M context)", ID: "claude-opus-5[1m]"}
	in.ContextWindow = &statusline.ContextWindow{UsedPercentage: 42.7}
	return in
}

// line writes one row of the frame. Raw mode leaves the cursor where it is
// on "\n", hence the explicit carriage return.
func line(w *bufio.Writer, s string) {
	w.WriteString(s)
	w.WriteString("\x1b[K\r\n")
}
