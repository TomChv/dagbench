package main

import (
	"dagger/dagbench/internal/dagger"
	"fmt"
	"strings"
)

type CLI struct {
	//+private
	Ctr *dagger.Container
}

const (
	baseImage  = "alpine:3.22"
	daggerRepo = "https://github.com/dagger/dagger"
)

func daggerCtr(daggerVersion string) *dagger.Container {
	daggerRepo := dag.Git(daggerRepo)

	var daggerRef *dagger.GitRef

	// Naive version parsing
	if strings.HasPrefix(daggerVersion, "v") && strings.Contains(daggerVersion, ".") {
		daggerRef = daggerRepo.Tag(daggerVersion)
	} else if !strings.ContainsAny(daggerVersion, "0123456789") || strings.ContainsAny(daggerVersion, "-_/") {
		daggerRef = daggerRepo.Branch(daggerVersion)
	} else {
		daggerRef = daggerRepo.Commit(daggerVersion)
	}

	daggerSource := daggerRef.Tree(dagger.GitRefTreeOpts{DiscardGitDir: true})

	return dag.
		DaggerDev(dagger.DaggerDevOpts{
			Source: daggerSource,
		}).Dev()
}

func newCLI(binary *dagger.File, daggerVersion string) (*CLI, error) {
	return &CLI{
		Ctr: daggerCtr(daggerVersion).
			WithMountedFile("/bin/dagbench", binary).
			WithEntrypoint([]string{"/bin/dagbench"}),
	}, nil
}

func (c *CLI) New(
	//+default="bench"
	name string,

	//+default="go"
	sdk string,

	//+default="json"
	format string,
) *dagger.File {
	return c.Ctr.
		WithExec(
			[]string{"new", name, "--sdk", sdk, "--format", format},
			dagger.ContainerWithExecOpts{UseEntrypoint: true},
		).
		File(fmt.Sprintf("%s.%s", name, format))
}
