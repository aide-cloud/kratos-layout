// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/data"
	"github.com/go-kratos/kratos-layout/internal/server"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(string2 string, logger log.Logger) (*kratos.App, func(), error) {
	bootstrap := conf.GetConfig(string2)
	confServer := bootstrap.Server
	confData := bootstrap.Data
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	graphqlService := service.NewGraphqlService(logger)
	greeterGraphqlService := service.NewGreeterGraphqlService(greeterService)
	root := service.NewRoot(greeterGraphqlService, logger)
	engine := server.GetGinEngine(confServer, graphqlService, root, logger)
	httpServer := server.NewHTTPServer(confServer, engine, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
