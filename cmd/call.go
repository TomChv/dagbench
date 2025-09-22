package cmd

import (
	"errors"
	"fmt"
	"quartz/dagbenchmark.io/config"
	execwrapper "quartz/dagbenchmark.io/exec-wrapper"
	"quartz/dagbenchmark.io/sdk"
	"strings"

	"github.com/spf13/cobra"
)

var ErrMissingArgs = errors.New("missing args")

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "benchmark dagger call",
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

		if err := run(config, sdk, func() *execwrapper.ExecWrapper {
			return sdk.Call(args)
		}); err != nil {
			return fmt.Errorf("failed to call %s: %w", strings.Join(args, " "), err)
		}

		return nil
	},
}

func parseCallArgs(args []string) ([]string, error) {
	if len(args) == 0 {
		return nil, ErrMissingArgs
	}

	var cmdArgs []string
	for _, arg := range args {
		cmdArgs = append(cmdArgs, strings.Split(arg, " ")...)
	}

	return cmdArgs, nil
}
