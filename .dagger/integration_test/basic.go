package main

import (
	"context"

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

	_, err := getTestCLI("test-basic-run-from-ctr").
		Run().
		WithName("test").
		WithAutoInit().
		WithSDK("go").
		WithIteration(2).
		WithCommand([]string{"call", "container-echo", "--string-arg=hello"}).
		WithSpan("containerEcho").
		Exec().
		Contents(ctx)

	return err
}

func testBasicRunFromConfigFile(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run from config file")
	defer span.End()

	cli := getTestCLI("basic-run-from-config-file")

	configFile := cli.
		New().
		WithName("test").
		WithSDK("go").
		WithRecipe("sdk").
		WithAutoInit().
		Exec()

	_, err := cli.
		Run().
		WithConfigFile(configFile).
		WithIteration(2).
		Exec().
		Contents(ctx)

	return err
}

func testBasicRunFromExistingModule(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run from existing module")
	defer span.End()

	cli := getTestCLI("basic-run-from-existing-module")

	module := cli.
		Container().
		WithWorkdir("/test-module").
		WithExec([]string{"dagger", "init", "--sdk=go", "--name=test", "--source=."}).
		Directory("/test-module")

	_, err := cli.
		Run().
		WithModule("/test-module", module).
		WithIteration(2).
		WithName("test").
		WithCommand([]string{"call", "container-echo", "--string-arg=hello"}).
		WithSpan("contianerEcho").
		Exec().
		Contents(ctx)

	return err
}

func testBasicRunWithTemplate(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run with template")
	defer span.End()

	templateDir := dag.Directory().WithNewFile(
		"main.go",
		`package main

type Test struct{}

func (t *Test) Foo() string {
  return "Foo"
}`)

	_, err := getTestCLI("basic-run-with-template").
		Run().
		WithName("test").
		WithAutoInit().
		WithSDK("go").
		WithTemplateDir(templateDir).
		WithIteration(2).
		WithCommand([]string{"call", "foo"}).
		WithSpan("foo").
		Exec().
		Contents(ctx)

	return err
}

func testBasicRunWithWorkdirAndCleanup(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test run with workdir and auto cleanup")
	defer span.End()

	_, err := getTestCLI("basic-run-with-workdir-and-cleanup").
		Run().
		WithName("test").
		WithWorkdir("/tmp/workdir").
		WithIteration(2).
		WithCommand([]string{"init", "--sdk=go", "--name=test", "--source=."}).
		WithSpan("generatedContextDirectory").
		Exec().
		Contents(ctx)

	return err
}
