package main

import (
	"context"
	"testing"

	"github.com/dagger/testctx"
	"github.com/stretchr/testify/require"
)

type BasicSuite struct{}

func TestBasic(t *testing.T) {
	testctx.New(t, Middlewares()...).RunTests(BasicSuite{})
}

func (BasicSuite) TestRunFromCtr(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}

func (BasicSuite) TestRunFromConfigFile(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}

func (BasicSuite) TestRunFromExistingModule(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}

func (BasicSuite) TestRunWithTemplate(ctx context.Context, t *testctx.T) {
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

	require.NoError(t, err)
}

func (BasicSuite) TestRunWithWorkdirAndCleanup(ctx context.Context, t *testctx.T) {
	_, err := getTestCLI("basic-run-with-workdir-and-cleanup").
		Run().
		WithName("test").
		WithWorkdir("/tmp/workdir").
		WithIteration(2).
		WithCommand([]string{"init", "--sdk=go", "--name=test", "--source=."}).
		WithSpan("generatedContextDirectory").
		Exec().
		Contents(ctx)

	require.NoError(t, err)
}
