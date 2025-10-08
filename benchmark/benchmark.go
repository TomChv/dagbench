package benchmark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/otiai10/copy"

	"quartz/dagbench.io/config"
	"quartz/dagbench.io/dagger"
	"quartz/dagbench.io/hook"
)

func Run(conf *config.Config) (*Result, error) {
	daggerCLI := dagger.NewCLI(conf)

	if conf.DoAutoInit() {
		fmt.Printf("Auto initializing module at (sdk=%s ; workdir=%s)\n\n", conf.Init.SDK, conf.Workdir)

		if err := daggerCLI.PruneCache(); err != nil {
			return nil, fmt.Errorf("failed to prune cache: %w", err)
		}

		outputHook := hook.NewCaptureOutput()
		if err := daggerCLI.Init(conf.Init.Name, conf.Init.SDK, outputHook); err != nil {
			errorReportPath := filepath.Join(conf.DebugDir(), "error-auto-init")

			if rerr := os.WriteFile(errorReportPath, []byte(outputHook.Output()), 0o644); rerr != nil {
				return nil, fmt.Errorf("failed to init dagger: %w", err)
			}

			return nil, fmt.Errorf("failed to init dagger: %w\nwriting error report to %s", err, errorReportPath)
		}

		if conf.Init.TemplateDir != "" {
			if err := copy.Copy(conf.Init.TemplateDir, conf.Workdir); err != nil {
				return nil, fmt.Errorf("failed to copy template dir: %w", err)
			}
		}
	}

	result := make(Result, conf.Iteration)
	for i := 0; i < conf.Iteration; i++ {
		result[i] = make(map[string]time.Duration)
	}

	for _, cmd := range conf.Commands {
		for i := 0; i < conf.Iteration; i++ {
			fmt.Printf("Running command %v (iteration=%d)", cmd.Args, i)

			if err := daggerCLI.PruneCache(); err != nil {
				return nil, fmt.Errorf("failed to prune cache: %w", err)
			}

			spanHook := hook.NewSpanMarker(cmd.SpanNames)
			outputHook := hook.NewCaptureOutput()

			cmdStart := time.Now()
			err := daggerCLI.Exec(cmd.Args, spanHook, outputHook)
			cmdDuration := time.Since(cmdStart)
			if err != nil {
				errorReportPath := filepath.Join(conf.DebugDir(), fmt.Sprintf("error-%s-%d.txt", cmd.Args[0], i))

				if rerr := os.WriteFile(errorReportPath, []byte(outputHook.Output()), 0o644); rerr != nil {
					return nil, fmt.Errorf("failed to execute dagger: %w", err)
				}

				return nil, fmt.Errorf("failed to execute dagger: %w\nwriting error report to %s", err, errorReportPath)
			}

			if conf.IsDebug() {
				outputPath := filepath.Join(conf.DebugDir(), fmt.Sprintf("output-%s-%d.txt", cmd.Args[0], i))
				_ = os.WriteFile(outputPath, []byte(outputHook.Output()), 0o644)
			}

			var inlineRes strings.Builder

			fmt.Fprintf(&inlineRes, ": ")
			for spanName, duration := range spanHook.Values() {
				resultName := fmt.Sprintf("%s/%s", strcase.ToLowerCamel(cmd.Args[0]), strcase.ToLowerCamel(spanName))
				fmt.Fprintf(&inlineRes, "[%s: %.1fs] ", strcase.ToLowerCamel(spanName), duration.Seconds())

				result[i][resultName] = duration
			}

			// Also add the command duration
			fmt.Fprintf(&inlineRes, "[%s: %.1fs] ", "cmdDuration", cmdDuration.Seconds())
			result[i][fmt.Sprintf("%s/%s", strcase.ToLowerCamel(cmd.Args[0]), "cmdDuration")] = cmdDuration

			fmt.Printf("%s\n", inlineRes.String())
		}
	}

	return &result, nil
}
