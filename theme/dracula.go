package theme

// Dracula — official palette (https://draculatheme.com/contribute). Included
// as a second theme to validate the registry and to give users a built-in
// alternative without an extra config file.
func init() {
	Register(Palette{
		Name:   "dracula",
		Red:    NewColor(255, 85, 85),   // #ff5555
		Green:  NewColor(80, 250, 123),  // #50fa7b
		Yellow: NewColor(241, 250, 140), // #f1fa8c
		Blue:   NewColor(189, 147, 249), // #bd93f9 — Dracula has no pure blue; purple stands in
		Purple: NewColor(255, 121, 198), // #ff79c6 — "pink" in Dracula naming
		Cyan:   NewColor(139, 233, 253), // #8be9fd
		White:  NewColor(248, 248, 242), // #f8f8f2
		Gutter: NewColor(98, 114, 164),  // #6272a4 — "comment"
		Bg2:    NewColor(68, 71, 90),    // #44475a — "current line"
		Bg3:    NewColor(40, 42, 54),    // #282a36 — base background
		ChipBg: NewColor(33, 34, 44),    // #21222c — deepest backdrop
		ChipFg: FgEscape(40, 42, 54),    // dark text on bright chips
	})
}
