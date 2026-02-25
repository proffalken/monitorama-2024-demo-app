package main

import (
	"context"
	"io"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	traceapi "go.opentelemetry.io/otel/trace"
)

var tracer traceapi.Tracer

func newTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx)
}

func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			attribute.String("service.name", "otel-basic-go-client"),
			attribute.String("service.version", "1.1.0"),
			attribute.String("service.namespace", "oteldemostack"),
			attribute.String("deployment.environment.name", "production"),
		),
	)
}

func newTraceProvider(exp sdktrace.SpanExporter, res *resource.Resource) *sdktrace.TracerProvider {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
}

func newMetricsExporter(ctx context.Context) (metric.Exporter, error) {
	return otlpmetrichttp.New(ctx)
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := newMetricsExporter(context.Background())
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(3*time.Second))),
	), nil
}

func main() {
	ctx := context.Background()

	res, err := newResource()
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	traceExporter, err := newTraceExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize trace exporter: %v", err)
	}

	tp := newTraceProvider(traceExporter, res)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("basicotel/goclient")

	meterProvider, err := newMeterProvider(res)
	if err != nil {
		log.Fatalf("failed to initialize meter provider: %v", err)
	}
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	for {
		ctx, span := tracer.Start(ctx, "call_remote_service", traceapi.WithSpanKind(traceapi.SpanKindClient))

		resp, err := otelhttp.Get(ctx, "http://localhost:8898/space_json/")
		if err != nil {
			span.RecordError(err)
			span.End()
			log.Printf("request failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			span.RecordError(err)
			span.End()
			log.Printf("read body failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Message: %s", string(body))
		span.End()
		time.Sleep(2 * time.Second)
	}
}
