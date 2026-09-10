package config

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ANSI sequences used by the editor. Kept here so ui.go reads as layout.
const (
	enterAltScreen = "\x1b[?1049h\x1b[?25l" // alternate screen, cursor hidden
	leaveAltScreen = "\x1b[?25h\x1b[?1049l"
	clearScreen    = "\x1b[H\x1b[2J"

	sgrReset   = "\x1b[0m"
	sgrBold    = "\x1b[1m"
	sgrDim     = "\x1b[2m"
	sgrReverse = "\x1b[7m"
)

// errNoTTY is returned when stdin is not a terminal, so `config` cannot run.
var errNoTTY = errors.New("`config` needs an interactive terminal (stdin is not a tty)")

// rawMode switches the terminal to raw mode so single keypresses arrive
// without a newline and are not echoed. It shells out to stty rather than
// calling tcsetattr, which keeps the binary free of any dependency beyond
// the standard library. The returned function restores the saved mode.
func rawMode() (restore func(), err error) {
	saved, err := stty("-g")
	if err != nil {
		return nil, errNoTTY
	}
	if _, err := stty("raw", "-echo"); err != nil {
		return nil, err
	}
	return func() { _, _ = stty(saved) }, nil
}

// stty runs stty against the real terminal and returns its trimmed output.
func stty(args ...string) (string, error) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// key is one editor command, decoded from the raw byte stream.
type key int

const (
	keyNone key = iota
	keyUp
	keyDown
	keyToggle
	keySave
	keyQuit  // q or Esc — asks before discarding changes
	keyAbort // Ctrl-C, Ctrl-D or a closed stdin — leaves at once
)

// readKey blocks until the next keypress and maps it to a key. Unknown input
// yields keyNone, and the caller redraws and waits again.
func readKey(r *bufio.Reader) key {
	b, err := r.ReadByte()
	if err != nil {
		return keyAbort
	}
	switch b {
	case 'k', 'K':
		return keyUp
	case 'j', 'J':
		return keyDown
	case ' ', '\r', '\n':
		return keyToggle
	case 's', 'S':
		return keySave
	case 'q', 'Q':
		return keyQuit
	case 0x03, 0x04: // Ctrl-C, Ctrl-D
		return keyAbort
	case 0x1b:
		// A bare Esc arrives alone; an arrow key arrives as ESC [ A/B, so the
		// rest of the sequence is already buffered by the time we get here.
		if r.Buffered() == 0 {
			return keyQuit
		}
		if intro, err := r.ReadByte(); err != nil || (intro != '[' && intro != 'O') {
			return keyNone
		}
		final, err := r.ReadByte()
		if err != nil {
			return keyAbort
		}
		switch final {
		case 'A':
			return keyUp
		case 'B':
			return keyDown
		}
		return keyNone
	}
	return keyNone
}
