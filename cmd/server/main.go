package main

import (
	"flag"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"os"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	flagConf string

	id, _    = os.Hostname()
	Metadata = make(map[string]string)
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	cfg := conf.GetConfig(flagConf)
	SetEnv(cfg.GetEnv())
	app, cleanup, err := wireApp(cfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
