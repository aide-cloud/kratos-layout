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
	"os"
)

func SetEnv(env *conf.Env, logger log.Logger) []kratos.Option {
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

func GetETCD(registrar *conf.Registrar) *etcd.Registry {
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
	return etcd.New(client)
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

func newApp(gs *grpc.Server, hs *http.Server, etcdRegistry *etcd.Registry, opts ...kratos.Option) *kratos.App {
	opts = append(opts, kratos.Server(gs, hs), kratos.Registrar(etcdRegistry))
	return kratos.New(opts...)
}
