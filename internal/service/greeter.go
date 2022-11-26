package service

import (
	"context"
	pb "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos/v2/log"
)

type (
	GreeterService struct {
		pb.UnimplementedGreeterServer

		logger *log.Helper
		logic  GreeterLogicInterface
	}

	GreeterGraphqlService struct {
		*GreeterService
	}

	GreeterLogicInterface interface {
		SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error)
	}
)

func NewGreeterService(logic GreeterLogicInterface, logger log.Logger) *GreeterService {
	return &GreeterService{logger: log.NewHelper(logger), logic: logic}
}

func NewGreeterGraphqlService(s *GreeterService) *GreeterGraphqlService {
	return &GreeterGraphqlService{GreeterService: s}
}

func (s *GreeterGraphqlService) SayHello(ctx context.Context, args struct {
	In *pb.HelloRequest
}) (*pb.HelloReply, error) {
	return s.GreeterService.SayHello(ctx, args.In)
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return Validate(ctx, req, s.logic.SayHello)
}
