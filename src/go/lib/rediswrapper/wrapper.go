package rediswrapper

import (
  "context"

  "github.com/go-redis/redis/v8"
)


var RedisClient *redis.Client
func ConnectRedis(address string) *redis.Client {
  return redis.NewClient(&redis.Options{
    Addr: address,
  })
}

var ctx = context.Background()

func GetCounter(counter string) int64{
  msg, err := RedisClient.Incr(ctx, counter).Result()
  if err != nil {
    panic(err)
  }
  return msg
}

func Publish(topic string, json []byte) error{
  if err := RedisClient.Publish(ctx, topic, json).Err(); err != nil {
    panic(err)
  }
  return nil
}

func Enque(queue string, data []byte) error{
  if err := RedisClient.LPush(ctx, queue, data).Err(); err != nil {
    panic(err)
  }
  return nil
}

func WaitFor(queue string) ([]string, error){
  msg, err := RedisClient.BLPop(ctx, 0, queue).Result()
  if err != nil {
    panic(err)
  }
  return msg, nil
}
