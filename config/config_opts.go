package config

// OptFunc is a function that modifies a Config.
type OptFunc func(config *Config)

// WithInit sets the init command of the config.
func WithInit(name, sdk, templateDir string) OptFunc {
	return func(config *Config) {
		config.Init = &Init{
			Name:        name,
			SDK:         sdk,
			TemplateDir: templateDir,
		}
	}
}

// WithModule sets the module of the config.
func WithModule(module string) OptFunc {
	return func(config *Config) {
		config.Module = module
	}
}

// EnableCloud enables the cloud mode of the config.
func EnableCloud() OptFunc {
	return func(config *Config) {
		config.Cloud = true
	}
}

// WithCommand sets the command of the config.
func WithCommand(spanNames []string, args []string) OptFunc {
	return func(config *Config) {
		config.Commands = append(config.Commands, &Command{
			SpanNames: spanNames,
			Args:      args,
		})
	}
}

// EnableDebug enables the debug mode of the config.
func EnableDebug() OptFunc {
	return func(config *Config) {
		config.debug = true
	}
}
