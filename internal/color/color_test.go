package color_test

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/go-spass/baum/internal/color"
)

const testLSColors = "di=01;34:ln=01;36:ex=01;32:*.go=00;36:fi=00;37"

func TestParseLSColors(t *testing.T) {
	cases := []struct {
		input string
		key   string
		want  string
	}{
		{"di=01;34:ln=01;36", "di", "01;34"},
		{"di=01;34:ln=01;36", "ln", "01;36"},
		{"*.go=00;36:*.md=00;37", "*.go", "00;36"},
		{"*.go=00;36:*.md=00;37", "*.md", "00;37"},
		{"", "di", ""},                              // empty input
		{"di=:ln=01;36", "di", ""},                  // empty value ignored
		{"malformed:ln=01;36", "malformed", ""},      // no '='
		{"di=01;34:ln=01;36", "ex", ""},             // missing key
	}
	for _, tc := range cases {
		m := color.ParseLSColors(tc.input)
		if got := m[tc.key]; got != tc.want {
			t.Errorf("ParseLSColors(%q)[%q] = %q, want %q", tc.input, tc.key, got, tc.want)
		}
	}
}

func TestApply_never(t *testing.T) {
	var buf bytes.Buffer
	clr := color.NewFromLSColors("never", &buf, testLSColors)

	cases := []struct {
		name string
		mode fs.FileMode
	}{
		{"mydir", fs.ModeDir | 0755},
		{"link.txt", fs.ModeSymlink | 0777},
		{"script.sh", 0755},
		{"main.go", 0644},
	}
	for _, tc := range cases {
		if got := clr.Apply(tc.name, tc.mode); got != tc.name {
			t.Errorf("color never: Apply(%q) = %q, want unchanged", tc.name, got)
		}
	}
}

func TestApply_always(t *testing.T) {
	var buf bytes.Buffer
	clr := color.NewFromLSColors("always", &buf, testLSColors)

	cases := []struct {
		name    string
		mode    fs.FileMode
		wantSGR string
	}{
		{"mydir", fs.ModeDir | 0755, "01;34"},
		{"link.txt", fs.ModeSymlink | 0777, "01;36"},
		{"script.sh", 0755, "01;32"},   // executable bit set
		{"main.go", 0644, "00;36"},     // *.go mapping
		{"plain.txt", 0644, "00;37"},   // no *.txt → falls back to fi
		{"README", 0644, "00;37"},      // no extension → falls back to fi
	}
	for _, tc := range cases {
		got := clr.Apply(tc.name, tc.mode)
		want := "\033[" + tc.wantSGR + "m" + tc.name + "\033[0m"
		if got != want {
			t.Errorf("Apply(%q, %04o) = %q, want %q", tc.name, tc.mode, got, want)
		}
	}
}

func TestApply_autoNonTTY(t *testing.T) {
	// bytes.Buffer is not a TTY — auto should disable color.
	var buf bytes.Buffer
	clr := color.NewFromLSColors("auto", &buf, testLSColors)

	if got := clr.Apply("mydir", fs.ModeDir|0755); got != "mydir" {
		t.Errorf("color auto (non-TTY): Apply should be unchanged, got %q", got)
	}
}

func TestApply_emptyLSColors(t *testing.T) {
	// No LS_COLORS mappings — Apply should return name unchanged even with always.
	var buf bytes.Buffer
	clr := color.NewFromLSColors("always", &buf, "")

	if got := clr.Apply("mydir", fs.ModeDir|0755); got != "mydir" {
		t.Errorf("empty lsColors: Apply should be unchanged, got %q", got)
	}
}
