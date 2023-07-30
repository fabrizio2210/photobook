package main

import (
    "os"

    "Api/controllers"
    "Api/routes"
    "Lib/db"
    "Lib/filemanager"
    "Lib/rediswrapper"

    "github.com/gin-gonic/gin"
)


func setupRouter() *gin.Engine {
  r := gin.Default()
  routes.PhotoRoute(r)
  routes.UidRoute(r)
  return r
}

func main() {
  filemanager.Init()
  rediswrapper.RedisClient = rediswrapper.ConnectRedis(os.Getenv("REDIS_HOST") + ":6379")
  db.DB = db.ConnectDB()
  controllers.GuestApiURL = os.Getenv("GUEST_API_URL")
  controllers.EventCollection = db.GetCollection("events")
  router := setupRouter()
  router.Run("0.0.0.0:5000")
}

