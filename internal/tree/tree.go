package tree

// Options controls how the tree walk behaves.
type Options struct {
	MaxDepth int
	All      bool
	DirsOnly bool
	Color    string
}

// Walk traverses root and renders the directory tree to stdout.
func Walk(root string, opts Options) error {
	return nil
}
