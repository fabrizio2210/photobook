package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func UserInfoRoute(router *gin.Engine) {
  router.GET("/api/user_info/:userId", controllers.GetUserInfo())
}
