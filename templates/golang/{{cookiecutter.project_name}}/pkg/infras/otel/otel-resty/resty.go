// Package otelresty provides OpenTelemetry middleware for resty
package otelresty

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("resty-client")

// RequestMiddleware 用于在 resty 发起请求前，记录请求相关 tracing 信息
func RequestMiddleware(_ *resty.Client, req *resty.Request) error {
	ctx, span := tracer.Start(req.Context(), fmt.Sprintf("HTTP %s", req.Method),
		trace.WithAttributes(
			attribute.String("http.url", req.URL),
			attribute.String("http.method", req.Method),
		),
	)
	ctx = context.WithValue(ctx, "otel-span", span)
	req.SetContext(ctx)

	return nil
}

// ResponseMiddleware 用于在 resty 发起请求后，记录响应相关 tracing 信息
func ResponseMiddleware(_ *resty.Client, resp *resty.Response) error {
	span := resp.Request.Context().Value("otel-span").(trace.Span)
	defer span.End()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode()),
	)
	if resp.IsError() {
		span.RecordError(errors.Errorf("HTTP error: %s", resp.String()))
	}
	return nil
}
