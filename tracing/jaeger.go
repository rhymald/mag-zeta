package tracing

import (
  "go.opentelemetry.io/otel/exporters/jaeger"
  "go.opentelemetry.io/otel/sdk/resource"
  "go.opentelemetry.io/otel/sdk/trace"
  semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func JaegerTraceProvider() (*trace.TracerProvider, error) {
  exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))
  if err != nil { return nil, err }
  tp := trace.NewTracerProvider(
    trace.WithBatcher(exp),
    trace.WithResource(resource.NewWithAttributes(
      semconv.SchemaURL,
      semconv.ServiceNameKey.String("mag"),
      semconv.DeploymentEnvironmentKey.String("production"),
    )),
  )
  return tp, nil
}
