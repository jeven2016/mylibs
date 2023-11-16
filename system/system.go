package system

import (
	"github.com/jeven2016/mylibs/cache"
	"github.com/jeven2016/mylibs/config"
	"github.com/jeven2016/mylibs/db"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var system *System
var sysOnce = sync.Once{}
var lock = sync.Mutex{}

type System struct {
	RedisClient *cache.Redis

	MongoClient *db.Mongo

	//ServiceRegister Register

	Config config.Config

	TaskPool *ants.Pool

	collectionMap map[string]*mongo.Collection

	startupParams *StartupParams
}

func (s *System) RegisterService(cfg *config.ServerConfig) error {
	//registerParam := &RegisterParam{
	//	ServiceName:           cfg.ApplicationName,
	//	ServiceHost:           cfg.Http.Address,
	//	ServicePort:           cfg.Http.Port,
	//	RefreshSeconds:        cfg.Registration.Etcd.RefreshSeconds,
	//	ConnectTimeoutSeconds: cfg.Registration.Etcd.ConnectTimeoutSeconds,
	//}
	//
	////register service to etcd
	//if register, err := NewRegister(cfg.Registration.Etcd.Endpoints, s.Log, registerParam, s.TaskPool); err != nil {
	//	return err
	//} else {
	//	s.ServiceRegister = register
	//	err = s.ServiceRegister.Register(context.Background())
	//
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

//func (s *System) GetServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
//	return s.ServiceRegister.ListServiceAddresses(ctx, serviceName)
//}

func (s *System) GetCollection(name string) *mongo.Collection {
	sysOnce.Do(func() {
		s.collectionMap = make(map[string]*mongo.Collection)
	})
	if collection, ok := s.collectionMap[name]; ok {
		return collection
	} else {
		lock.Lock()
		defer lock.Unlock()

		if collection, ok = s.collectionMap[name]; ok {
			return collection
		}
		s.collectionMap[name] = s.MongoClient.Db.Collection(name)
	}
	return s.collectionMap[name]
}

func GetSystem() *System {
	return system
}

func SetSystem(sys *System) {
	system = sys
}
