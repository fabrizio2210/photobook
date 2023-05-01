package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func UidRoute(router *gin.Engine) {
  router.GET(   "/api/uid",         controllers.GetUid())
}
