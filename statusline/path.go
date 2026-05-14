package statusline

import (
	"os"
	"strings"

	"claude-statusline/theme"
)

const iconFolder = "" // nf-fa-folder

func init() {
	Register(Segment{Name: "path", Render: renderPath})
}

func renderPath(in Input, p theme.Palette) (string, bool) {
	cwd := in.Workspace.CurrentDir
	if cwd == "" {
		cwd = "~"
	} else if home := os.Getenv("HOME"); home != "" && strings.HasPrefix(cwd, home) {
		cwd = "~" + cwd[len(home):]
	}
	return ChipOutline(iconFolder+" "+cwd, p.White, p.ChipBg), true
}
