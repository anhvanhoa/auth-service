package bootstrap

import (
	"github.com/anhvanhoa/service-core/bootstrap/db"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	q "github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/utils"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env    *Env
	DB     *pg.DB
	Log    *log.LogGRPCImpl
	Cache  cache.CacheI
	Queue  q.QueueClient
	Helper utils.Helper
}

func App() *Application {
	env := Env{}
	NewEnv(&env)
	logConfig := log.NewConfig()
	log := log.InitLogGRPC(logConfig, zapcore.DebugLevel, env.IsProduction())
	db := db.NewPostgresDB(db.ConfigDB{
		URL:  env.UrlDb,
		Mode: env.NodeEnv,
	})
	configRedis := cache.NewConfigCache(
		env.DbCache.Addr,
		env.DbCache.Password,
		env.DbCache.Db,
		env.DbCache.Network,
		env.DbCache.MaxIdle,
		env.DbCache.MaxActive,
		env.DbCache.IdleTimeout,
	)
	cache := cache.NewCache(configRedis)
	cfgQueue := q.NewDefaultConfig(
		env.Queue.Addr,
		env.Queue.Network,
		env.Queue.Password,
		env.Queue.Db,
		nil,
		5,
	)
	queue := q.NewQueueClient(cfgQueue)
	helper := utils.NewHelper()
	return &Application{
		Env:    &env,
		DB:     db,
		Log:    log,
		Cache:  cache,
		Queue:  queue,
		Helper: helper,
	}
}
