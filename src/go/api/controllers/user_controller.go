package controllers

import (
  "context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "time"

  "Api/responses"

  "net/http"
  "github.com/gin-gonic/gin"
)


func GetUserInfo() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    uid := c.Param("userId")
    res, err := getGuest(ctx, uid)

    if err != nil {
      res.StatusCode = http.StatusBadGateway
      res.Status = fmt.Sprintf("Error: %v", err)
    }


    log.Printf("Response: %+v", res)
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      log.Printf("Error in reading body: %s", err)
    }
    var anyJson map[string]interface{}
    json.Unmarshal(body, &anyJson)
    log.Printf("Returned from guest api: %s", body)
    c.JSON(
      res.StatusCode,
      responses.Response{
        Status: res.StatusCode,
        Message: res.Status,
        Data: anyJson,
      },
    )
  }
}

