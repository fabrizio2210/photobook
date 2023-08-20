package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func AdminRoute(router *gin.Engine) {
  router.POST("/api/admin/toggle_upload", controllers.ToggleUpload())
  router.POST("/api/admin/upload_cover", controllers.PostCover())
}
