package controllers

import (
  "context"
  "log"
  "net/http"
  "time"

  "Api/db"
  "Api/models"
  "Api/responses"

  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "github.com/go-playground/validator/v10"
)

var eventCollection *mongo.Collection = db.GetCollection("events")
var validate = validator.New()

func EditPhoto() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var data models.PhotoEdit
    defer cancel()

    if err := c.BindJSON(&data); err != nil {
      c.JSON(
        http.StatusBadRequest,
        responses.EventResponse{
          Status: http.StatusBadRequest,
          Message: "error",
          Data: map[string]interface{}{"data": err.Error()},
        },
      )
      return
    }

    if validationErr := validate.Struct(&data); validationErr != nil {
      c.JSON(
        http.StatusBadRequest,
        responses.EventResponse{
          Status: http.StatusBadRequest,
          Message: "error",
          Data: map[string]interface{}{"data": validationErr.Error()},
        },
      )
      return
    }

    var photo models.PhotoEvent
    photoId := c.Param("photoId")
    err := eventCollection.FindOne(ctx, bson.M{"photo_id": photoId}).Decode(&photo)
    if err != nil {
      c.JSON(
        http.StatusNotFound,
        responses.EventResponse{
          Status: http.StatusNotFound,
          Message: "error", Data: map[string]interface{}{"data": err.Error()},
        },
      )
      return
    }

    log.Printf("AuthorID:%s, author_id:%v", photo.Author_id, c.Query("author_id"))
    if c.Query("author_id") != photo.Author_id {
      c.JSON(
        http.StatusUnauthorized,
        responses.EventResponse{
          Status: http.StatusUnauthorized,
          Message: "Not authorized",
        },
      )
      return
    }

    c.JSON(
      http.StatusOK,
      responses.EventResponse{
        Status: http.StatusOK,
        Message: "success",
        Data: map[string]interface{}{"data": photo},
      },
    )
  }
}
