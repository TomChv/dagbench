package config

type ConfigOptFunc func(config *Config)

// WithWorkdir sets the workdir for the config
func WithWorkdir(workdir string) ConfigOptFunc {
	return func(config *Config) {
		config.Workdir = workdir
	}
}

func WithTemplateDir(templateDir string) ConfigOptFunc {
	return func(config *Config) {
		config.TemplateDir = templateDir
	}
}

func WithCommands(commands map[string][]string) ConfigOptFunc {
	return func(config *Config) {
		config.Commands = commands
	}
}

func DisableInit() ConfigOptFunc {
	return func(config *Config) {
		config.RunInit = false
	}
}

func EnableCloud() ConfigOptFunc {
	return func(config *Config) {
		config.Cloud = true
	}
}
