package controllers

import (
  "context"
  "time"

  "Api/responses"

  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
)


func GetUid() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    id := uuid.New()

    c.JSON(
      http.StatusOK,
      responses.Response{
        Status: http.StatusOK,
        Message: "success",
        Data: map[string]interface{}{"uid": id.String()},
      },
    )
  }
}
