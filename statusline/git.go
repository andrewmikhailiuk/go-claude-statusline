package statusline

import (
	"os/exec"
	"strconv"
	"strings"

	"claude-statusline/icons"
	"claude-statusline/theme"
)

const worktreeSubdir = "/.git/worktrees/"

func init() {
	Register(Segment{Name: "git", Order: 10, Render: renderGit})
}

func renderGit(in Input, p theme.Palette, ic icons.Set) (string, bool) {
	cwd := in.Workspace.CurrentDir
	if cwd == "" {
		return "", false
	}
	if runIn(cwd, "git", "rev-parse", "--is-inside-work-tree") != "true" {
		return "", false
	}

	branch := runIn(cwd, "git", "branch", "--show-current")
	if branch == "" {
		// Detached HEAD or fresh repo — skip rather than print a blank chip.
		return "", false
	}

	text := branch
	// Worktree marker: gitdir lives under the parent's /.git/worktrees/.
	if gitDir := runIn(cwd, "git", "rev-parse", "--git-dir"); strings.Contains(gitDir, worktreeSubdir) && ic.Worktree != "" {
		text += " " + ic.Worktree
	}

	// Ahead/behind vs upstream. Missing upstream makes git exit non-zero;
	// runIn swallows that and we just skip the markers.
	if lr := runIn(cwd, "git", "rev-list", "--left-right", "--count", "@{u}...HEAD"); lr != "" {
		parts := strings.Split(lr, "\t")
		if len(parts) == 2 {
			behind, _ := strconv.Atoi(parts[0])
			ahead, _ := strconv.Atoi(parts[1])
			if ahead > 0 {
				text += " " + ic.Ahead + strconv.Itoa(ahead)
			}
			if behind > 0 {
				text += " " + ic.Behind + strconv.Itoa(behind)
			}
		}
	}

	return ChipOutline(icons.Prefix(ic.Branch, text), p.Green, p.Bg3), true
}

// runIn executes name with args in the given working directory, discarding
// stderr. Returns the trimmed stdout, or "" on any error — segments treat
// empty as "no data" and degrade gracefully.
func runIn(cwd, name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
