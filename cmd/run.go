package cmd

import (
	"fmt"
	"os"
	"quartz/dagbench.io/benchmark"
	"quartz/dagbench.io/config"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Config file
	configFile string

	// Common flags
	name       string
	outputPath string
	daggerBin  string

	// Run flags
	iteration int
	command   string
	spanName  string
	useCloud  bool

	// Init flags
	autoInit    bool
	moduleName  string
	sdk         string
	workdir     string
	templateDir string

	// Module flags
	module string

	// Debug mode
	debug bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a dagger benchmark",
	RunE: func(cmd *cobra.Command, args []string) error {
		var configuration *config.Config
		if configFile != "" {
			conf, err := configFromFile()
			if err != nil {
				return fmt.Errorf("failed to load config from file: %w", err)
			}

			configuration = conf
		} else {
			configuration = configFromFlag()
		}

		if err := configuration.Verify(); err != nil {
			return fmt.Errorf("failed to verify configuration: %w", err)
		}

		fmt.Printf("Running benchmark %s with %d iterations\n\n", configuration.Name, configuration.Iteration)

		result, err := benchmark.Run(configuration)
		if err != nil {
			return fmt.Errorf("failed to run benchmark: %w", err)
		}

		fmtResult := result.GoBenchFormat(configuration.Name, configuration.Metadatas())
		fmt.Printf("\nBenchmark result:\n%s\n", fmtResult)

		if err := os.WriteFile(outputPath, []byte(fmtResult), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Benchmark saved to %s\n", outputPath)

		return nil
	},
}

func configFromFile() (*config.Config, error) {
	conf, err := config.NewFromFile(configFile)
	if err != nil {
		return nil, err
	}

	if name != "" {
		conf.Name = name
	}

	if conf.Init != nil {
		if moduleName == "" {
			moduleName = name
		}

		if templateDir != "" {
			conf.Init.TemplateDir = templateDir
		}

		if moduleName != "" {
			conf.Init.Name = moduleName
		}

		if sdk != "" {
			conf.Init.SDK = sdk
		}
	}

	if useCloud {
		conf.Cloud = true
	}

	if module != "" {
		conf.Module = module
	}

	if daggerBin != "" {
		conf.BinPath = daggerBin
	}

	if iteration != 0 {
		conf.Iteration = iteration
	}

	if debug {
		conf.EnableDebug()
	}

	return conf, nil
}

func configFromFlag() *config.Config {
	configOpts := []config.ConfigOptFunc{
		config.WithCommand([]string{spanName}, strings.Split(command, " ")),
	}

	if autoInit {
		if moduleName == "" {
			moduleName = name
		}

		configOpts = append(configOpts, config.WithInit(moduleName, sdk, templateDir))
	}

	if module != "" {
		configOpts = append(configOpts, config.WithModule(module))
	}

	if useCloud {
		configOpts = append(configOpts, config.EnableCloud())
	}

	if debug {
		configOpts = append(configOpts, config.EnableDebug())
	}

	return config.New(name, daggerBin, iteration, configOpts...)
}

func init() {
	// Config file
	runCmd.Flags().StringVar(&configFile, "config", "", "Config file to use")

	// Common flag
	runCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the benchmark")
	runCmd.Flags().StringVarP(&daggerBin, "dagger-bin", "d", "dagger", "Dagger binary to use")
	runCmd.Flags().StringVarP(&outputPath, "output", "o", "out.txt", "Output file for the report")

	// Run flag
	runCmd.Flags().IntVarP(&iteration, "iteration", "i", 10, "Number of iterations to run")
	runCmd.Flags().StringVarP(&spanName, "span", "s", "", "Span name to record")
	runCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run")
	runCmd.Flags().BoolVar(&useCloud, "use-cloud", false, "If enable, --cloud will be set")

	// Init flag
	runCmd.Flags().BoolVar(&autoInit, "auto-init", false, "Automatically init the module using provided flags")
	runCmd.Flags().StringVar(&moduleName, "module-name", "", "Name of the module to init")
	runCmd.Flags().StringVar(&workdir, "workdir", "", "Working directory for the benchmark")
	runCmd.Flags().StringVar(&sdk, "sdk", "", "Language to use for benchmark")
	runCmd.Flags().StringVar(&templateDir, "template-dir", "", "Template directory for the benchmark")

	// Module flag
	runCmd.Flags().StringVarP(&module, "module", "m", "", "Module to use for the benchmark")

	// Debug mode
	runCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode")
}
