package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/google/uuid"
	"github.com/jeven2016/mylibs/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"reflect"
	"time"
)

const RedisStreamDataVar = "data"

type Redis struct {
	Client *redis.Client
	config *config.RedisConfig
}

func NewRedis(ctx context.Context, redisCfg *config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Address,
		Password:     redisCfg.Password,
		DB:           redisCfg.DefaultDb,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  time.Duration(redisCfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(redisCfg.WriteTimeout) * time.Second,
		PoolSize:     redisCfg.PoolSize,
		PoolTimeout:  time.Duration(redisCfg.PoolTimeout) * time.Second,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	rd := &Redis{
		Client: client,
		config: redisCfg,
	}
	return rd, nil
}

func (rd *Redis) EnsureConsumeGroupCreated(ctx context.Context, streamName string, group string) error {
	if groups, err := rd.Client.XInfoGroups(context.Background(), streamName).Result(); err == nil {
		return nil
	} else {
		//当无法获取到group信息时，创建一个消费group
		for _, g := range groups {
			if g.Name == group {
				return nil
			}
		}

		//You can use the XGROUP CREATE command with MKSTREAM option, to create an empty stream
		//XGroupCreate 方法要求先有stream的存在才能创建group
		if err = rd.Client.XGroupCreateMkStream(ctx, streamName, group, "0").Err(); err != nil {
			return err
		}
	}
	return nil
}

func (rd *Redis) PublishMessage(ctx context.Context, data interface{}, streamName string) error {
	if data == nil {
		return errors.New(fmt.Sprintf("cannot publis empty data, stream is %v", streamName))
	}
	var json string
	var err error
	dataType := reflect.TypeOf(data)

	if dataType.Kind() == reflect.String {
		json = data.(string)
	} else {
		json, err = convertor.ToJson(data)
	}

	if err != nil {
		return errors.New(fmt.Sprintf("unable convert data into json, stream: %s", streamName))
	}

	//just send the json data into stream since it's too complicated to map a struct to map, there would be
	//different kind of exceptions that need to handle
	err = rd.Client.XAdd(ctx, &redis.XAddArgs{
		Stream:     streamName,
		NoMkStream: false, // * 默认false,当为false时,key不存在，会新建
		MaxLen:     10000, // * 指定stream的最大长度,当队列长度超过上限后，旧消息会被删除，只保留固定长度的新消息
		Approx:     false, // * 默认false,当为true时,模糊指定stream的长度
		ID:         "*",   // 消息 id，我们使用 * 表示由 redis 生成
		// MinID: "id",            // * 超过阈值，丢弃设置的小于MinID消息id【基本不用】
		// Limit: 1000,            // * 限制长度【基本不用】
		Values: map[string]string{
			RedisStreamDataVar: json,
		}}).Err()

	return err
}

func (rd *Redis) Consume(ctx context.Context, streamName string,
	consumerGroup string, msgChan chan<- interface{}) error {
	defer func() {
		zap.S().Info("closing redis stream")
		close(msgChan)
		if err := recover(); err != nil {
			zap.S().Errorf("an unexpected error occurs during fetching data form stream, %v", err)
		}
	}()

	prefix := uuid.New().String()[:8]
	consumerId := streamName + ":consumer:" + prefix
loop:
	for {
		select {
		case <-ctx.Done():
			zap.S().Info("stop fetching while context canceled ")
			break loop
		default:
			entries, err := rd.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerId,

				Streams: []string{streamName, ">"},
				Count:   1,
				Block:   0,
			}).Result()
			if err != nil {
				zap.S().Errorf("failed to handle XReadGroup, %v", err)
				return err
			}

			//如果系统退出，此处chan中的已经Ack过的消息将无法保证正确处理，会出现丢失 todo enhancement
			for i := 0; i < len(entries[0].Messages); i++ {
				messageID := entries[0].Messages[i].ID
				jsonData := entries[0].Messages[i].Values[RedisStreamDataVar]
				msgChan <- jsonData
				zap.L().Info("retrieve a message into channel")

				rd.Client.XAck(ctx, streamName, consumerGroup, messageID)
			}
		}
	}
	return nil
}

// Len returns the current stream length
func (rd *Redis) Len(ctx context.Context, streamName string) (int64, error) {
	streamLen, err := rd.Client.XLen(ctx, streamName).Result()
	if err != nil {
		return 0, err
	}
	return streamLen, err
}
