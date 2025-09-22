package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/config"
	execwrapper "quartz/dagbenchmark.io/exec-wrapper"
	"quartz/dagbenchmark.io/sdk"

	"github.com/spf13/cobra"
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "benchmark the complete flow of a dagger module",
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

		args, err = parseCallArgs(args)
		if err != nil {
			return err
		}

		err = runFlow(config, sdk, []func() *execwrapper.ExecWrapper{
			sdk.Develop,
			sdk.Functions,
			func() *execwrapper.ExecWrapper {
				return sdk.Call(args)
			},
		})
		if err != nil {
			return fmt.Errorf("failed to run flow: %w", err)
		}

		return nil
	},
}
