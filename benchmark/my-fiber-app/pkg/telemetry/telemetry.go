package telemetry

import (
	"context"
	"log"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitMetrics(app *fiber.App) {
	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping"})            // Optional: Remove some paths from metrics
	prometheus.SetIgnoreStatusCodes([]int{401, 403, 404}) // Optional: Skip metrics for these status codes

	app.Use(prometheus.Middleware)

	app.Get("/metrics-fiber", monitor.New())

}

var TracerProvider *sdktrace.TracerProvider

func InitTracing(ctx context.Context) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlptracehttp.WithInsecure(), // para localhost
	)
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("my-fiber-app"),
		)),
	)

	otel.SetTracerProvider(TracerProvider)

}
