package bootstrap

import (
	"strings"

	"github.com/anhvanhoa/service-core/boostrap/config"
	"github.com/anhvanhoa/service-core/domain/grpc_client"
)

type jwtSecret struct {
	Access  string
	Refresh string
	Verify  string
	Forgot  string
}

type dbCache struct {
	Addr        string
	DB          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
}

type queue struct {
	Addr        string
	DB          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
	Concurrency int
	Queues      map[string]int
}

type Env struct {
	NODE_ENV string

	URL_DB string

	NAME_SERVICE   string
	PORT_GRPC      int
	HOST_GRPC      string
	INTERVAL_CHECK string
	TIMEOUT_CHECK  string

	DB_CACHE *dbCache

	SECRET_OTP string

	QUEUE *queue

	JWT_SECRET *jwtSecret

	FRONTEND_URL string

	MAIL_SERVICE_ADDR string

	GRPC_CLIENTS []*grpc_client.ConfigGrpc
}

func NewEnv(env any) {
	setting := config.DefaultSettingsConfig()
	if setting.IsProduction() {
		setting.SetFile("prod.config")
	} else {
		setting.SetFile("dev.config")
	}
	config.NewConfig(setting, env)
}

func (env *Env) IsProduction() bool {
	return strings.ToLower(env.NODE_ENV) == "production"
}
