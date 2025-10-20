package main

import (
	"context"
	"fmt"

	"dagger/dagbench/internal/dagger"

	"github.com/google/uuid"
)

// Create a new Dagger Engine container using the given tag as the image tag
func daggerContainerFromTag(
	ctx context.Context,

	// The cache volume name to use.
	cacheVolumeName string,

	// The tag of the engine image to use.
	tag string,
) (*dagger.Container, error) {
	image := fmt.Sprintf("registry.dagger.io/engine:%s", tag)

	engine := dag.Container().From(image).
		WithExposedPort(1234, dagger.ContainerWithExposedPortOpts{Protocol: dagger.NetworkProtocolTcp}).
		WithMountedCache("/var/lib/dagger", dag.CacheVolume(cacheVolumeName)).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"--addr", "tcp://0.0.0.0:1234",
				"--network-name", "dagger-dev",
				"--network-cidr", "10.88.0.0/16",
			},
			UseEntrypoint:            true,
			InsecureRootCapabilities: true,
		})

	endpoint, err := engine.Endpoint(ctx, dagger.ServiceEndpointOpts{Scheme: "tcp"})
	if err != nil {
		return nil, err
	}

	return dag.
		Container().
		From(image).
		WithServiceBinding("dagger-engine", engine).
		WithEnvVariable("_EXPERIMENTAL_DAGGER_RUNNER_HOST", endpoint), nil
}

func buildUniqueCacheVolumeName(tag string) string {
	return fmt.Sprintf("dev-engine-%s-%s", tag, uuid.NewString())
}
