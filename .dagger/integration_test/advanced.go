package main

import (
	"context"
	"testing"

	"dagger/dagbench-test/internal/dagger"

	"github.com/dagger/testctx"
	"github.com/stretchr/testify/require"
)

type AdvancedSuite struct{}

func TestAdvanced(t *testing.T) {
	testctx.New(t, Middlewares()...).RunTests(AdvancedSuite{})
}

func (AdvancedSuite) TestUseContainer(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}

func (AdvancedSuite) TestBenchInitFromExistingConfigWithModule(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}
