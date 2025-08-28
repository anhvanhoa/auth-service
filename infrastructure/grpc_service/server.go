package grpcservice

import (
	"auth-service/bootstrap"

	grpc_server "github.com/anhvanhoa/service-core/boostrap/grpc"
	"github.com/anhvanhoa/service-core/domain/log"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(env *bootstrap.Env, authService proto_auth.AuthServiceServer, log *log.LogGRPCImpl) *grpc_server.GRPCServer {
	config := &grpc_server.GRPCServerConfig{
		IsProduction: env.IsProduction(),
		PortGRPC:     env.PORT_GRPC,
		NameService:  env.NAME_SERVICE,
	}
	return grpc_server.NewGRPCServer(
		config,
		log,
		func(server *grpc.Server) {
			proto_auth.RegisterAuthServiceServer(server, authService)
		},
	)
}
