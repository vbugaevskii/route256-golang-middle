package logger

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)

	if err != nil {
		Warn(
			"gRPC",
			zap.String("server", info.FullMethod),
			zap.String("method", info.FullMethod),
			zap.String("err", err.Error()),
		)
		return nil, err
	}

	Debug(
		"gRPC",
		zap.String("server", info.FullMethod),
		zap.String("method", info.FullMethod),
	)

	return res, nil
}
