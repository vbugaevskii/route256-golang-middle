package tracing

import (
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
