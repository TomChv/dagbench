package integration

import (
	"dagger/dagbench/internal/dagger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type CLI interface {
	Container() *dagger.Container
}

func Tracer() trace.Tracer {
	return otel.Tracer("dagger.io/sdk.go")
}
