package cmd

import (
	"errors"
	"fmt"
	"quartz/dagbenchmark.io/config"
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

		if err := sdk.PruneCache(); err != nil {
			return fmt.Errorf("failed to prune cache: %w", err)
		}

		args, err = parseCallArgs(args)
		if err != nil {
			return err
		}

		fmt.Printf("Executed command: %s\n", strings.Join(args, " "))

		execReport, err := sdk.Call(args).Exec()
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
