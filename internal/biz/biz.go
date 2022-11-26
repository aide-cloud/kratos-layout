package biz

import (
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewGreeterUsecase,
	wire.Bind(new(service.GreeterLogicInterface), new(*GreeterUsecase)),
)
