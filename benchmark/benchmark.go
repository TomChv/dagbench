package benchmark

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"quartz/dagbench.io/config"
	"quartz/dagbench.io/dagger"
	"quartz/dagbench.io/hook"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/otiai10/copy"
)

// Run a benchmark using the given configuration.
//
// If a call fails, it will save the error report in
// `/tmp/dagbench-report/<name>/error-<command>-<iteration>.txt`
// and return an error.
// If debug is enabled, it will also save the output of the command
// in `/tmp/dagbench-report/<name>/output-<command>-<iteration>.txt`.
//
//nolint:funlen
func Run(ctx context.Context, conf *config.Config) (*Result, error) {
	if conf.DoAutoInit() {
		if err := autoInit(ctx, conf); err != nil {
			return nil, fmt.Errorf("failed to auto init: %w", err)
		}
	}

	daggerCLI := dagger.NewCLI(conf)
	result := make(Result, conf.Iteration)
	for i := range conf.Iteration {
		result[i] = make(map[string]time.Duration)
	}

	for i := range conf.Iteration {
		for _, cmd := range conf.Commands {
			fmt.Printf("Running command %v (iteration=%d)", cmd.Args, i)

			if err := daggerCLI.PruneCache(ctx); err != nil {
				return nil, fmt.Errorf("failed to prune cache: %w", err)
			}

			spanHook := hook.NewSpanMarker(cmd.SpanNames)
			outputHook := hook.NewCaptureOutput()

			cmdStart := time.Now()
			err := daggerCLI.Exec(ctx, cmd.Args, spanHook, outputHook)
			cmdDuration := time.Since(cmdStart)
			if err != nil {
				errorReportPath := filepath.Join(
					conf.DebugDir(),
					fmt.Sprintf("error-%s-%d.txt", cmd.Args[0], i),
				)

				return nil, trySaveErrReport(
					errorReportPath,
					outputHook.Output(),
					fmt.Errorf("failed to execute dagger: %w", err),
				)
			}

			if conf.IsDebug() {
				outputPath := filepath.Join(
					conf.DebugDir(),
					fmt.Sprintf("output-%s-%d.txt", cmd.Args[0], i),
				)
				_ = os.WriteFile(outputPath, []byte(outputHook.Output()), 0o600)
			}

			var inlineRes strings.Builder

			_, _ = fmt.Fprintf(&inlineRes, ": ")
			for spanName, duration := range spanHook.Values() {
				resultName := fmt.Sprintf(
					"%s/%s",
					strcase.ToLowerCamel(cmd.Args[0]),
					strcase.ToLowerCamel(spanName),
				)
				_, _ = fmt.Fprintf(
					&inlineRes,
					"[%s: %.1fs] ",
					strcase.ToLowerCamel(spanName),
					duration.Seconds(),
				)

				result[i][resultName] = duration
			}

			// Also add the command duration
			_, _ = fmt.Fprintf(&inlineRes, "[%s: %.1fs] ", "cmdDuration", cmdDuration.Seconds())
			result[i][fmt.Sprintf("%s/%s", strcase.ToLowerCamel(cmd.Args[0]), "cmdDuration")] = cmdDuration

			fmt.Printf("%s\n", inlineRes.String())
		}

		if conf.CleanUpAfterIteration() {
			if err := os.RemoveAll(conf.Workdir); err != nil {
				return nil, fmt.Errorf("failed to clean up workdir: %w", err)
			}

			if err := os.MkdirAll(conf.Workdir, 0o750); err != nil {
				return nil, fmt.Errorf("failed to create workdir: %w", err)
			}
		}
	}

	return &result, nil
}

func autoInit(ctx context.Context, conf *config.Config) error {
	daggerCLI := dagger.NewCLI(conf)

	fmt.Printf(
		"Auto initializing module at (sdk=%s ; workdir=%s)\n\n",
		conf.Init.SDK, conf.Workdir,
	)

	if err := daggerCLI.PruneCache(ctx); err != nil {
		return fmt.Errorf("failed to prune cache: %w", err)
	}

	outputHook := hook.NewCaptureOutput()
	if err := daggerCLI.Init(ctx, conf.Init.Name, conf.Init.SDK, outputHook); err != nil {
		errorReportPath := filepath.Join(conf.DebugDir(), "error-auto-init")

		return trySaveErrReport(
			errorReportPath,
			outputHook.Output(),
			fmt.Errorf("failed to init dagger: %w", err),
		)
	}

	if conf.Init.TemplateDir != "" {
		if err := copy.Copy(conf.Init.TemplateDir, conf.Workdir); err != nil {
			return fmt.Errorf("failed to copy template dir: %w", err)
		}
	}

	return nil
}

func trySaveErrReport(path string, output string, err error) error {
	if rerr := os.WriteFile(path, []byte(output), 0o600); rerr != nil {
		return err
	}

	return fmt.Errorf(
		"%w\nwriting error report to %s",
		err,
		path,
	)
}
