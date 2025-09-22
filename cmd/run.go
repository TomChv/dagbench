package cmd

import (
	"fmt"
	"quartz/dagbenchmark.io/config"
	execwrapper "quartz/dagbenchmark.io/exec-wrapper"
	"quartz/dagbenchmark.io/report"
	"quartz/dagbenchmark.io/sdk"
)

// Wrap any execution in a run to keep a consistent behaviour accross commands
func run(config *config.Config, sdk sdk.SDK, exWrapperBuilder func() *execwrapper.ExecWrapper) (err error) {
	fmt.Println("********* CONFIG **********")
	fmt.Println(config)

	fmt.Println("\n********* EXECUTION **********")

	var reports []*report.Report
	for range runs {
		if err := sdk.PruneCache(); err != nil {
			return fmt.Errorf("failed to prune cache: %w", err)
		}

		runReport, err := exWrapperBuilder().Exec()
		if err != nil {
			return err
		}

		reports = append(reports, runReport)
	}

	var result *report.Report
	if runs == 1 {
		result = reports[0]
	} else {
		fmt.Printf("Merging %d reports\n", len(reports))

		result, err = report.Merge(reports...)
		if err != nil {
			return fmt.Errorf("failed to merge reports: %w", err)
		}
	}

	fmt.Println("\n********* REPORTING **********")
	fmt.Println(result)

	if saveReportDir != "" {
		if err := result.SaveAsCSVAt(saveReportDir); err != nil {
			return err
		}
	}

	if saveOutputDir != "" && runs == 1 {
		if err := result.SaveOutputAt(saveOutputDir); err != nil {
			return err
		}
	}

	return nil
}
