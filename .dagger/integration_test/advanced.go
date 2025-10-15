package main

import (
	"context"

	"dagger/dagbench-test/internal/dagger"

	"golang.org/x/sync/errgroup"
)

func testAdvanced(ctx context.Context) error {
	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error { return testAdvanceBenchInitFromExistingConfigWithModule(gctx) })

	return eg.Wait()
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

	ctr := getTestCtr("advanced-test-advance-bench-init-from-existing-config-with-module")

	mod := ctr.
		WithWorkdir("/tmp/module").
		WithExec([]string{"dagger", "init", "--sdk=go", "--name=benchmark", "--source=."}).
		Directory("/tmp/module")

	_, err := ctr.
		WithDirectory(
			"/module", mod,
			dagger.ContainerWithDirectoryOpts{
				Exclude: []string{"dagger.json", "dagger.gen.go", "internal"},
			}).
		WithNewFile("/config.json", config).
		WithExec(
			[]string{
				"run", "--workdir", "/module",
				"--config", "/config.json",
				"-i", "2", "-o", "test.txt",
			},
			dagger.ContainerWithExecOpts{
				UseEntrypoint: true,
			},
		).
		File("test.txt").
		Contents(ctx)

	return err
}
