package grpc

import (
	"github.com/rjahon/img-service/config"
	"github.com/rjahon/img-service/genproto/img_service"
	"github.com/rjahon/img-service/grpc/client"
	"github.com/rjahon/img-service/grpc/service"
	"github.com/rjahon/img-service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, strg storage.StorageI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {
	opt := grpc.MaxConcurrentStreams(100)
	grpcServer = grpc.NewServer(opt)

	img_service.RegisterServiceServer(grpcServer, service.NewImgService(cfg, strg, svcs))

	reflection.Register(grpcServer)

	return
}
