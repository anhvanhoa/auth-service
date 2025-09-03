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
	Db          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
}

type queue struct {
	Addr        string
	Db          int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Network     string
	Concurrency int
	Queues      map[string]int
}

type Env struct {
	NodeEnv string `mapstructure:"node_env"`

	UrlDb string `mapstructure:"url_db"`

	NameService   string `mapstructure:"name_service"`
	PortGrpc      int    `mapstructure:"port_grpc"`
	HostGrpc      string `mapstructure:"host_grpc"`
	IntervalCheck string `mapstructure:"interval_check"`
	TimeoutCheck  string `mapstructure:"timeout_check"`

	DbCache *dbCache `mapstructure:"db_cache"`

	SecretOtp string `mapstructure:"secret_otp"`

	Queue *queue `mapstructure:"queue"`

	JwtSecret *jwtSecret `mapstructure:"jwt_secret"`

	FrontendUrl string `mapstructure:"frontend_url"`

	MailServiceAddr string `mapstructure:"mail_service_addr"`

	GrpcClients []*grpc_client.ConfigGrpc `mapstructure:"grpc_clients"`
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
	return strings.ToLower(env.NodeEnv) == "production"
}
