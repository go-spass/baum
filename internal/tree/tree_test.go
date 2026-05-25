package tree_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-spass/baum/internal/tree"
)

func TestWalk(t *testing.T) {
	cases := []struct {
		name string
		opts tree.Options
	}{
		{"basic", tree.Options{}},
		{"empty", tree.Options{}},
		{"hidden", tree.Options{All: true}},
		{"dirs-only", tree.Options{DirsOnly: true}},
		{"depth-limit", tree.Options{MaxDepth: 2}},
		{"symlink", tree.Options{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			scenario, err := filepath.Abs(filepath.Join("testdata", tc.name))
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join(scenario, "golden.txt")
			root := filepath.Join(scenario, "root")

			t.Chdir(root)

			var buf bytes.Buffer
			if err := tree.Walk(&buf, ".", tc.opts); err != nil {
				t.Fatalf("Walk: %v", err)
			}
			got := buf.Bytes()

			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.WriteFile(golden, got, 0644); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("output mismatch\nwant:\n%s\ngot:\n%s", want, got)
			}
		})
	}
}
