package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "graphql-tools",
		Short: "GraphQL Tools",
	}
	schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Set of tools to process schemas.",
	}
)

func init() {
	rootCmd.AddCommand(schemaCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
