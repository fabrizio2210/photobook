package main

import (
    "context"
    "os"

    "Api/controllers"
    "Api/routes"
    "Api/db"

    "github.com/go-redis/redis/v8"
    "github.com/gin-gonic/gin"
)


var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_HOST") + ":6379",
})

func setupRouter() *gin.Engine {
  r := gin.Default()
  routes.PhotoRoute(r)
  routes.UidRoute(r)
  return r
}

func main() {
  db.DB = db.ConnectDB()
  controllers.EventCollection = db.GetCollection("events")
  router := setupRouter()
  router.Run("0.0.0.0:5000")
}
