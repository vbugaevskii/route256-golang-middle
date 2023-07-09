package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
)

func Init(serviceName string) error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	_, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return err
	}

	return nil
}

func MarkSpanWithError(ctx context.Context, err error) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return err
	}

	ext.Error.Set(span, true)
	span.LogKV("error", err.Error())

	return err
}
