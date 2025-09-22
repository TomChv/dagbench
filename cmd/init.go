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

		sdk, err := sdk.NewSDK(config)
		if err != nil {
			return fmt.Errorf("failed to create sdk: %w", err)
		}

		if err := run(config, sdk, sdk.Init); err != nil {
			return fmt.Errorf("failed to init module: %w", err)
		}

		if saveConfigDir != "" {
			if err := config.Save(saveConfigDir); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&language, "language", "l", "go", "Language to use for benchmark")
	initCmd.Flags().StringVar(&saveConfigDir, "save-config-dir", "", "Save config at path")
}
