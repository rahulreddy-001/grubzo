package cmd

import (
	"context"
	"grubzo/internal/config"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.uber.org/zap"
)

func initTracer(service string, c *config.Config, logger *zap.Logger) (func(context.Context) error, error) {
	ctx := context.Background()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(c.TempoHost),
	)

	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(service),
		),
	)
	if err != nil {
		logger.Error("error initilizing tracer", zap.Error(err))
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(renameProcessor{}),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(

		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return tp.Shutdown, nil
}

type renameProcessor struct {
}

func (p renameProcessor) OnStart(ctx context.Context, s sdktrace.ReadWriteSpan) {
	attrs := s.Attributes()

	var dbSystem, dbStatement string
	for _, attr := range attrs {
		switch string(attr.Key) {
		case "db.system":
			dbSystem = attr.Value.AsString()
		case "db.statement":
			dbStatement = attr.Value.AsString()
		case "db.system.name":
			dbSystem = attr.Value.AsString()
		case "db.query.text":
			dbStatement = attr.Value.AsString()
		}
	}

	if dbSystem == "redis" && dbStatement != "" {
		cmd := strings.Fields(dbStatement)
		if len(cmd) > 0 {
			s.SetName("Redis." + cmd[0])
		}
	}
}

func (p renameProcessor) OnEnd(s sdktrace.ReadOnlySpan) {}

func (p renameProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (p renameProcessor) ForceFlush(ctx context.Context) error {
	return nil
}
