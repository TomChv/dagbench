package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"quartz/dagbench.io/config"
	"quartz/dagbench.io/daggerexec"
	"quartz/dagbench.io/report"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var round int

var out string

var benchCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "bench <config-path>",
	Short: "run a dagger benchmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := args[0]

		config, err := config.NewFromFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if config.RunInit {
			// Create a unique workdir for the run if we handle the init
			config.Workdir = filepath.Join(config.Workdir, uuid.NewString())

			fmt.Printf("Initializing module\n")
			if _, err := daggerexec.Run(config, "init"); err != nil {
				return fmt.Errorf("failed to init module: %w", err)
			}
		} else {
			fmt.Printf("Skipping module initialization\n")
		}

		reports := []*report.Report{}
		for name, cmd := range config.Commands {
			for i := 0; i < round; i++ {
				fmt.Printf("Running %s (iteration=%d)\n", name, i)

				report, err := daggerexec.Run(config, name, cmd...)
				if err != nil {
					// If the command failed but some output was generated, save it for debugging.
					if !report.HasOutput() {
						return fmt.Errorf("failed to run command: %w", err)
					}

					fmt.Printf("Command %s failed, saving output at ./dagbench-error.txt\n", name)
					if rerr := report.SaveOutput("./dagbench-error.txt"); rerr != nil {
						return fmt.Errorf("failed to save error output: %w", rerr)
					}

					return fmt.Errorf("failed to run command: %w", err)
				}

				reports = append(reports, report)
			}
		}

		benchResult := report.FormatGoBenchReport(config, reports)
		fmt.Println(benchResult)

		if out != "" {
			if err := os.WriteFile(out, []byte(benchResult), 0644); err != nil {
				return fmt.Errorf("failed to write report to file: %w", err)
			}

			fmt.Printf("Benchmark saved at %s\n", out)
		}

		return nil
	},
}

func init() {
	benchCmd.Flags().IntVarP(&round, "round", "r", 1, "Number of rounds to run each commands")
	benchCmd.Flags().StringVarP(&out, "out", "o", "", "Output file for the report")
}
