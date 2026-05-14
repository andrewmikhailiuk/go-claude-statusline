package theme

// Atom One Dark — values copied verbatim from the original JS statusline so
// the new binary is pixel-identical to the previous bun-driven script.
func init() {
	Register(Palette{
		Name:   "atom-one-dark",
		Red:    NewColor(224, 108, 117),
		Green:  NewColor(152, 195, 121),
		Yellow: NewColor(229, 192, 123),
		Blue:   NewColor(97, 175, 239),
		Purple: NewColor(198, 120, 221),
		Cyan:   NewColor(86, 182, 194),
		White:  NewColor(171, 178, 191),
		Gutter: NewColor(99, 109, 131),
		Bg2:    NewColor(53, 59, 69),
		Bg3:    NewColor(30, 33, 39),
		ChipBg: NewColor(38, 42, 50),
		ChipFg: FgEscape(40, 44, 52),
	})
}
