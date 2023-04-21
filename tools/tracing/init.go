package tracing

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
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

	return tp
}
