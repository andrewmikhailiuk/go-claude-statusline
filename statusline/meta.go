package statusline

import (
	"fmt"
	"regexp"
	"strings"

	"claude-statusline/theme"
)

const iconModel = "" // nf-fa-cube

func init() {
	Register(Segment{Name: "meta", Render: renderMeta})
}

// modelVerRe captures the major/minor version digits from a model id like
// "claude-opus-4-7[1m]" → ["4", "7"].
var modelVerRe = regexp.MustCompile(`(?i)(?:opus|sonnet|haiku)[_-](\d+)[_-](\d+)`)

// parseModel resolves a human-friendly "Family Major.Minor" label from
// Claude Code's display_name and id fields. Returns "" when neither yields
// a recognisable family — caller should drop the chip in that case.
func parseModel(displayName, id string) string {
	var label string
	for _, fam := range []string{"Opus", "Sonnet", "Haiku"} {
		if strings.Contains(displayName, fam) {
			label = fam
			break
		}
	}
	if label == "" {
		lower := strings.ToLower(id)
		for _, fam := range []string{"opus", "sonnet", "haiku"} {
			if strings.Contains(lower, fam) {
				label = strings.ToUpper(fam[:1]) + fam[1:]
				break
			}
		}
	}
	if label == "" {
		return ""
	}
	if m := modelVerRe.FindStringSubmatch(id); len(m) == 3 {
		return fmt.Sprintf("%s %s.%s", label, m[1], m[2])
	}
	return label
}

func renderMeta(in Input, p theme.Palette) (string, bool) {
	var parts []string

	if m := parseModel(in.Model.DisplayName, in.Model.ID); m != "" {
		parts = append(parts, ChipOutline(iconModel+" "+m, p.Yellow, p.Bg2))
	}

	if in.ContextWindow != nil {
		pct := in.ContextWindow.UsedPercentage
		var color theme.Color
		switch {
		case pct >= 90:
			color = p.Red
		case pct >= 70:
			color = p.Yellow
		case pct >= 40:
			color = p.Blue
		default:
			color = p.Green
		}
		var pctStr string
		if pct >= 90 {
			pctStr = fmt.Sprintf("%.1f", pct)
		} else {
			pctStr = fmt.Sprintf("%d", int(pct+0.5))
		}
		parts = append(parts, Chip(pctStr+"%", color, p))
	}

	if len(parts) == 0 {
		return "", false
	}
	return strings.Join(parts, ""), true
}
