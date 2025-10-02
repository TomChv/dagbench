package cmd

import (
	"fmt"
	"os"
	"quartz/dagbench.io/config"

	"github.com/spf13/cobra"
)

var (
	sdk         string
	bin         string
	workdir     string
	templatedir string
	format      string
	disableInit bool
	useCloud    bool
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Args:  cobra.ExactArgs(1),
	Short: "initialiaze a new dagger benchmark configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := config.StringToFormat(format)
		if err != nil {
			return fmt.Errorf("failed to parse --format argument: %w", err)
		}

		// Create config
		configOpts := []config.ConfigOptFunc{
			config.WithWorkdir(workdir),
			config.WithTemplateDir(templatedir),
		}

		if disableInit {
			configOpts = append(configOpts, config.DisableInit())
		}

		if useCloud {
			configOpts = append(configOpts, config.EnableCloud())
		}

		name := args[0]
		conf, err := config.New(name, sdk, bin, configOpts...)
		if err != nil {
			return fmt.Errorf("failed to create configuration: %w", err)
		}

		// Save config
		workdir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get workdir: %w", err)
		}

		configPath, err := conf.Save(workdir, format)
		if err != nil {
			return err
		}

		fmt.Printf("Config saved at %s\n", configPath)

		return nil
	},
}

func init() {
	newCmd.Flags().StringVar(&sdk, "sdk", "", "Language to use for benchmark")
	newCmd.Flags().StringVar(&bin, "bin", "dagger", "Name of the dagger binary")
	newCmd.Flags().StringVar(&workdir, "workdir", "", "Working directory for the benchmark")
	newCmd.Flags().StringVar(&templatedir, "templatedir", "", "Template directory for the benchmark")
	newCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Format to save the config in")
	newCmd.Flags().BoolVar(&disableInit, "disable-init", false, "Disable the init command")
	newCmd.Flags().BoolVar(&useCloud, "use-cloud", false, "If enable, --cloud will be set")

	if err := newCmd.MarkFlagRequired("sdk"); err != nil {
		panic(err)
	}
}
