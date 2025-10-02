package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dagbench",
	Short: "Dagger benchmark CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(benchCmd)
}
