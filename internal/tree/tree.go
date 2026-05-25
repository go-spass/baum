package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-spass/baum/internal/color"
)

// Options controls how the tree walk behaves.
type Options struct {
	MaxDepth int
	All      bool
	DirsOnly bool
	Color    string
}

// Walk traverses root and renders the directory tree to w.
func Walk(w io.Writer, root string, opts Options) error {
	clr := color.New(opts.Color, w)
	fmt.Fprintln(w, clr.Apply(root, os.ModeDir))
	var dirs, files int
	if err := walkDir(w, root, "", 1, opts, clr, &dirs, &files); err != nil {
		return err
	}
	fmt.Fprintln(w, "")
	writeSummary(w, dirs, files, opts.DirsOnly)
	return nil
}

func walkDir(w io.Writer, dir, prefix string, depth int, opts Options, clr *color.Colorizer, dirs, files *int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var visible []os.DirEntry
	for _, e := range entries {
		if !opts.All && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if opts.DirsOnly && !e.IsDir() {
			continue
		}
		visible = append(visible, e)
	}

	for i, entry := range visible {
		isLast := i == len(visible)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}

		name := entry.Name()
		isSymlink := entry.Type()&os.ModeSymlink != 0

		var mode os.FileMode
		if info, err := entry.Info(); err == nil {
			mode = info.Mode()
		}

		if isSymlink {
			coloredName := clr.Apply(name, mode)
			if target, err := os.Readlink(filepath.Join(dir, name)); err == nil {
				name = coloredName + " -> " + target
			} else {
				name = coloredName
			}
			*files++
		} else if entry.IsDir() {
			*dirs++
			name = clr.Apply(name, mode)
		} else {
			*files++
			name = clr.Apply(name, mode)
		}

		fmt.Fprintf(w, "%s%s%s\n", prefix, connector, name)

		if entry.IsDir() && (opts.MaxDepth == 0 || depth < opts.MaxDepth) {
			if err := walkDir(w, filepath.Join(dir, entry.Name()), childPrefix, depth+1, opts, clr, dirs, files); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeSummary(w io.Writer, dirs, files int, dirsOnly bool) {
	dw := "directories"
	if dirs == 1 {
		dw = "directory"
	}
	if dirsOnly {
		fmt.Fprintf(w, "%d %s\n", dirs, dw)
		return
	}
	fw := "files"
	if files == 1 {
		fw = "file"
	}
	fmt.Fprintf(w, "%d %s, %d %s\n", dirs, dw, files, fw)
}
