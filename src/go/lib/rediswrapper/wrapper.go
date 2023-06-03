package rediswrapper

import (
  "context"
  "os"

  "github.com/go-redis/redis/v8"
)


var redisClient = redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_HOST") + ":6379",
})

var ctx = context.Background()

func GetCounter(counter string) int64{
  msg, err := redisClient.Incr(ctx, counter).Result()
  if err != nil {
    panic(err)
  }
  return msg
}

func Publish(topic string, json []byte) error{
  if err := redisClient.Publish(ctx, topic, json).Err(); err != nil {
    panic(err)
  }
  return nil
}

func Enque(queue string, data []byte) error{
  if err := redisClient.LPush(ctx, queue, data).Err(); err != nil {
    panic(err)
  }
  return nil
}
