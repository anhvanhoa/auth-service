package grpcservice

import (
	"auth-service/bootstrap"

	grpc_service "github.com/anhvanhoa/service-core/bootstrap/grpc"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/user_context"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(
	env *bootstrap.Env,
	cacher cache.CacheI,
	log *log.LogGRPCImpl,
	authService proto_auth.AuthServiceServer,
) *grpc_service.GRPCServer {
	config := &grpc_service.GRPCServerConfig{
		IsProduction: env.IsProduction(),
		PortGRPC:     env.PortGrpc,
		NameService:  env.NameService,
	}
	middleware := grpc_service.NewMiddleware()
	return grpc_service.NewGRPCServer(
		config,
		log,
		func(server *grpc.Server) {
			proto_auth.RegisterAuthServiceServer(server, authService)
		},
		middleware.AuthorizationInterceptor(
			env.SecretService,
			func(action string, resource string) bool {
				hasPermission, err := cacher.Get(resource + "." + action)
				if err != nil {
					return false
				}
				return hasPermission != nil && string(hasPermission) == "true"
			},
			func(id string) *user_context.UserContext {
				userData, err := cacher.Get(id)
				if err != nil || userData == nil {
					return nil
				}
				uCtx := user_context.NewUserContext()
				uCtx.FromBytes(userData)
				return uCtx
			},
		),
	)
}
