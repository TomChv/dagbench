package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/config"
	"quartz/dagbenchmark.io/sdk"

	"github.com/spf13/cobra"
)

var language string
var configDir string
var alias string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "benchmark dagger init",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := config.NewConfig(language, config.WithName(alias))
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

		if configDir != "" {
			if err := config.Save(configDir); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&language, "language", "l", "go", "Language to use for benchmark")
	initCmd.Flags().StringVar(&configDir, "config-dir", "", "Directory to save config file to")
	initCmd.Flags().StringVar(&alias, "alias", "", "Alias for the benchmark")
}
