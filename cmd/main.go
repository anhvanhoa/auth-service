package main

import (
	"auth-service/bootstrap"
	"auth-service/infrastructure/grpc_client"
	grpcservice "auth-service/infrastructure/grpc_service"
	"context"

	"github.com/anhvanhoa/service-core/domain/discovery"
	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	cache := app.Cache
	queueClient := app.Queue
	discoveryConfig := discovery.DiscoveryConfig{
		ServiceName:   env.NameService,
		ServicePort:   env.PortGrpc,
		ServiceHost:   env.HostGrpc,
		IntervalCheck: env.IntervalCheck,
		TimeoutCheck:  env.TimeoutCheck,
	}
	discoveryClient, err := discovery.NewDiscovery(&discoveryConfig)
	if err != nil {
		log.Fatal("Failed to create discovery client: " + err.Error())
	}
	discoveryClient.Register()
	defer discoveryClient.Close(env.NameService)

	clientFactory := gc.NewClientFactory(env.GrpcClients...)
	client := clientFactory.GetClient(env.MailServiceAddr)
	mailService := grpc_client.NewMailService(client)

	authService := grpcservice.NewAuthService(db, env, log, mailService, queueClient, cache)
	grpcSrv := grpcservice.NewGRPCServer(env, authService, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
