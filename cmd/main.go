package main

import (
	"auth-service/bootstrap"
	"auth-service/infrastructure/grpc_client"
	grpcservice "auth-service/infrastructure/grpc_service"
	"context"

	gc "github.com/anhvanhoa/service-core/domain/grpc_client"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	log := app.Log
	db := app.DB
	cache := app.Cache
	queueClient := app.Queue

	clientFactory := gc.NewClientFactory(env.GrpcClients...)
	mailService := grpc_client.NewMailService(clientFactory.GetClient(env.MailServiceAddr))
	permissionClient, err := grpc_client.NewPermissionClient(clientFactory.GetClient(env.PermissionServiceAddr))
	if err != nil {
		log.Error("Failed to create permission client: " + err.Error())
	}

	authService := grpcservice.NewAuthService(db, env, log, mailService, permissionClient, queueClient, cache)
	grpcSrv := grpcservice.NewGRPCServer(env, cache, log, authService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	permissions := app.Helper.ConvertResourcesToPermissions(grpcSrv.GetResources())
	if _, err := permissionClient.PermissionService().RegisterPermission(ctx, permissions); err != nil {
		log.Error("Failed to register permission: " + err.Error())
	}
	if err := grpcSrv.Start(ctx); err != nil {
		log.Fatal("gRPC server error: " + err.Error())
	}
}
