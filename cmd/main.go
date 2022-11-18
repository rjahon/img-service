package main

import (
	"context"
	"log"
	"net"

	"github.com/rjahon/img-service/config"
	"github.com/rjahon/img-service/grpc"
	"github.com/rjahon/img-service/grpc/client"
	"github.com/rjahon/img-service/storage/postgres"
)

func main() {
	cfg := config.Load()

	pg, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", err.Error())
	}
	defer pg.CloseDB()

	svcs, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", err.Error())
	}

	grpcServer := grpc.SetUpServer(cfg, pg, svcs)

	lis, err := net.Listen("tcp", cfg.ServicePort)
	if err != nil {
		log.Panic("net.Listen", err.Error())
	}

	log.Printf("GRPC started on port %s", cfg.ServicePort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", err.Error())
	}
}
