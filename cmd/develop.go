package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/config"
	"quartz/dagbenchmark.io/sdk"

	"github.com/spf13/cobra"
)

var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "benchmark dagger develop",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			return ErrMissingConfigFile
		}

		config, err := config.NewConfigFromFile(configFile)
		if err != nil {
			return err
		}

		sdk, err := sdk.NewSDK(config)
		if err != nil {
			return err
		}

		if err := sdk.PruneCache(); err != nil {
			return fmt.Errorf("failed to prune cache: %w", err)
		}

		execReport, err := sdk.Develop().Exec()
		if err != nil {
			return err
		}

		if saveReportDir != "" {
			if err := execReport.SaveAsCSVAt(saveReportDir); err != nil {
				return err
			}
		}

		fmt.Println(execReport)

		if saveOutputDir != "" {
			if err := execReport.SaveOutputAt(saveOutputDir); err != nil {
				return err
			}
		}

		return nil
	},
}
