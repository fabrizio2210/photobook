package main

import (
    "context"
    "Api/routes"
    "os"

    "github.com/go-redis/redis/v8"
    "github.com/gin-gonic/gin"
)


var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_HOST") + ":6379",
})


func main() {
  router := gin.Default()
  routes.PhotoRoute(router)

  router.Run("0.0.0.0:5000")
}
