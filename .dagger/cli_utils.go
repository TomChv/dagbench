package main

import (
	"dagger/dagbench/internal/dagger"
	"strings"
)

func fetchSourceFromGit(daggerVersion string) *dagger.Directory {
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

	return daggerRef.Tree(dagger.GitRefTreeOpts{DiscardGitDir: true})
}
