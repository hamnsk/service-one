package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	otrace "go.opentelemetry.io/otel/trace"
	"net/http"
)

func trace(c *fiber.Ctx) error {
	log.Info("Send Span")
	opts := []otrace.SpanStartOption{
		otrace.WithSpanKind(otrace.SpanKindServer),
	}
	httpReq := &http.Request{}
	_ = fasthttpadaptor.ConvertRequest(c.Context(), httpReq, true)

	_, _, spanContext := otelhttptrace.Extract(c.Context(), httpReq)
	reqCtx := otrace.ContextWithSpanContext(c.Context(), spanContext)

	tp := otel.GetTracerProvider()

	tr := tp.Tracer("audit-service")
	_, span := tr.Start(reqCtx, "audit-span", opts...)

	defer span.End()

	return c.Next()
}
