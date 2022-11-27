package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/google/wire"
)

var GetConfigProviderSet = wire.NewSet(
	wire.FieldsOf(new(*Bootstrap), "Data"),
	wire.FieldsOf(new(*Bootstrap), "Server"),
	wire.FieldsOf(new(*Bootstrap), "Registrar"),
	wire.FieldsOf(new(*Bootstrap), "Env"),
	wire.FieldsOf(new(*Bootstrap), "Log"),
	wire.FieldsOf(new(*Bootstrap), "Discovery"),
	wire.FieldsOf(new(*Bootstrap), "Trace"),
)

func GetConfig(flagconf string) *Bootstrap {
	c := config.New(config.WithSource(file.NewSource(flagconf)))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	return &bc
}
