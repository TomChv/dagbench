package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/config"
	"quartz/dagbenchmark.io/sdk"

	"github.com/spf13/cobra"
)

var language string
var saveConfigDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "benchmark dagger init",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.NewConfig(language)
		if err != nil {
			return fmt.Errorf("failed to create configuration: %w", err)
		}

		fmt.Println(config)

		sdk, err := sdk.NewSDK(config)
		if err != nil {
			return fmt.Errorf("failed to create sdk: %w", err)
		}

		if err := sdk.PruneCache(); err != nil {
			return fmt.Errorf("failed to prune cache: %w", err)
		}

		execReport, err := sdk.Init().Exec()
		if err != nil {
			return fmt.Errorf("failed to execute setup command: %w", err)
		}

		if saveConfigDir != "" {
			if err := config.Save(saveConfigDir); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		if saveReportDir != "" {
			if err := execReport.SaveAsCSVAt(saveReportDir); err != nil {
				return fmt.Errorf("failed to save report: %w", err)
			}
		}

		fmt.Println(execReport)

		if saveOutputDir != "" {
			if err := execReport.SaveOutputAt(saveOutputDir); err != nil {
				return fmt.Errorf("failed to save output: %w", err)
			}
		}

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&language, "language", "l", "go", "Language to use for benchmark")
	initCmd.Flags().StringVar(&saveConfigDir, "save-config-dir", "", "Save config at path")
}
