package bootstrap

import (
	"time"

	"github.com/anhvanhoa/service-core/boostrap/db"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	q "github.com/anhvanhoa/service-core/domain/queue"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap/zapcore"
)

type Application struct {
	Env   *Env
	DB    *pg.DB
	Log   *log.LogGRPCImpl
	Cache cache.CacheI
	Queue q.QueueClient
}

func App() *Application {
	env := Env{}
	NewEnv(&env)
	logConfig := log.NewConfig()
	log := log.InitLogGRPC(logConfig, zapcore.DebugLevel, env.IsProduction())
	db := db.NewPostgresDB(db.ConfigDB{
		URL:  env.URL_DB,
		Mode: env.NODE_ENV,
	})
	configRedis := cache.NewConfigCache(
		env.DB_CACHE.Addr,
		env.DB_CACHE.Password,
		env.DB_CACHE.DB,
		env.DB_CACHE.Network,
		env.DB_CACHE.MaxIdle,
		env.DB_CACHE.MaxActive,
		env.DB_CACHE.IdleTimeout,
	)
	cache := cache.NewCache(configRedis)
	cfgQueue := q.NewDefaultConfig(
		env.QUEUE.Addr,
		env.QUEUE.Network,
		env.QUEUE.Password,
		env.QUEUE.DB,
		time.Minute*2,
		nil,
		5,
	)
	queue := q.NewQueueClient(cfgQueue)
	return &Application{
		Env:   &env,
		DB:    db,
		Log:   log,
		Cache: cache,
		Queue: queue,
	}
}
