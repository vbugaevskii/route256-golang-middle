package metrics

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ServerMetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	MetricRequestCounter.Inc()

	tsBegin := time.Now()
	res, err := handler(ctx, req)
	tsElapsed := time.Since(tsBegin)

	retStatus, _ := status.FromError(err)
	retStatusStr := fmt.Sprint(retStatus.Code())

	MetricResponseCounter.WithLabelValues(retStatusStr).Inc()
	MetricResponseTimeHistogram.WithLabelValues(retStatusStr).Observe(tsElapsed.Seconds())

	return res, err
}

func ClientMetricsInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	tsBegin := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	tsElapsed := time.Since(tsBegin)

	retStatus, _ := status.FromError(err)
	retStatusStr := fmt.Sprint(retStatus.Code())

	MetricClientResponseTimeHistogram.WithLabelValues(retStatusStr).Observe(tsElapsed.Seconds())

	return err
}
