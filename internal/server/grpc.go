package server

import (
	prometheus "github.com/aide-cloud/prom"
	v1 "github.com/go-kratos/kratos-layout/api/ping/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, pingService *service.PingService, tp *traceSdk.TracerProvider, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(tracing.WithTracerProvider(tp)),
			ratelimit.Server(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(prometheus.MetricSeconds)),
				metrics.WithRequests(prom.NewCounter(prometheus.MetricRequests)),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterPingServer(srv, pingService)
	return srv
}
