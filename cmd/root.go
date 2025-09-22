package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var ErrMissingConfigFile = errors.New("missing config file")

var saveReportDir string
var saveOutputDir string
var configFile string
var runs int

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
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "Config file")
	rootCmd.PersistentFlags().IntVarP(&runs, "runs", "r", 1, "Number of runs (no stderr can be output in that case)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(developCmd)
	rootCmd.AddCommand(functionsCmd)
	rootCmd.AddCommand(callCmd)
}
