package jtrace

import (
	"blogfa/auth/config"
	"context"
	"io"

	zapLogger "blogfa/auth/pkg/logger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

var (
	tracer opentracing.Tracer
	Tracer itracer = &jtracer{}
)

type itracer interface {
	Connect() (io.Closer, error)
	GetTracer() opentracing.Tracer
	FromContext(ctx context.Context, startName string) opentracing.Span
	StartSpan(str string) opentracing.Span
	ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context
	SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
	ChildOf(span opentracing.Span, name string) opentracing.Span
}

type jtracer struct{}

func (j *jtracer) Connect() (io.Closer, error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: config.Global.Service.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.Global.Jaeger.LogSpans,
			LocalAgentHostPort: config.Global.Jaeger.HostPort,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	var closer io.Closer
	var err error
	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		logger := zapLogger.GetZapLogger(false)
		zapLogger.Prepare(logger).Development().Level(zap.InfoLevel).Add("msg", "during Listen jaeger err").Commit(err.Error())

		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

func (j *jtracer) GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}

func (j *jtracer) FromContext(ctx context.Context, startName string) opentracing.Span {

	// if context has a span for tracing then use spanFromContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pctx := parent.Context()
		if trc := opentracing.GlobalTracer(); trc != nil {
			spn := trc.StartSpan(startName, opentracing.ChildOf(pctx))
			return spn
		}
	}

	// if we havent span in context, create new span
	return opentracing.GlobalTracer().StartSpan(startName)
}

func (j *jtracer) StartSpan(str string) opentracing.Span {
	return opentracing.GlobalTracer().StartSpan(str)
}

func (j *jtracer) ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	if qr := ctx.Value("span"); qr != nil {
		ctx := context.Background()
		return opentracing.ContextWithSpan(ctx, span)
	}
	return opentracing.ContextWithSpan(ctx, span)
}

func (j *jtracer) SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, name, opts...)
}

func (j *jtracer) ChildOf(span opentracing.Span, name string) opentracing.Span {
	return opentracing.StartSpan(name, opentracing.ChildOf(span.Context()))
}
