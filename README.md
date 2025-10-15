# Dagger Benchmark

A CLI tool to benchmark dagger modules.

## Installation

You can either install the binary using the `dagger module`

```shell
dagger -m github.com/TomChv/dagbenchmark -c 'bin | export dagbench'
```

Or directly use `dagbench` inside the CLI container.

```shell
dagger -m github.com/TomChv/dagbenchmark -c 'cli | container | terminal' 
```

Or use it as a [module dependency](https://daggerverse.dev/mod/github.com/TomChv/dagbench) to create benchmarks
shaped on your needs.

```shell
dagger install github.com/TomChv/dagbenchmark
```

You can also specify a dagger engine version using shell:

```shell
dagger -m github.com/TomChv/dagbenchmark -c 'cli --dagger-tag v0.16.0 | container | terminal' 
```

Or provide the dagger engine container directly:

```shell
dagger -m github.com/TomChv/dagbenchmark -c 'cli --dagger-ctr $(github.com/dagger/dagger@main | dev) | container | terminal' 
```

## Usage

`dagbench` produces a go benchmark report compatible with [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) so you
can compare outputs.

See [examples](#examples) for more details on how to produce these report using `dagbench`.

`dagbench` measure the duration of a specific span and the total duration of the given command.

Example of a benchmark line for the command `"functions"` with the measure span named `"load module"`:

```
BenchmarkExample/functions/loadModule   1  20.7 s/op
BenchmarkExample/functions/cmdDuration  1  21.6 s/op
```


## CLI

```
Run a dagger benchmark

Usage:
  dagbench run [flags]

Flags:
      --auto-init             Automatically init the module using provided flags
  -c, --command string        Command to run
      --config string         Config file to use
  -d, --dagger-bin string     Dagger binary to use (default "dagger")
      --debug                 Enable debug mode
  -h, --help                  help for run
  -i, --iteration int         Number of iterations to run (default 10)
  -m, --module string         Module to use for the benchmark
      --module-name string    Name of the module to init
  -n, --name string           Name of the benchmark
  -o, --output string         Output file for the report (default "out.txt")
      --sdk string            Language to use for benchmark
  -s, --span string           Span name to record
      --template-dir string   Template directory for the benchmark
      --use-cloud             If enable, --cloud will be set
      --workdir string        Working directory for the benchmark
```

### Examples

#### Run a standalone benchmark from CLI flags

You can run a standalone benchmark from CLI flags that will auto initialize a module.

```shell
dagbench run \
--name example \
--auto-init \
--sdk go \
-o example.txt \
--command "call container-echo --string-arg=hello" \
--span "containerEcho" \
--iteration 2
```

#### Run a benchmark on an existing module

You can run a benchmark on an existing module by using the `--module` (or `-m`) flag.

```shell
dagger init --sdk=go --name=example --source=.

dagbench run \
--name example \
-m . \
--command "call container-echo --string-arg=hello" \
--span "containerEcho" \
--iteration 2 \
-o example.txt
```

#### Run a benchmark from a config file

You can run a benchmark from a config file by using the `--config` (or `-c`) flag.

This is useful if you want to make more complex workflow for your benchmark like testing the full SDK performances.
Some pre-defined configuration (called recipes) are available using the `--recipe` (or `-r`) flag.

**1. Create a config file**

```shell
dagbench new --name example --sdk go --recipe sdk --auto-init
```

This should create a file `example.json` in the current directory.

```json
{
  "name": "example",
  "iteration": 10,
  "binPath": "dagger",
  "version": "dagger v0.19.0 (image://registry.dagger.io/engine:v0.19.0) darwin/arm64/v8",
  "init": {
    "name": "example",
    "sdk": "go"
  },
  "commands": [
    {
      "spanNames": [
        "develop"
      ],
      "args": [
        "develop"
      ]
    },
    {
      "spanNames": [
        "load module"
      ],
      "args": [
        "functions"
      ]
    },
    {
      "spanNames": [
        "load module",
        "containerEcho"
      ],
      "args": [
        "call",
        "container-echo",
        "--string-arg=hello"
      ]
    }
  ],
  "cloud": false
}
```

**2. Run using the config file**

```
dagbench run --config example.json -i 2
```


**3. Override the config file using CLI flags**

You can override any file of the config file using CLI flags, since is useful to test different SDK or version
without having to create a new config file.

For example, let's run the same benchmark but with the TypeScript SDK.

```shell
dagbench run --config example.json -i 2 --sdk typescript
```

## CI

**Run integration tests**

```shell
dagger -m .dagger/integration_test -c all
```

**Build the binary**

```shell
dagger -m .dagger/ci -c 'build | export dagbench'
```

**Lint project**

```shell
dagger -m .dagger/ci -c 'lint'
```

**Format project**

```shell
dagger -m .dagger/ci -c 'fmt | export .'
```