package cmd

import (
	"github.com/spf13/cobra"
)

var saveReportDir string
var saveOutputDir string

var rootCmd = &cobra.Command{
	Use:   "dagbenchmark",
	Short: "Dagger benchmark CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&saveReportDir, "save-report-dir", "", "Save report directory")
	rootCmd.PersistentFlags().StringVar(&saveOutputDir, "save-output-dir", "", "Save output directory")

	rootCmd.AddCommand(initCmd)
}
