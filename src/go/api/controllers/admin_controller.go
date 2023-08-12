package controllers

import (
  "context"
  "encoding/json"
  "net/http"
  "time"

  "Api/responses"
  "Lib/db"
  "Lib/models"
  "Lib/rediswrapper"

  "github.com/gin-gonic/gin"
)

func ToggleUpload() gin.HandlerFunc {
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
    if db.IsUploadBlocked() {
      db.UnblockUpload()
    } else {
      db.BlockUpload()
      encodedJson, err := json.Marshal(models.MessageEvent{
        Message: "Stop al televoto",
        Type: "stopped_upload",
        })
      if err != nil {
        panic(err)
      }
      // Notify all the clients.
      if err := rediswrapper.Publish("sse", encodedJson); err != nil {
          panic(err)
      }
    }
    c.JSON(
      http.StatusOK,
      responses.Response{
        Status: http.StatusOK,
        Message: "success",
        Data: map[string]interface{}{"upload_status": db.IsUploadBlocked()},
      },
    )
  }
}

