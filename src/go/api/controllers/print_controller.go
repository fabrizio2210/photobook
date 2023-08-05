package controllers

import (
  "context"
  "net/http"
  "time"

  "Api/responses"
  "Lib/rediswrapper"

  "github.com/gin-gonic/gin"
)

func PostNewPrint() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if (! isEditor(ctx, c.Query("author_id"))) {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: "Not authorized",
        },
      )
      return
    }
    rediswrapper.Enque("in_print", []byte{1})
  }
}
