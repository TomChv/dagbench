package main

import (
	"context"

	"dagger/dagbench-test/internal/dagger"
	"golang.org/x/sync/errgroup"
)

func testBasic(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testBasicRunFromCtr(gctx) })
	eg.Go(func() error { return testBasicRunFromConfigFile(gctx) })
	eg.Go(func() error { return testBasicRunFromExistingModule(gctx) })
	eg.Go(func() error { return testBasicRunWithTemplate(gctx) })
	eg.Go(func() error { return testBasicRunWithWorkdirAndCleanup(gctx) })

	return eg.Wait()
}

func testBasicRunFromCtr(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run from ctr flag")
	defer span.End()

	ctr := getTestCtr("basic-run-from-ctr")

	_, err := ctr.
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

func testBasicRunFromConfigFile(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run from config file")
	defer span.End()

	ctr := getTestCtr("basic-run-from-config-file")

	configFile := ctr.
		WithExec(
			[]string{"new", "--name", "test", "--sdk", "go", "--recipe", "sdk", "--auto-init"},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.json")

	_, err := ctr.
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

func testBasicRunFromExistingModule(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run from existing module")
	defer span.End()

	ctr := getTestCtr("basic-run-from-existing-module")

	_, err := ctr.
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

func testBasicRunWithTemplate(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run with template")
	defer span.End()

	ctr := getTestCtr("basic-run-with-template")

	_, err := ctr.
		WithDirectory(
			"/template",
			dag.Directory().WithNewFile("main.go",
				`package main

type Test struct{}

func (t *Test) Foo() string {
  return "Foo"
}`),
		).
		WithExec(
			[]string{
				"run",
				"--name", "test", "--auto-init",
				"--sdk", "go", "--template-dir", "/template",
				"--iteration", "2",
				"--command", "call foo",
				"--span", "foo",
				"-o", "test.txt",
			},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.txt").
		Contents(ctx)

	return err
}

func testBasicRunWithWorkdirAndCleanup(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run with workdir and auto cleanup")
	defer span.End()

	ctr := getTestCtr("basic-run-with-workdir-and-cleanup")

	_, err := ctr.
		WithDirectory("/tmp/workdir", dag.Directory()).
		WithExec(
			[]string{
				"run", "-i", "2", "--name", "test", "--workdir", "/tmp/workdir",
				"--command", "init --sdk=go --name=test --source=.",
				"--span", "generatedContextDirectory",
				"-o", "test.txt",
			},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			}).
		File("test.txt").
		Contents(ctx)

	return err
}
