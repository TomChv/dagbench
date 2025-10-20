package main

import (
	"context"

	"dagger/dagbench-test/internal/dagger"

	"golang.org/x/sync/errgroup"
)

func testAdvanced(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testUseContainer(gctx) })
	eg.Go(func() error { return testAdvanceBenchInitFromExistingConfigWithModule(gctx) })

	return eg.Wait()
}

func testUseContainer(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test use container")
	defer span.End()

	_, err := getTestCLI("test-use-container").
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

func testAdvanceBenchInitFromExistingConfigWithModule(ctx context.Context) error {
	ctx, span := Tracer().Start(ctx, "test benchmark init with an existing module")
	defer span.End()

	config := `{
  "name": "test",
  "commands": [
    {
      "spanNames": [
        "generatedContextDirectory"
      ],
      "args": [
        "init",
        "--sdk=go",
        "--name=benchmark",
        "--source=."
      ]
    },
    {
      "spanNames": [
        "develop"
      ],
      "args": [
        "develop"
      ]
    }
  ]
}`

	cli := getTestCLI("advanced-test-advance-bench-init-from-existing-config-with-module")

	mod := cli.
		Container().
		WithWorkdir("/tmp/module").
		WithExec([]string{"dagger", "init", "--sdk=go", "--name=benchmark", "--source=."}).
		Directory("/tmp/module")

	_, err := cli.
		Run().
		WithWorkdir("/module", dagger.DagbenchCmdRunWithWorkdirOpts{
			Directory: dag.Directory().WithDirectory(".", mod, dagger.DirectoryWithDirectoryOpts{
				Exclude: []string{"dagger.json", "dagger.gen.go", "internal"},
			}),
		}).
		WithConfigFile(dag.File("config.json", config)).
		WithIteration(2).
		Exec().
		Contents(ctx)

	return err
}
