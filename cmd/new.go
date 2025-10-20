package cmd

import (
	"fmt"
	"quartz/dagbench.io/benchmark/recipe"
	"quartz/dagbench.io/config"

	"github.com/spf13/cobra"
)

var (
	format     string
	initRecipe string
	outDir     string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new benchmark configuration",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cmd.SilenceUsage = true
		configOpts := []config.OptFunc{}

		if autoInit {
			if moduleName == "" {
				moduleName = name
			}

			configOpts = append(configOpts, config.WithInit(moduleName, sdk, templateDir))
		}

		if module != "" {
			configOpts = append(configOpts, config.WithModule(module))
		}

		if workdir != "" {
			configOpts = append(configOpts, config.WithWorkdir(workdir))
		}

		if useCloud {
			configOpts = append(configOpts, config.EnableCloud())
		}

		configuration := config.New(name, daggerBin, 10, configOpts...)

		if initRecipe != "" {
			configuration.Commands = recipe.Get(recipe.Recipe(initRecipe))
		}

		ctx := cmd.Context()
		if err := configuration.Verify(ctx); err != nil {
			return fmt.Errorf("failed to verify configuration: %w", err)
		}

		if autoInit {
			configuration.Workdir = ""
		}

		path, err := configuration.Save(outDir, config.Format(format))
		if err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Configuration saved to %s\n", path)

		return nil
	},
}

func init() {
	// Common flag
	newCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the benchmark")
	newCmd.Flags().StringVarP(&daggerBin, "dagger-bin", "d", "dagger", "Dagger binary to use")
	newCmd.Flags().StringVarP(&outDir, "output", "o", ".", "Output directory for the configuration")
	newCmd.Flags().StringVar(&format, "format", "json", "Output format for the configuration")
	newCmd.Flags().
		StringVarP(&initRecipe, "recipe", "r", "", "Recipe to use for the benchmark (available recipes: sdk)")

	// Run flag
	newCmd.Flags().BoolVar(&useCloud, "use-cloud", false, "If enable, --cloud will be set")

	// Init flag
	newCmd.Flags().
		BoolVar(&autoInit, "auto-init", false, "Automatically init the module using provided flags")
	newCmd.Flags().StringVar(&moduleName, "module-name", "", "Name of the module to init")
	newCmd.Flags().StringVar(&workdir, "workdir", "", "Working directory for the benchmark")
	newCmd.Flags().StringVar(&sdk, "sdk", "", "Language to use for benchmark")
	newCmd.Flags().
		StringVar(&templateDir, "template-dir", "", "Template directory for the benchmark")

	// Module flag
	newCmd.Flags().StringVarP(&module, "module", "m", "", "Module to use for the benchmark")

	if err := newCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
}
