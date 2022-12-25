package service

import (
	"context"
	pb "github.com/go-kratos/kratos-layout/api/ping/v1"
	"github.com/go-kratos/kratos/v2/log"
)

type (
	PingService struct {
		pb.UnimplementedPingServer

		logger *log.Helper
		logic  PingLogicInterface
	}

	PingGraphqlService struct {
		*PingService
	}

	PingLogicInterface interface {
		Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error)
	}
)

func NewPingService(logic PingLogicInterface, logger log.Logger) *PingService {
	return &PingService{logger: log.NewHelper(logger), logic: logic}
}

func NewPingGraphqlService(s *PingService) *PingGraphqlService {
	return &PingGraphqlService{PingService: s}
}

func (s *PingGraphqlService) Ping(ctx context.Context, args struct {
	In *pb.PingRequest
}) (*pb.PingResponse, error) {
	return s.PingService.Ping(ctx, args.In)
}

func (s *PingService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return Validate(ctx, req, s.logic.Ping)
}
