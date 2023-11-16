package db

import (
	"context"
	"github.com/jeven2016/mylibs/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Mongo struct {
	Client *mongo.Client
	Db     *mongo.Database
	Config *config.MongoConfig
}

func NewMongo(ctx context.Context, config *config.MongoConfig) (*Mongo, error) {
	var mg Mongo
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Uri))
	if err != nil {
		return nil, err
	}

	// 检测MongoDB是否连接成功
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	mg.Client = client
	mg.Config = config

	//初始化全局Db
	if config.Database != "" {
		mg.Db = client.Database(config.Database)
	}

	return &mg, nil
}
