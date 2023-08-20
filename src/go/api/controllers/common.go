package controllers

import (
  "log"
  "net/http"

  "Api/responses"

  "github.com/gin-gonic/gin"
)

func maybeGetForm(c *gin.Context, data any) bool {
  if err := c.Bind(data); err != nil {
    log.Printf("Error in parsing: %v", err.Error())
    c.JSON(
      http.StatusBadRequest,
      responses.Response{
        Status: http.StatusBadRequest,
        Message: "error",
        Data: map[string]interface{}{"event": err.Error()},
      },
    )
    return false
  }


  if err := validate.Struct(data); err != nil {
    log.Printf("Error in validation: %v", err.Error())
    c.JSON(
      http.StatusBadRequest,
      responses.Response{
        Status: http.StatusBadRequest,
        Message: "error",
        Data: map[string]interface{}{"event": err.Error()},
      },
    )
    return false
  }
  return true
}
