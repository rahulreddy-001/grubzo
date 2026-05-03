package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName          = "grubzo.mcp"
	payloadPreviewLimit = 2048
)

type payloadSummary struct {
	size      int
	preview   string
	truncated bool
	ok        bool
}

func Start(ctx context.Context, name string, payload any, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	allAttrs := append(
		[]attribute.KeyValue{attribute.String("mcp.span.name", name)},
		callerAttributes(1)...,
	)
	allAttrs = append(allAttrs, attrs...)
	allAttrs = append(allAttrs, payloadAttributes(payload)...)
	addEventToParent(ctx, name, allAttrs)
	return otel.Tracer(tracerName).Start(ctx, name, trace.WithAttributes(allAttrs...))
}

func StartServer(ctx context.Context, name string, payload any, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	allAttrs := append(
		[]attribute.KeyValue{attribute.String("mcp.span.name", name)},
		callerAttributes(1)...,
	)
	allAttrs = append(allAttrs, attrs...)
	allAttrs = append(allAttrs, payloadAttributes(payload)...)
	addEventToParent(ctx, name, allAttrs)
	return otel.Tracer(tracerName).Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(allAttrs...),
	)
}

func StartClient(
	ctx context.Context,
	name string,
	method string,
	rawURL string,
	payload any,
	attrs ...attribute.KeyValue,
) (context.Context, trace.Span) {
	allAttrs := append(
		[]attribute.KeyValue{attribute.String("mcp.span.name", name)},
		callerAttributes(1)...,
	)
	allAttrs = append(allAttrs, attrs...)
	if method != "" {
		allAttrs = append(allAttrs, semconv.HTTPRequestMethodKey.String(method))
	}
	if rawURL != "" {
		allAttrs = append(allAttrs, semconv.URLFull(rawURL))
		if parsedURL, err := url.Parse(rawURL); err == nil && parsedURL.Hostname() != "" {
			allAttrs = append(allAttrs, semconv.ServerAddress(parsedURL.Hostname()))
		}
	}
	summary := summarizePayload(payload)
	if summary.ok {
		allAttrs = append(
			allAttrs,
			semconv.HTTPRequestBodySize(summary.size),
			attribute.Int("mcp.payload.size", summary.size),
			attribute.String("mcp.payload.preview", summary.preview),
		)
		if summary.truncated {
			allAttrs = append(allAttrs, attribute.Bool("mcp.payload.truncated", true))
		}
	}
	addEventToParent(ctx, name, allAttrs)

	return otel.Tracer(tracerName).Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(allAttrs...),
	)
}

func End(span trace.Span, err *error) {
	if err != nil && *err != nil {
		RecordError(span, *err)
	}
	span.End()
}

func RecordError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(semconv.ErrorTypeKey.String(fmt.Sprintf("%T", err)))
}

func InjectHTTPHeaders(ctx context.Context, header http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
}

func callerAttributes(skip int) []attribute.KeyValue {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}

	functionName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		functionName = fn.Name()
	}

	return []attribute.KeyValue{
		semconv.CodeFilepath(filepath.ToSlash(file)),
		semconv.CodeFunction(functionName),
		semconv.CodeLineNumber(line),
	}
}

func payloadAttributes(payload any) []attribute.KeyValue {
	summary := summarizePayload(payload)
	if !summary.ok {
		return nil
	}

	attrs := []attribute.KeyValue{
		attribute.Int("mcp.payload.size", summary.size),
		attribute.String("mcp.payload.preview", summary.preview),
	}
	if summary.truncated {
		attrs = append(attrs, attribute.Bool("mcp.payload.truncated", true))
	}
	return attrs
}

func summarizePayload(payload any) payloadSummary {
	if payload == nil {
		return payloadSummary{}
	}

	var raw []byte
	switch typed := payload.(type) {
	case []byte:
		raw = typed
	case json.RawMessage:
		raw = typed
	case string:
		raw = []byte(typed)
	default:
		bytes, err := json.Marshal(payload)
		if err != nil {
			raw = []byte(fmt.Sprintf("%T", payload))
		} else {
			raw = bytes
		}
	}

	if len(raw) == 0 {
		return payloadSummary{}
	}

	preview := string(raw)
	truncated := false
	if len(preview) > payloadPreviewLimit {
		preview = preview[:payloadPreviewLimit]
		truncated = true
	}

	return payloadSummary{
		size:      len(raw),
		preview:   preview,
		truncated: truncated,
		ok:        true,
	}
}

func addEventToParent(ctx context.Context, name string, attrs []attribute.KeyValue) {
	parent := trace.SpanFromContext(ctx)
	if parent == nil || !parent.IsRecording() {
		return
	}
	parent.AddEvent("mcp."+name, trace.WithAttributes(attrs...))
}
