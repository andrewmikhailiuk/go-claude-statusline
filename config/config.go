// Package config holds the persisted statusline settings — the on-disk shape
// of ~/.claude/claude-statusline.json — plus the interactive editor behind
// the `claude-statusline config` subcommand.
//
// A missing or partial file is always valid: every unset field falls back to
// a built-in default, and a segment absent from the toggle map renders. That
// way a config written today keeps working when new segments are added.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Built-in defaults, used whenever the config file leaves a field unset.
const (
	DefaultTheme    = "atom-one-dark"
	DefaultIconFont = "plain"
)

// EnvPath overrides the location of the settings file.
const EnvPath = "CLAUDE_STATUSLINE_CONFIG"

const fileName = "claude-statusline.json"

// Config is the settings file. Fields are omitted when empty so a
// hand-edited file stays readable.
type Config struct {
	Theme    string `json:"theme,omitempty"`
	IconFont string `json:"icon_font,omitempty"`

	// Segments maps a segment name to whether its row is rendered. A segment
	// missing from the map is enabled — newly added segments show up without
	// anyone having to edit the file.
	Segments map[string]bool `json:"segments,omitempty"`
}

// SegmentEnabled reports whether the named segment should render.
func (c Config) SegmentEnabled(name string) bool {
	if on, ok := c.Segments[name]; ok {
		return on
	}
	return true
}

// Path returns the settings file location: the EnvPath override, else
// $HOME/.claude/claude-statusline.json.
func Path() string {
	if p := os.Getenv(EnvPath); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fileName
	}
	return filepath.Join(home, ".claude", fileName)
}

// Load reads the settings file. A missing file is not an error — it yields
// the zero Config, which resolves to the built-in defaults. A malformed file
// is an error, so the caller can complain instead of silently ignoring it.
func Load(path string) (Config, error) {
	var c Config
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("%s: %w", Display(path), err)
	}
	return c, nil
}

// Save writes c to path as indented JSON, creating parent directories. The
// bytes land in a temp file first so an interrupted write cannot truncate a
// working config.
func Save(path string, c Config) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// Display shortens path for humans by collapsing $HOME to "~".
func Display(path string) string {
	if home, err := os.UserHomeDir(); err == nil && home != "" && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}
	return path
}
