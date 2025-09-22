package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var ErrMissingConfigFile = errors.New("missing config file")

var reportDir string
var outputDir string
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
	// TODO: improve flag separation
	rootCmd.PersistentFlags().StringVar(&reportDir, "report-dir", "", "Directory to save reports in")
	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Directory to save command output in")
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "Config file")
	rootCmd.PersistentFlags().IntVarP(&runs, "runs", "r", 1, "Number of runs (no stderr can be output in that case)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(developCmd)
	rootCmd.AddCommand(functionsCmd)
	rootCmd.AddCommand(callCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(flowCmd)
}
