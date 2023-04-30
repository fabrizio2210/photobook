package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func PhotoRoute(router *gin.Engine) {
  router.PUT(   "/api/photo/:photoId", controllers.EditPhoto())
  router.DELETE("/api/photo/:photoId", controllers.DeletePhoto())
  router.GET(   "/api/photo/:photoId", controllers.GetPhoto())
}
