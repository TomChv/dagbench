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

	if reportDir != "" {
		if err := result.SaveAsCSVAt(reportDir); err != nil {
			return err
		}
	}

	if outputDir != "" && runs == 1 {
		if err := result.SaveOutputAt(outputDir); err != nil {
			return err
		}
	}

	return nil
}

func runFlow(config *config.Config, sdk sdk.SDK, exWrapperBuilders []func() *execwrapper.ExecWrapper) (err error) {
	fmt.Println("********* CONFIG **********")
	fmt.Println(config)

	fmt.Println("\n********* EXECUTION **********")

	reportsMap := make(map[int][]*report.Report)
	for range runs {
		for i, exWrapperBuilder := range exWrapperBuilders {
			if err := sdk.PruneCache(); err != nil {
				return fmt.Errorf("failed to prune cache: %w", err)
			}

			runReport, err := exWrapperBuilder().Exec()
			if err != nil {
				return err
			}

			reportsMap[i] = append(reportsMap[i], runReport)
		}
	}

	var reports []*report.Report
	for _, reportsMap := range reportsMap {
		if len(reportsMap) == 1 {
			reports = append(reports, reportsMap[0])
			continue
		}

		fmt.Printf("Merging %d reports\n", len(reportsMap))

		result, err := report.Merge(reportsMap...)
		if err != nil {
			return fmt.Errorf("failed to merge reports: %w", err)
		}

		reports = append(reports, result)
	}

	fmt.Println("\n********* REPORTING **********")
	for _, report := range reports {
		fmt.Println(report)

		if reportDir != "" {
			if err := report.SaveAsCSVAt(reportDir); err != nil {
				return err
			}
		}
	}

	return nil
}
