package color

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Colorizer applies LS_COLORS-based ANSI SGR sequences to filenames.
type Colorizer struct {
	enabled  bool
	lsColors map[string]string
}

// New creates a Colorizer reading LS_COLORS from the environment.
func New(mode string, w io.Writer) *Colorizer {
	return NewFromLSColors(mode, w, os.Getenv("LS_COLORS"))
}

// NewFromLSColors creates a Colorizer with an explicit LS_COLORS string.
// Intended for tests that need deterministic color output without touching env vars.
func NewFromLSColors(mode string, w io.Writer, lsColors string) *Colorizer {
	return &Colorizer{
		enabled:  isEnabled(mode, w),
		lsColors: ParseLSColors(lsColors),
	}
}

// Enabled reports whether color output is active.
func (c *Colorizer) Enabled() bool { return c.enabled }

// Apply wraps name in the appropriate ANSI SGR sequence based on its file mode.
// mode must include permission bits (from entry.Info().Mode()), not just type bits.
// Returns name unchanged when color is disabled or no mapping exists.
func (c *Colorizer) Apply(name string, mode fs.FileMode) string {
	if !c.enabled {
		return name
	}
	sgr := c.sgrFor(name, mode)
	if sgr == "" {
		return name
	}
	return "\033[" + sgr + "m" + name + "\033[0m"
}

func (c *Colorizer) sgrFor(name string, mode fs.FileMode) string {
	switch {
	case mode&fs.ModeDir != 0:
		return c.lsColors["di"]
	case mode&fs.ModeSymlink != 0:
		return c.lsColors["ln"]
	case mode&0111 != 0:
		return c.lsColors["ex"]
	default:
		if ext := filepath.Ext(name); ext != "" {
			if sgr, ok := c.lsColors["*"+ext]; ok {
				return sgr
			}
		}
		return c.lsColors["fi"]
	}
}

// ParseLSColors parses an LS_COLORS string into a lookup map.
// Format: "di=01;34:ln=01;36:*.go=00;36:..."
// Entries with empty values are silently dropped.
func ParseLSColors(env string) map[string]string {
	m := make(map[string]string)
	for _, pair := range strings.Split(env, ":") {
		k, v, ok := strings.Cut(pair, "=")
		if ok && k != "" && v != "" {
			m[k] = v
		}
	}
	return m
}

func isEnabled(mode string, w io.Writer) bool {
	switch mode {
	case "never":
		return false
	case "always":
		return true
	default: // "auto" or empty
		if os.Getenv("NO_COLOR") != "" {
			return false
		}
		if os.Getenv("TERM") == "dumb" {
			return false
		}
		return isTTY(w)
	}
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice != 0
}
