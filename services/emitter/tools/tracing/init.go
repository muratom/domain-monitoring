package tracing

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer(serviceName, jaegerCollectorURL string) *tracesdk.TracerProvider {
	exporter, err := NewJaegerExporter(jaegerCollectorURL)
	if err != nil {
		logrus.Fatalf("failed to create the jaeger exporter: %v", err)
	}

	tp, err := NewTraceProvider(serviceName, exporter)
	if err != nil {
		logrus.Fatalf("failed to create the tracer provider: %v", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}
