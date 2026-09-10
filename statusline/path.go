package statusline

import (
	"os"
	"strings"

	"claude-statusline/icons"
	"claude-statusline/theme"
)

func init() {
	Register(Segment{Name: "path", Order: 20, Render: renderPath})
}

func renderPath(in Input, p theme.Palette, ic icons.Set) (string, bool) {
	cwd := in.Workspace.CurrentDir
	if cwd == "" {
		cwd = "~"
	} else if home := os.Getenv("HOME"); home != "" && strings.HasPrefix(cwd, home) {
		cwd = "~" + cwd[len(home):]
	}
	return ChipOutline(icons.Prefix(ic.Folder, cwd), p.White, p.ChipBg), true
}
