package data

import (
	"context"
	v1 "github.com/go-kratos/kratos-layout/api/ping/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type (
	PingRepo struct {
		data   *Data
		logger *log.Helper
	}
)

func (p *PingRepo) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	p.logger.Infof("PingRepo.Ping() called with args: req = %v", req)
	err := p.data.db.Exec("SELECT 1").Error
	if err != nil {
		return nil, err
	}

	err = p.data.cache.Set("ping", "pong", 60*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return &v1.PingResponse{Message: "pong"}, nil
}

var _ biz.PingRepo = (*PingRepo)(nil)

func NewPingRepo(data *Data, logger log.Logger) *PingRepo {
	return &PingRepo{data: data, logger: log.NewHelper(logger)}
}
