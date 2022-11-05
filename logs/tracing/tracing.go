package tracing

import (
	"context"
	"fmt"
	"github.com/taorzhang/toolkit/errs"
	"github.com/taorzhang/toolkit/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	log = logs.NewLogger("module", "tracing")
)

type TraceCfg struct {
	endpoint         string
	serviceNamespace string
	serviceName      string
	sampler          sdk_trace.Sampler
	propagator       propagation.TextMapPropagator
	shutdown         func()
}

func NewTraceCfg(endpoint, serviceNamespace, serviceName string) *TraceCfg {
	return &TraceCfg{endpoint: endpoint, serviceName: serviceName, serviceNamespace: serviceNamespace}
}

type TraceCfgOpt func(t *TraceCfg) error

func WithSampler(fraction float64) TraceCfgOpt {
	return func(t *TraceCfg) error {
		t.sampler = sdk_trace.TraceIDRatioBased(fraction)
		return nil
	}
}

func InitOpenTrace(ctx context.Context, cfg *TraceCfg, opts ...TraceCfgOpt) errs.IError {
	for idx := range opts {
		if err := opts[idx](cfg); err != nil {
			return errs.New(errs.InvalidParams, fmt.Sprintf("run trace cft opt failed, and err is %v", err))
		}
	}
	if cfg.endpoint == "" || cfg.serviceNamespace == "" || cfg.serviceName == "" {
		return errs.New(errs.InvalidParams, fmt.Sprintf("endpoint: %s, serviceNamespace: %s, serviceName: %s", cfg.endpoint, cfg.serviceNamespace, cfg.serviceName))
	}
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.endpoint)))
	if err != nil {
		return errs.New(errs.InvalidInitComponent, fmt.Sprintf("new jaeger exporter failed, endpoints:%s, err message: %s", cfg.endpoint, err))
	}
	traceProvider := sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNamespaceKey.String(cfg.serviceNamespace),
			semconv.ServiceNameKey.String(cfg.serviceName))),
		sdk_trace.WithSampler(cfg.sampler))
	otel.SetTracerProvider(traceProvider)
	if cfg.propagator != nil {
		otel.SetTextMapPropagator(cfg.propagator)
	} else {
		otel.SetTextMapPropagator(propagation.TraceContext{})
	}

	// waiting shutdown
	go func() {
		for {
			select {
			// ctx.Done是幂等
			case <-ctx.Done():
				if err = traceProvider.Shutdown(ctx); err != nil {
					log.Warn(ctx, "trace provider is ready to shutdown", "ctx_err", ctx.Err())
					return
				}
			}
		}
	}()
	return nil
}

func StartSpan(ctx context.Context, traceName, spanName string) (spanCtx context.Context, span trace.Span) {
	return otel.Tracer(traceName).Start(ctx, spanName)
}
