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

		if err := run(config, sdk, sdk.Develop); err != nil {
			return fmt.Errorf("failed to run develop: %w", err)
		}

		return nil
	},
}
