package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func PhotoRoute(router *gin.Engine) {
  router.PUT(   "/api/photo/:photoId", controllers.EditPhoto())
  router.DELETE("/api/photo/:photoId", controllers.DeletePhoto())
  router.GET(   "/api/photo/:photoId", controllers.GetPhotoLatestEvent())
  router.GET(   "/api/events",         controllers.GetAllPhotoEvents())
  router.POST(  "/api/new_photo",      controllers.PostNewPhoto())
  router.PUT(  "/api/new_photo",      controllers.PutNewPhoto())
  router.GET(  "/api/new_photo",      controllers.GetNewPhoto())
}
