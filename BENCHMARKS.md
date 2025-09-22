# Benchmarks

This document describes the steps to benchmark a dagger module with `dagbenchmark`

### 1. Initialize the module

Create a new dagger module using `dagbenchmark`, this will create a module on a temporary directory,
measure the initilization time and save a config file that you will need for the next steps.

```bash
dagbenchmark init --config-dir ./.configs/my-module --report-dir ./.reports/my-module
```

You should see something like this:

```bash
********* CONFIG **********
Dagger version  : dagger vX.XX.XX (image://registry.dagger.io/engine:vX.XX.XX) darwin/arm64/v8
Dagger path     : /Users/tomchauveau/.asdf/shims/dagger
Language        : go
Temp dir        : /var/folders/g5/2mpcw9d56cv1xybf6ns3vhzm0000gn/T/dagger-benchmark/75fa767f-cd81-4297-bcac-4786d7b9d02f

********* EXECUTION **********
Pruning cache...
Executing  init-go-vX.XX.XX

********* REPORTING **********
Report init-go-vX.XX.XX results
moduleSource    : 200ms
run codegen     : 9.5s
generatedContextDirectory       : 9.7s

Saving report at XXX/.reports/my-module/init-go-vX.XX.XX-results.csv
Saving config at XXX/.configs/my-module/go-vX.XX.XX.json
```

The config file is saved in the directory you specified with `--config-dir`. It contains required information to run further commands.

Example of the config file:

```json
{
  "language": "go",
  "binPath": "XXX/.asdf/shims/dagger",
  "tempDir": "XXXXX/T/dagger-benchmark/75fa767f-cd81-4297-bcac-4786d7b9d02f",
  "version": "vX.XX.XX"
}
```

### 2. Measure develop, functions and call times

Now that the module is initialized, you can run the develop, functions and call commands using the config file you created in the init.
You can set the flag `-r` to run multiple times, this will compute the average time in the final report.


You can run the develop, functions and call commands using the config file you created in the init in one call with `flow`.

```bash
go run ./main.go --config-file .configs/my-module/go-v0.18.19.json --report-dir ./.reports/my-module flow 'container-echo --string-arg="foo"'
```

OR you can run them one by one:

- Functions

```bash
dagbenchmark --config-file ./.configs/my-module/go-v0.18.19.json --report-dir ./.reports/my-module functions
```

- Develop

```
dagbenchmark --config-file ./.configs/my-module/go-v0.18.19.json --report-dir ./.reports/my-module develop 
```

- Call

```bash
dagbenchmark --config-file ./.configs/my-module/go-v0.18.19.json --report-dir ./.reports/my-module call 'container-echo --string-arg="foo"'
```

### 3. Measure the dev engine

You can rerun step 1 and 2 with the dev engine by configuring your environment.

Here's an example of .envrc after running `./hack/dev` on the `dagger/dagger` repo:

```bash
export DAGGER_REPO_ROOT=$HOME/Documents/github.com/dagger/dagger
export _EXPERIMENTAL_DAGGER_CLI_BIN=$DAGGER_REPO_ROOT/bin/dagger
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://dagger-engine.dev
export PATH=$DAGGER_REPO_ROOT/bin:$PATH
```

You should see different informations in the `CONFIG` section:

```bash
dagbenchmark --save-config-dir ./.configs/my-module --save-report-dir ./.reports/my-module init 

# ********* CONFIG **********
# Dagger version  : dagger vX.XX.XX-XXXXXXXXX-XXXXXXXXX (docker-container://dagger-engine.dev) darwin/arm64/v8
# Dagger path     : /XXXXXX/github.com/dagger/dagger/bin/dagger
# Language        : go
# Temp dir        : /XXXXXX/T/dagger-benchmark/baaf58fb-003c-4c16-bda5-e63236328500
```

Then run step 2 with the dev engine, you should end up with the following outputs:

```bash
❯ ls .reports/my-module 
call-containerEcho-go-v0.18.18-250911085307-97a2c3ee1a6a-results.csv functions-go-v0.18.19-results.csv
call-containerEcho-go-v0.18.19-results.csv                           init-go-v0.18.18-250911085307-97a2c3ee1a6a-results.csv
functions-go-v0.18.18-250911085307-97a2c3ee1a6a-results.csv          init-go-v0.18.19-results.csv
 
❯ ls .configs/my-module 
go-v0.18.18-250911085307-97a2c3ee1a6a.json go-v0.18.19.json
```

### 4. Measure the difference between versions

Now that you have record the 2 versions of the module, you can run the `diff` command to compare the reports.

```bash
dagbenchmark diff .reports/my-module/functions-go-v0.18.18-250911085307-97a2c3ee1a6a-results.csv .reports/my-module/functions-go-v0.18.19-results.csv 

********* DIFFERENCE **********
Report diff-functions-go-v0.18.18-250911085307-97a2c3ee1a6a-results.csv-on-functions-go-v0.18.19-results.csv results
loading type definitions        : 100ms
load module     : 100ms
finding module configuration    : 0s
getModDef       : 0s
initializing module     : 0s
```