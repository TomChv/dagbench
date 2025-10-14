package cmd

import (
	"fmt"
	"quartz/dagbench.io/benchmark"
	"quartz/dagbench.io/plots"
	"strings"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"

	"github.com/spf13/cobra"
)

var (
	useCmdDuration bool
	outputFile     string
)

var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Create a visualization for benchmark results",
}

var plotSingleBarCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "single-bar [file]",
	Short: "Create a bar visualization for a single benchmark file",
	RunE: func(_ *cobra.Command, args []string) error {
		file := args[0]
		opts := benchmark.ParserOpts{
			UseCmdDuration: useCmdDuration,
		}

		benchmarkFile, err := benchmark.ParseFile(file, opts)
		if err != nil {
			return err
		}

		aggregation := benchmarkFile.Aggregate()
		if name == "" {
			name = benchmarkFile.Name()
		}

		benchPlot, err := plots.SingleBar(aggregation, name)
		if err != nil {
			return err
		}

		plotfile := strings.TrimSuffix(file, ".txt") + ".svg"

		if err := benchPlot.Save(1.5*font.Length(len(aggregation))*vg.Inch, 6*vg.Inch, plotfile); err != nil {
			return fmt.Errorf("failed to save plot: %w", err)
		}

		fmt.Printf("Plot saved to %s\n", plotfile)

		return nil
	},
}

var plotMultiBarCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "multi-bar [files...]",
	Short: "Create a multi bar visualization for multiple benchmark files",
	Long: `Create a multi bar vizualisation from multiple benchmark files.

You can use it to compare the same module but with different setup.
For instance, comparing a TypeScript module but with Node, Bun or Deno runtime.`,
	RunE: func(_ *cobra.Command, args []string) error {
		opts := benchmark.ParserOpts{
			UseCmdDuration: useCmdDuration,
		}

		dataset := make(plots.DataSet, len(args))
		for _, file := range args {
			benchmarkFile, err := benchmark.ParseFile(file, opts)
			if err != nil {
				return err
			}

			aggregation := benchmarkFile.Aggregate()
			dataset[benchmarkFile.Name()] = aggregation
		}

		benchPlot, err := plots.MultiBar(dataset, name)
		if err != nil {
			return err
		}

		if err := plots.SaveWithOutsideLegend(
			benchPlot,
			1.3*font.Length(len(dataset))*vg.Inch,
			6*vg.Inch,
			1*vg.Inch,
			outputFile,
		); err != nil {
			return fmt.Errorf("failed to save plot: %w", err)
		}

		fmt.Printf("Plot saved to %s\n", outputFile)

		return nil
	},
}

var plotVersionsCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "versions-diff [files...]",
	Short: "Create a linear visualization for multiple benchmark files over different version",
	Long: `Create a linear vizualisation from multiple benchmark files.

You can use it to compare the same module over time to see the evolution of the performance.`,
	RunE: func(_ *cobra.Command, args []string) error {
		opts := benchmark.ParserOpts{
			UseCmdDuration: useCmdDuration,
		}

		dataset := make(plots.DataSet, len(args))
		for _, file := range args {
			benchmarkFile, err := benchmark.ParseFile(file, opts)
			if err != nil {
				return err
			}

			version := benchmarkFile.Version()
			aggregation := benchmarkFile.Aggregate()

			for name, value := range aggregation {
				if dataset[name] == nil {
					dataset[name] = make(plots.Data)
				}

				dataset[name][version] = value
			}
		}

		benchPlot, err := plots.MultiLinearByVersion(dataset, name)
		if err != nil {
			return err
		}

		if err := benchPlot.Save(3*font.Length(len(dataset))*vg.Inch, 6*vg.Inch, outputFile); err != nil {
			return fmt.Errorf("failed to save plot: %w", err)
		}

		fmt.Printf("Plot saved to %s\n", outputFile)

		return nil
	},
}

func init() {
	plotCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Name of the benchmark")
	plotCmd.PersistentFlags().
		BoolVar(&useCmdDuration, "use-cmd-duration", false, "Use cmdDuration entries (will ignore other entries)")

	plotMultiBarCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output of the benchmark")
	_ = plotMultiBarCmd.MarkFlagRequired("name")

	plotVersionsCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output of the benchmark")
	_ = plotVersionsCmd.MarkFlagRequired("name")

	// Subcommand
	plotCmd.AddCommand(plotSingleBarCmd)
	plotCmd.AddCommand(plotMultiBarCmd)
	plotCmd.AddCommand(plotVersionsCmd)
}
