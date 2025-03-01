package server

import (
	prometheus "github.com/aide-cloud/prom"
	"github.com/gin-gonic/gin"
	kGin "github.com/go-kratos/gin"
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
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

type GraphqlServer interface {
	RegisterGraphqlGinRouter(root *service.Root, r *gin.Engine)
}

var _ GraphqlServer = (*service.GraphqlService)(nil)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, engine *gin.Engine, _ log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// swagger api
	srv.HandlePrefix("/q/", openapiv2.NewHandler())

	// graphql 使用gin作为统一路由
	srv.HandlePrefix("/", engine)
	return srv
}

// GetGinEngine 获取gin引擎
func GetGinEngine(c *conf.Server, pingService *service.PingService, graphqlServer *service.GraphqlService, root *service.Root, tp *traceSdk.TracerProvider, logger log.Logger) *gin.Engine {
	if c.Http.GetMode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine := gin.New()
	ginEngine.Use(
		kGin.Middlewares(
			recovery.Recovery(),
			logging.Server(logger),
			validate.Validator(),
			tracing.Server(tracing.WithTracerProvider(tp)),
			ratelimit.Server(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(prometheus.MetricSeconds)),
				metrics.WithRequests(prom.NewCounter(prometheus.MetricRequests)),
			),
		),
	)

	ginEngine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	graphqlServer.RegisterGraphqlGinRouter(root, ginEngine)
	v1.RegisterPingGinHTTPServer(v1.NewPing(pingService, v1.WithRouter(ginEngine)))
	return ginEngine
}
