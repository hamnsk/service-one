package main

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"os"
)

type AppTracer struct {
	*tracesdk.TracerProvider
}

func InitTracing() (AppTracer, error) {
	host := os.Getenv("JAEGER_AGENT_HOST")
	port := os.Getenv("JAEGER_AGENT_PORT")
	if host == "" || port == "" {
		log.Fatal("Could not parse Jaeger env vars. Please set JAEGER_AGENT_HOST & JAEGER_AGENT_PORT ")
	}

	tp, err := tracerProvider(fmt.Sprintf("http://%s:%s/api/traces", host, port),
		"audit-service",
		"prod",
		1,
	)
	if err != nil {
		log.Fatal("ERROR: cannot init Jaeger: " + err.Error())
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return AppTracer{tp}, nil

}

func tracerProvider(url, service, environment string, id int64) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	bsp := tracesdk.NewBatchSpanProcessor(exp)
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSpanProcessor(bsp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}
