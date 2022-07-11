package lib

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"time"
)

const (
	instrumentName    = "lib"
	instrumentVersion = "v1.2.3"
)

var (
	tracer = otel.Tracer(
		instrumentName,
		trace.WithInstrumentationVersion(instrumentVersion),
		trace.WithSchemaURL(semconv.SchemaURL))
)

func DoSomething(ctx context.Context, value string) {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("DoSomething with value %s", value))
	defer span.End()

	time.Sleep(time.Duration(rand.Int63n(10)))

}
