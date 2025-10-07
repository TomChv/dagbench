package config

type ConfigOptFunc func(config *Config)

func WithInit(name, sdk, templateDir string) ConfigOptFunc {
	return func(config *Config) {
		config.Init = &Init{
			Name:        name,
			SDK:         sdk,
			TemplateDir: templateDir,
		}
	}
}

func WithModule(module string) ConfigOptFunc {
	return func(config *Config) {
		config.Module = module
	}
}

func EnableCloud() ConfigOptFunc {
	return func(config *Config) {
		config.Cloud = true
	}
}

func WithCommand(spanNames []string, args []string) ConfigOptFunc {
	return func(config *Config) {
		config.Commands = append(config.Commands, &Command{
			SpanNames: spanNames,
			Args:      args,
		})
	}
}

func EnableDebug() ConfigOptFunc {
	return func(config *Config) {
		config.debug = true
	}
}
