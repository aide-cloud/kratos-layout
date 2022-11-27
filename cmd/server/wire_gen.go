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
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(string2 string) (*kratos.App, func(), error) {
	bootstrap := conf.GetConfig(string2)
	confServer := bootstrap.Server
	confData := bootstrap.Data
	registrar := bootstrap.Registrar
	client := GetETCD(registrar)
	discovery := bootstrap.Discovery
	clientConn := data.GetRPCConn(client, discovery)
	greeterClient := data.GetGreeterClient(clientConn)
	log := bootstrap.Log
	logger := GetLogger(log)
	dataData, cleanup, err := data.NewData(confData, greeterClient, logger)
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
	registry := GetETCDRegistrar(client)
	env := bootstrap.Env
	v := GetEnv(env, logger)
	app := newApp(grpcServer, httpServer, registry, v...)
	return app, func() {
		cleanup()
	}, nil
}
