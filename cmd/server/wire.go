//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/server"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			conf.GetConfigProviderSet,
			server.ProviderSet,
			service.ProviderSet,
			GetEnv,
			GetLogger,
			GetETCD,
			GetETCDRegistrar,
			GetTrace,
			newApp,
		),
	)
}
