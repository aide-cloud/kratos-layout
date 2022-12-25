package biz

import (
	"context"
	v1 "github.com/go-kratos/kratos-layout/api/ping/v1"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2/log"
)

type (
	PingRepo interface {
		Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error)
	}

	PingLogic struct {
		repo   PingRepo
		logger *log.Helper
	}
)

func (p *PingLogic) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	return p.repo.Ping(ctx, req)
}

var _ service.PingLogicInterface = (*PingLogic)(nil)

func NewPingLogic(repo PingRepo, logger log.Logger) *PingLogic {
	return &PingLogic{repo: repo, logger: log.NewHelper(logger)}
}
