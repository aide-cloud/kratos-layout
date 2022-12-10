package main

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

func GetEnv(env *conf.Env, logger log.Logger) []kratos.Option {
	Name = env.GetName()
	Metadata = env.GetMetadata()
	opts := []kratos.Option{
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(Metadata),
		kratos.Logger(logger),
	}
	return opts
}

func SetEnv(env *conf.Env) {
	Name = env.GetName()
	Metadata = env.GetMetadata()
	if Version == "" {
		Version = env.GetVersion()
	}
}

func GetETCD(registrar *conf.Registrar) *clientV3.Client {
	endpoints := registrar.GetEtcd().GetEndpoints()
	if len(endpoints) == 0 {
		panic("etcd endpoints is empty")
	}
	// new etcd client
	client, err := clientV3.New(clientV3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		panic(err)
	}
	// new reg with etcd client
	return client
}

func GetETCDRegistrar(etcdClient *clientV3.Client) *etcd.Registry {
	return etcd.New(etcdClient)
}

func GetLogger(conf *conf.Log) log.Logger {
	return log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}

func GetTrace(conf *conf.Trace) *traceSdk.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.GetEndpoint())))
	if err != nil {
		panic(err)
	}
	tp := traceSdk.NewTracerProvider(
		traceSdk.WithBatcher(exp),
		traceSdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
			semconv.ServiceVersionKey.String(Version),
		)),
	)
	return tp
}

func newApp(gs *grpc.Server, hs *http.Server, etcdRegistry *etcd.Registry, opts ...kratos.Option) *kratos.App {
	opts = append(opts, kratos.Server(gs, hs))
	if etcdRegistry != nil {
		//opts = append(opts, kratos.Registrar(etcdRegistry))
	}
	return kratos.New(opts...)
}
