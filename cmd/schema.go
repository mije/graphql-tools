package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mije/graphql-tools/pkg/schema/compare"
	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two schemas",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flag("in").Value.String() != "sdl" {
			return fmt.Errorf("unsupported input format")
		}
		if cmd.Flag("out").Value.String() != "txt" {
			return fmt.Errorf("unsupported output format")
		}

		x, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer x.Close()

		y, err := os.Open(args[1])
		if err != nil {
			return err
		}
		defer y.Close()

		res, err := compare.Schema(x, y)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "PATH\tSEVERITY\tTYPE\tDESCRIPTION\t")
		for _, c := range res.Changes() {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", c.Path, c.Severity.Level, c.Type, c.Message)
		}
		return w.Flush()
	},
}

func init() {
	compareCmd.Flags().StringP("in", "i", "sdl", "input format")
	compareCmd.Flags().StringP("out", "o", "txt", "output format")

	schemaCmd.AddCommand(compareCmd)
}
