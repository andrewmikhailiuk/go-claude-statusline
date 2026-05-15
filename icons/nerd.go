package icons

// Nerd set — Private Use Area glyphs from a Nerd Font
// (https://www.nerdfonts.com). Renders correctly only when the user's
// terminal font is a Nerd Font patch.
func init() {
	Register(Set{
		Name:     "nerd",
		Branch:   "", // nf-dev-git_branch
		Folder:   "", // nf-fa-folder
		Model:    "", // nf-fa-cube
		Worktree: "", // nf-oct-file_submodule
		Ahead:    "", // nf-fa-arrow_up
		Behind:   "", // nf-fa-arrow_down
	})
}
