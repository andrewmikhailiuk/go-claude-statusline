package statusline

import "claude-statusline/theme"

// Chip renders a solid color block: dark palette text on the given color's
// background, padded with one space on each side. Used for high-emphasis
// elements like the context-percentage indicator.
func Chip(text string, c theme.Color, p theme.Palette) string {
	return c.Bg + p.ChipFg + " " + text + " " + theme.Reset
}

// ChipOutline renders a coloured-text-on-darker-background block. The bg
// argument is the backdrop (usually p.Bg2 or p.Bg3); pass a different value
// to glue chips together on a shared row.
func ChipOutline(text string, fg theme.Color, bg theme.Color) string {
	return bg.Bg + fg.Fg + " " + text + " " + theme.Reset
}
