package rediswrapper

import (
  "context"
  "log"

  "github.com/go-redis/redis/v8"
)


var RedisClient *redis.Client
func ConnectRedis(address string) *redis.Client {
  log.Printf("Connecting to \"%s\" for Redis", address)
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

func HSet(key string, subkey string, value []byte) error{
  if err := RedisClient.HSet(ctx, key, subkey, value).Err(); err != nil {
    return err
  }
  return nil
}

func HMGet(key string, subkeys []string) ([]string, error){
  var vals []string
  res, err := RedisClient.HMGet(ctx, key, subkeys...).Result();
  if err != nil {
    return nil, err
  }
  for _, r := range res {
    if val, ok := r.(string); ok {
      vals = append(vals, val)
    }
  }
  return vals, nil
}

func HDel(key string, fields ...string) error {
  return RedisClient.HDel(ctx, key, fields...).Err()
}

func Del(key string) error {
  return RedisClient.Del(ctx, key).Err()
}

