package main

import (
	"context"
	"diLesson/application/contract"
	"diLesson/config"
	"diLesson/server"
	"github.com/golobby/container/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
)

func main() {

	log.Println("config: start")

	p, err := filepath.Abs("configs/main.yaml")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := config.Parse(p)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("config: end")

	log.Println("tracing configuration: start")

	res, err := resource.New(
		context.Background(),
		resource.WithProcess(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("pay-service"),
			semconv.ServiceVersionKey.String("0.0.1"),
			semconv.ServiceInstanceIDKey.String(uuid.New().String()),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(conf.Trace.Host+":"+conf.Trace.HttpPort),
		otlptracehttp.WithInsecure(),
	)
	defer func(exporter *otlptrace.Exporter, ctx context.Context) {
		err := exporter.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}(exporter, context.Background())

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(traceProvider)

	//TODO: is it necessary?
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer := otel.Tracer(
		"application",
		trace.WithInstrumentationVersion("v1.2.3"),
		trace.WithSchemaURL(semconv.SchemaURL),
	)

	log.Println("tracing configuration: end")

	log.Println("building dependencies: start")
	err = config.BuildDI(conf)
	if err != nil {
		log.Fatal(err)
	}

	var service contract.Charge
	err = container.Resolve(&service)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("building dependencies: end")

	s := server.NewServer(service)
	ls, err := net.Listen("tcp", ":"+conf.Grpc.Port)

	gs := grpc.NewServer()
	server.RegisterPayServiceServer(gs, s)

	log.Println("server: start")

	_, span := tracer.Start(context.Background(), "server listens")

	if err = gs.Serve(ls); err != nil {
		log.Fatal(err)
	}

	span.End()

	log.Println("server: end")
}
