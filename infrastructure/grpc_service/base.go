package grpcservice

import (
	"auth-service/bootstrap"
	"auth-service/domain/usecase"
	"auth-service/infrastructure/grpc_client"
	"auth-service/infrastructure/repo"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/goid"
	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/domain/saga"
	"github.com/anhvanhoa/service-core/domain/token"
	"github.com/anhvanhoa/service-core/domain/transaction"
	proto_auth "github.com/anhvanhoa/sf-proto/gen/auth/v1"

	"github.com/go-pg/pg/v10"
)

type authService struct {
	proto_auth.UnimplementedAuthServiceServer
	env              *bootstrap.Env
	log              *log.LogGRPCImpl
	uuid             goid.GoUUID
	mailService      *grpc_client.MailService
	checkTokenUc     usecase.CheckTokenUsecase
	loginUc          usecase.LoginUsecase
	registerUc       usecase.RegisterUsecase
	refreshUc        usecase.RefreshUsecase
	logoutUc         usecase.LogoutUsecase
	verifyAccountUc  usecase.VerifyAccountUsecase
	forgotPasswordUc usecase.ForgotPasswordUsecase
	resetCodeUc      usecase.ResetPasswordByCodeUsecase
	resetTokenUc     usecase.ResetPasswordByTokenUsecase
	checkCodeUc      usecase.CheckCodeUsecase
}

func NewAuthService(
	db *pg.DB,
	env *bootstrap.Env,
	log *log.LogGRPCImpl,
	mailService *grpc_client.MailService,
	queueClient queue.QueueClient,
	cache cache.CacheI,
) proto_auth.AuthServiceServer {
	userRepo := repo.NewUserRepository(db)
	sessionRepo := repo.NewSessionRepository(db)
	tx := transaction.NewTransaction(db)
	saga := saga.NewSagaManager()
	argonService := hashpass.NewArgon()
	genUUID := goid.NewGoId().UUID()
	tokenAccess := token.NewToken(env.JwtSecret.Access)
	tokenRefresh := token.NewToken(env.JwtSecret.Refresh)
	tokenAuth := token.NewToken(env.JwtSecret.Verify)
	tokenForgot := token.NewToken(env.JwtSecret.Forgot)
	return &authService{
		env:          env,
		log:          log,
		mailService:  mailService,
		uuid:         genUUID,
		checkTokenUc: usecase.NewCheckTokenUsecase(sessionRepo),
		loginUc: usecase.NewLoginUsecase(
			userRepo,
			sessionRepo,
			tokenAccess,
			tokenRefresh,
			argonService,
			cache,
		),
		registerUc: usecase.NewRegisterUsecase(
			userRepo,
			sessionRepo,
			tokenAuth,
			tx,
			genUUID,
			argonService,
			cache,
			queueClient,
			saga,
		),
		refreshUc: usecase.NewRefreshUsecase(
			sessionRepo,
			tokenAccess,
			tokenRefresh,
			cache,
		),
		logoutUc: usecase.NewLogoutUsecase(
			sessionRepo,
			tokenAccess,
			cache,
		),
		verifyAccountUc: usecase.NewVerifyAccountUsecase(
			userRepo,
			sessionRepo,
			tokenAuth,
			cache,
		),
		forgotPasswordUc: usecase.NewForgotPasswordUsecase(
			userRepo,
			sessionRepo,
			tx,
			tokenForgot,
			cache,
			queueClient,
			saga,
			log,
		),
		resetCodeUc: usecase.NewResetPasswordCodeUsecase(
			userRepo,
			sessionRepo,
			cache,
			tokenForgot,
			argonService,
		),
		resetTokenUc: usecase.NewResetPasswordTokenUsecase(
			userRepo,
			sessionRepo,
			cache,
			tokenForgot,
			argonService,
		),
		checkCodeUc: usecase.NewCheckCodeUsecase(
			userRepo,
			sessionRepo,
		),
	}
}
