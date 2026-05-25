package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	fmt.Fprintln(w, root)
	var dirs, files int
	if err := walkDir(w, root, "", 1, opts, &dirs, &files); err != nil {
		return err
	}
	fmt.Fprintln(w, "")
	writeSummary(w, dirs, files, opts.DirsOnly)
	return nil
}

func walkDir(w io.Writer, dir, prefix string, depth int, opts Options, dirs, files *int) error {
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

		if isSymlink {
			if target, err := os.Readlink(filepath.Join(dir, name)); err == nil {
				name = name + " -> " + target
			}
			*files++
		} else if entry.IsDir() {
			*dirs++
		} else {
			*files++
		}

		fmt.Fprintf(w, "%s%s%s\n", prefix, connector, name)

		if entry.IsDir() && (opts.MaxDepth == 0 || depth < opts.MaxDepth) {
			if err := walkDir(w, filepath.Join(dir, entry.Name()), childPrefix, depth+1, opts, dirs, files); err != nil {
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
