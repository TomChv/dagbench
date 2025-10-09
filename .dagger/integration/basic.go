package integration

import (
	"context"

	"dagger/dagbench/internal/dagger"
	"golang.org/x/sync/errgroup"
)

func TestBasic(ctx context.Context, cli CLI) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testBasicRunFromCLI(gctx, cli) })
	eg.Go(func() error { return testBasicRunFromConfigFile(gctx, cli) })
	eg.Go(func() error { return testBasicRunFromExistingModule(gctx, cli) })

	return eg.Wait()
}

func testBasicRunFromCLI(ctx context.Context, cli CLI) error {
	ctx, span := Tracer().Start(ctx, "test run from CLI flag")
	defer span.End()

	_, err := cli.
		Container().
		WithExec(
			[]string{
				"run",
				"--name", "test", "--auto-init",
				"--sdk", "go", "--iteration", "2",
				"--command", "call container-echo --string-arg=hello",
				"--span", "containerEcho",
				"-o", "test.txt",
			},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.txt").
		Contents(ctx)

	return err
}

func testBasicRunFromConfigFile(ctx context.Context, cli CLI) error {
	ctx, span := Tracer().Start(ctx, "test run from config file")
	defer span.End()

	configFile := cli.
		Container().
		WithExec(
			[]string{"new", "--name", "test", "--sdk", "go", "--recipe", "sdk", "--auto-init"},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.json")

	_, err := cli.
		Container().
		WithFile("test.json", configFile).
		WithExec(
			[]string{"run", "--config", "test.json", "-i", "2", "-o", "test.txt"},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.txt").
		Contents(ctx)

	return err
}

func testBasicRunFromExistingModule(ctx context.Context, cli CLI) error {
	ctx, span := Tracer().Start(ctx, "test run from existing module")
	defer span.End()

	_, err := cli.
		Container().
		WithWorkdir("/test-module").
		WithExec([]string{"dagger", "init", "--sdk=go", "--name=test", "--source=."}).
		WithExec(
			[]string{
				"run", "-m", ".", "-i", "2", "--name", "test",
				"--command", "call container-echo --string-arg=hello",
				"--span", "containerEcho",
				"-o", "test.txt",
			},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.txt").
		Contents(ctx)

	return err
}
