package controllers

import (
  "context"
  "encoding/json"
  "net/http"
  "time"

  "Api/db"
  "Api/filemanager"
  "Api/models"
  "Api/responses"
  "Api/rediswrapper"

  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "github.com/go-playground/validator/v10"
  "github.com/google/uuid"
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
          Data: map[string]interface{}{"event": err.Error()},
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
          Data: map[string]interface{}{"event": validationErr.Error()},
        },
      )
      return
    }

    var photo models.PhotoEvent
    photoId := c.Param("photoId")
    err := eventCollection.Find(ctx,
      bson.M{"photo_id": photoId}).SetSort(
      bson.D{{"timestamp", 1}}).SetLimit(1).Decode(&photo)
    if err != nil {
      c.JSON(
        http.StatusNotFound,
        responses.EventResponse{
          Status: http.StatusNotFound,
          Message: "error", Data: map[string]interface{}{"event": err.Error()},
        },
      )
      return
    }

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
    
    new_event := photo
    if data.Author != "" {
      new_event.Author = data.Author
    }
    if data.Description != "" {
      new_event.Description = data.Description
    }
    id := uuid.New()
    new_event.Id = id.String()
    new_event.Event = "edit"
    new_event.Timestamp = rediswrapper.GetCounter("events_count")

    eventCollection.InsertOne(ctx, new_event)

    // Preparing for the public audiance.
    new_event.Author_id = ""
    new_event.Location = filemanager.LocationForClient(new_event.Photo_id)

    encodedEvent, err := json.Marshal(new_event)
    if err != nil {
      panic(err)
    }
    rediswrapper.Publish("sse", encodedEvent)

    c.JSON(
      http.StatusOK,
      responses.EventResponse{
        Status: http.StatusOK,
        Message: "success",
        Data: map[string]interface{}{"event": new_event},
      },
    )
  }
}
