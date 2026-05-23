package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type flags struct {
	level    int
	all      bool
	dirsOnly bool
	color    string
}

var f flags

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "baum [path]",
		Short:   "A rich terminal tree viewer",
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) == 1 {
				root = args[0]
			}
			fmt.Fprintf(cmd.OutOrStdout(), "baum — not yet implemented (root: %s)\n", root)
			return nil
		},
	}

	cmd.Flags().IntVarP(&f.level, "level", "L", 0, "maximum depth (0 = unlimited)")
	cmd.Flags().BoolVarP(&f.all, "all", "a", false, "include hidden files")
	cmd.Flags().BoolVarP(&f.dirsOnly, "dirs-only", "d", false, "list directories only")
	cmd.Flags().StringVar(&f.color, "color", "auto", "color output: always, auto, never")

	return cmd
}

func Execute(version string) {
	if err := newRootCmd(version).Execute(); err != nil {
		os.Exit(1)
	}
}
