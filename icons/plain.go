package icons

// Plain set — basic Unicode glyphs that any modern terminal font renders.
// No special icon font required. This is the default.
func init() {
	Register(Set{
		Name:     "plain",
		Branch:   "⎇",          // ⎇ BRANCH
		Folder:   "\U0001f4c1", // 📁 FILE FOLDER
		Model:    "\U0001f9ca", // 🧊 ICE
		Worktree: "⇟",          // ⇟ DOWNWARDS ARROW WITH DOUBLE STROKE
		Ahead:    "↑",          // ↑
		Behind:   "↓",          // ↓
	})
}
