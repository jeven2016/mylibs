package system

import (
	"context"
	"errors"
	"github.com/jeven2016/mylibs/cache"
	"github.com/jeven2016/mylibs/config"
	"github.com/jeven2016/mylibs/db"
	"github.com/jeven2016/mylibs/log"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var closed bool
var closeLock sync.Mutex

type StartupParams struct {
	EnableMongodb bool
	EnableRedis   bool
	EnableEtcd    bool
	Config        *config.ServerConfig
	PreShutdown   func() error
	PostShutdown  func() error
}

func (s *StartupParams) Validate() error {
	if s.Config == nil {
		return errors.New("start server failed: params and params.Config must be set")
	}
	if s.Config.ApplicationName == "" {
		return errors.New("start server failed: application name is required")
	}
	return nil
}

func Startup(ctx context.Context, params *StartupParams) *System {
	if params == nil {
		panic("start server failed: the params in method Startup(ctx, params) is required")
	}

	if err := params.Validate(); err != nil {
		panic(err.Error())
	}

	// 创建一个全局的App
	sys := &System{}
	sys.Config = params.Config
	sys.startupParams = params

	// log初始化
	log.SetupLog(params.Config.ApplicationName, params.Config.LogSetting)

	if params.EnableRedis {
		// 初始化redis
		redisClient, err := cache.NewRedis(ctx, params.Config.Redis)
		if err != nil {
			zap.L().Error("failed to initialize for redis", zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			zap.L().Info("Connecting to redis successfully")
			sys.RedisClient = redisClient
		}
	}

	if params.EnableMongodb {
		// 初始化Mongodb
		if mongoClient, err := db.NewMongo(ctx, params.Config.Mongo); err != nil {
			zap.L().Error("failed to connect mongodb", zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			zap.L().Info("Connecting to mongodb successfully")
			sys.MongoClient = mongoClient
		}
	}

	//init a routine pool
	pool, err := ants.NewPool(params.Config.TaskPoolSetting.Capacity)
	if err != nil {
		zap.L().Error("unable to init a routine pool", zap.Error(err))
		shutdown(ctx, sys, params)
		return nil
	} else {
		zap.L().Info("task pool initialized successfully")
	}
	sys.TaskPool = pool

	if params.EnableEtcd {
		//submit a task to register this service
		if err = sys.RegisterService(params.Config); err != nil {
			zap.L().Error("failed to register service in etcd", zap.String("app", params.Config.ApplicationName), zap.Error(err))
			shutdown(ctx, sys, params)
			return nil
		} else {
			zap.L().Info("service registered in etcd", zap.String("app", params.Config.ApplicationName))
		}
	}

	zap.L().Info("server starts successfully")
	exitChan := make(chan os.Signal)

	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(exitChan, syscall.SIGTERM, syscall.SIGINT)

	err = sys.TaskPool.Submit(func() {
		<-exitChan
		shutdown(ctx, sys, params)
	})
	if err != nil {
		zap.L().Info("unable to submit a shutdown hook", zap.Error(err))
		return nil
	}

	if err = sys.TaskPool.Submit(func() {
		<-ctx.Done()
		zap.S().Info("context is canceled")
		shutdown(ctx, sys, params)
	}); err != nil {
		zap.L().Info("unable to submit a shutdown hook", zap.Error(err))
	}
	SetSystem(sys)
	return sys
}

func Stop(ctx context.Context) {
	shutdown(ctx, GetSystem(), GetSystem().startupParams)
}

func shutdown(ctx context.Context, sys *System, params *StartupParams) {
	closeLock.Lock()
	defer func() {
		closed = true
		closeLock.Unlock()
	}()

	if closed {
		return
	}

	zap.L().Info("server is shutting down")

	if params.PreShutdown != nil {
		zap.S().Warn("call PreShutdown hook before exiting")
		if err := params.PreShutdown(); err != nil {
			zap.L().Warn("an error occurs while calling shutdown hook", zap.Error(err))
		}
	}

	if sys.RedisClient != nil {
		if err := sys.RedisClient.Client.Close(); err != nil {
			zap.L().Warn("an error occurs while closing redis's connection", zap.Error(err))
		} else {
			zap.S().Info("redis connections closed")
		}
	}

	if sys.MongoClient != nil {
		if err := sys.MongoClient.Client.Disconnect(ctx); err != nil {
			zap.L().Warn("an error occurs while closing mongodb's connection", zap.Error(err))
		} else {
			zap.S().Info("mongodb connections closed")
		}
	}

	//if sys.ServiceRegister != nil {
	//	if sys.ServiceRegister != nil {
	//		if err := sys.ServiceRegister.Cancel(ctx); err != nil {
	//			zap.L().Warn("an error occurs while closing etcd's connection", zap.Error(err))
	//		}
	//	}
	//}
	//
	if sys.TaskPool != nil {
		sys.TaskPool.Release()
		zap.S().Info("task pool released")
	}

	if params.PostShutdown != nil {
		zap.S().Info("call post shutdown hook before exiting")
		if err := params.PostShutdown(); err != nil {
			zap.L().Warn("an error occurs while calling post shutdown hook", zap.Error(err))
		}
	}
	zap.L().Info("shutdown completed")
}
