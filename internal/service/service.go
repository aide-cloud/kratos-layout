package service

import (
	"context"
	"embed"
	"github.com/aide-cloud/graphql-http"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewGraphqlService,
	NewRoot,
)

type GraphqlService struct {
	logger *log.Helper
}

type (
	RootInterface interface {
		Ping() string
		Check(ctx context.Context, args struct{ In string }) (string, error)
	}
)

type Root struct {
	logger *log.Helper
}

func (r *Root) Ping() string {
	return "pong"
}

func (r *Root) Check(ctx context.Context, args struct{ In string }) (string, error) {
	return args.In + " is ok!", nil
}

var _ RootInterface = (*Root)(nil)

// Content holds all the SDL file content.
//go:embed sdl
var content embed.FS

func NewRoot(logger log.Logger) *Root {
	return &Root{
		logger: log.NewHelper(logger),
	}
}

func NewGraphqlService(logger log.Logger) *GraphqlService {
	return &GraphqlService{
		logger: log.NewHelper(logger),
	}
}

func (g *GraphqlService) RegisterGraphqlGinRouter(root *Root, r *gin.Engine) {
	r.GET("/query", gin.WrapF(graphql.NewGraphQLNetHttpHandlerFunc("/graphql")))
	r.POST("/graphql", gin.WrapH(graphql.NewNetHttpHandler(root, content)))
}
