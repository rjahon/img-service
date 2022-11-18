package client

import (
	"github.com/rjahon/img-service/config"
)

type ServiceManagerI interface {
}

type grpcClients struct {
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {
	return &grpcClients{}, nil
}
