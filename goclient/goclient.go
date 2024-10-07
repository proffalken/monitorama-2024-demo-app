package main

import (
    "context"
    "io"
    "log"
    "time"

  	"go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	trace_api "go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/metric"
   	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
   	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
   	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
   	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

)

// Create the tracer
var tracer trace_api.Tracer

// Ensure we are exporting to HTTP(S)
func newTraceExporter(ctx context.Context) (trace.SpanExporter, error) {
	return otlptracehttp.New(ctx)
}

// Create the provider along with the resource attributes
func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.

    otel.SetTextMapPropagator(propagation.TraceContext{})

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("OTEL Basic Go Client"),
			semconv.ServiceVersion("0.1.0"),
			semconv.ServiceNamespace("oteldemostack"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

// Add more resource
func newResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("OTEL Basic Go Client"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func newMetricsExporter(ctx context.Context) (metric.Exporter, error) {
	return otlpmetrichttp.New(ctx)
}

func newMeterProvider(res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := newMetricsExporter(
        context.Background(),
    )
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}


func main() {

    ctx := context.Background()

    headers := map[string]string{
  "content-type": "application/json",
}

// Start tracing
exp, err := otlptrace.New(
  context.Background(),
  otlptracehttp.NewClient(
   otlptracehttp.WithEndpoint("localhost:4318"),
   otlptracehttp.WithHeaders(headers),
   otlptracehttp.WithInsecure(),
  ),
)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer("basicotel/goclient")

    res, err := newResource()

   	// Create a meter provider.
	// You can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	meterProvider, err := newMeterProvider(res)

    // Handle shutdown properly so nothing leaks.
	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
	}()

   	// Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used, which fails to generate data.
	otel.SetMeterProvider(meterProvider)

    for true {
        // Start our span
       ctx, span := tracer.Start(ctx, "call_remote_service",
                                trace_api.WithSpanKind(trace_api.SpanKindClient))

       // Set the attributes
       span.SetAttributes(attribute.String("kind", "CLIENT"))

       // Call the URI
       resp, err := otelhttp.Get(ctx, "http://localhost:8898/space_json/")
       if err != nil {
          log.Fatalln(err)
       }

    //We Read the response body on the line below.
       body, err := io.ReadAll(resp.Body)
       if err != nil {
          log.Fatalln(err)
       }

    //Convert the body to type string
       sb := string(body)
       // Print the message and context to STDOUT
       log.Printf("Message: %s, Context: %s", sb, ctx)
       // End the span
       span.End()
    // Sleep for 2 seconds
    time.Sleep(2000 * time.Millisecond)
    }
}
