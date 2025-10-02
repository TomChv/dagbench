package daggerexec

import (
	"fmt"
	"os"
	"quartz/dagbench.io/config"
	"quartz/dagbench.io/daggerexec/hook"
	"quartz/dagbench.io/report"
)

func Run(c *config.Config, commandName string, args ...string) (*report.Report, error) {
	if err := pruneCache(c); err != nil {
		return nil, fmt.Errorf("failed to prune cache: %w", err)
	}

	switch commandName {
	case "init":
		return initModule(c)
	case "functions":
		return functions(c)
	case "develop":
		return develop(c)
	default:
		return call(c, args)
	}
}

func initModule(c *config.Config) (*report.Report, error) {
	if err := os.MkdirAll(c.Workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workdir for module: %w", err)
	}

	if c.RunInit {
		report := report.New("init")
		outputHook := hook.NewCaptureOutput()

		if err := execDagger(
			c,
			[]string{"init", "--sdk", c.SDK, "--source=.", "--name", c.Name},
		); err != nil {
			report = report.WithOutput(outputHook.Output())

			return report, fmt.Errorf("failed to init module: %w", err)
		}

		return nil, nil
	}

	panic("not implemented yet")
}

func functions(c *config.Config) (*report.Report, error) {
	report := report.New("functions")

	spanHook := hook.NewSpanMarker([]string{"load module"})
	outputHook := hook.NewCaptureOutput()

	if err := execDagger(
		c,
		[]string{"functions"},
		spanHook,
		outputHook,
	); err != nil {
		report = report.WithOutput(outputHook.Output())

		return report, fmt.Errorf("failed to run functions: %w", err)
	}

	report = report.WithValues(spanHook.Values())

	return report, nil
}

func develop(c *config.Config) (*report.Report, error) {
	report := report.New("develop")

	spanHook := hook.NewSpanMarker([]string{"develop"})
	outputHook := hook.NewCaptureOutput()

	if err := execDagger(
		c,
		[]string{"develop"},
		spanHook,
		outputHook,
	); err != nil {
		report = report.WithOutput(outputHook.Output())

		return report, fmt.Errorf("failed to run develop: %w", err)
	}

	report = report.WithValues(spanHook.Values())

	return report, nil
}

func call(c *config.Config, args []string) (*report.Report, error) {
	report := report.New("call")

	spanHook := hook.NewSpanMarker([]string{
		"load module",
		convertFunctionNameToTraceMarker(args[1]),
	})
	outputHook := hook.NewCaptureOutput()

	if err := execDagger(
		c,
		args,
		spanHook,
		outputHook,
	); err != nil {
		report = report.WithOutput(outputHook.Output())

		return report, fmt.Errorf("failed to run call: %w", err)
	}

	report = report.WithValues(spanHook.Values())

	return report, nil
}
