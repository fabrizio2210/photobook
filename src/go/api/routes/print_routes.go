package routes

import (
  "Api/controllers"

  "github.com/gin-gonic/gin"
)

func PrintRoute(router *gin.Engine) {
  router.POST(  "/api/new_print",      controllers.PostNewPrint())
}
