package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/report"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <report1> <report2>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a diff between 2 reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		report1, err := report.NewFromCSV(args[0])
		if err != nil {
			return fmt.Errorf("failed to load report at %s: %w", args[0], err)
		}

		report2, err := report.NewFromCSV(args[1])
		if err != nil {
			return fmt.Errorf("failed to load report at %s: %w", args[1], err)
		}

		diffReport := report1.Diff(report2)

		if saveReportDir != "" {
			if err := diffReport.SaveAsCSVAt(saveReportDir); err != nil {
				return err
			}
		}

		fmt.Println("\n********* DIFFERENCE **********")
		fmt.Println(diffReport)

		return nil
	},
}
