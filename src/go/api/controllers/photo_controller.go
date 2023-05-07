package controllers

import (
  "context"
  "encoding/json"
  "log"
  "net/http"
  "os"
  "time"

  "Api/filemanager"
  "Api/models"
  "Api/responses"
  "Api/rediswrapper"

  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "github.com/go-playground/validator/v10"
  "github.com/google/uuid"
)

var EventCollection *mongo.Collection
var validate = validator.New()


func returnEvent(c *gin.Context, event *models.PhotoEvent) {
  c.JSON(
    http.StatusOK,
    responses.Response{
      Status: http.StatusOK,
      Message: "success",
      Data: map[string]interface{}{"event": event},
    },
  )
}

func maybeGetPhoto(ctx context.Context, c *gin.Context) *models.PhotoEvent {
  var photo *models.PhotoEvent
  photoId := c.Param("photoId")
  opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
  err := EventCollection.FindOne(ctx,
    bson.M{"photo_id": photoId},
    opts).Decode(&photo)
  if err != nil {
    c.JSON(
      http.StatusNotFound,
      responses.Response{
        Status: http.StatusNotFound,
        Message: "error", Data: map[string]interface{}{"event": err.Error()},
      },
    )
    return nil
  }
  if ctx.Value("private") == true {
    // Block if edit/deletion.
    if os.Getenv("BLOCK_UPLOAD") != "" {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: os.Getenv("BLOCK_UPLOAD_MSG"),
        },
      )
    }
    // Do not authorize if is not the author.
    if c.Query("author_id") != photo.Author_id {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: "Not authorized",
        },
      )
      return nil
    }
  }
  return photo
}

func insertEventDBAndPublish(ctx context.Context, c *gin.Context, event *models.PhotoEvent) {
  id := uuid.New()
  event.Id = id.String()
  event.Timestamp = rediswrapper.GetCounter("events_count")

  EventCollection.InsertOne(ctx, event)

  // Preparing for the public audience.
  event.StripPrivateInfo()
  event.Location = filemanager.LocationForClient(event.Photo_id)

  encodedEvent, err := json.Marshal(event)
  if err != nil {
    panic(err)
  }
  rediswrapper.Publish("sse", encodedEvent)
  log.Printf("Photo edited/deleted:%v", event)

  returnEvent(c, event)
}

func GetPhotoLatestEvent() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    context.WithValue(ctx, "private", false)
    event := maybeGetPhoto(ctx, c)
    if event == nil {
      // Photo not found.
      return
    }
    event.Location = filemanager.LocationForClient(event.Photo_id)
    event.Author_id = ""
    returnEvent(c, event)
  }
}

func DeletePhoto() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    context.WithValue(ctx, "private", true)
    new_event := maybeGetPhoto(ctx, c)
    if new_event == nil {
      // Photo not found or not authorized.
      return
    }

    new_event.Event = "deletion"
    insertEventDBAndPublish(ctx, c, new_event)

  }
}

func EditPhoto() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    var data models.PhotoEdit
    defer cancel()

    if err := c.BindJSON(&data); err != nil {
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
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
        responses.Response{
          Status: http.StatusBadRequest,
          Message: "error",
          Data: map[string]interface{}{"event": validationErr.Error()},
        },
      )
      return
    }

    context.WithValue(ctx, "private", true)
    new_event := maybeGetPhoto(ctx, c)
    if new_event == nil {
      // Photo not found or not authorized.
      return
    }

    if data.Author != "" {
      new_event.Author = data.Author
    }
    if data.Description != "" {
      new_event.Description = data.Description
    }
    new_event.Event = "edit"
    
    insertEventDBAndPublish(ctx, c, new_event)
  }
}

func GetAllPhotEvents() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    filter := bson.D{}
    if c.Query("author_id") != "" {
      filter = bson.D{{"author_id", c.Query("author_id")}}
    }
    opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
    cursor, err := EventCollection.Find(ctx, filter, opts)
    if err != nil {
      c.JSON(
        http.StatusNotFound,
        responses.Response{
          Status: http.StatusNotFound,
          Message: "error", Data: map[string]interface{}{"event": err.Error()},
        },
      )
      return
    }
    events := []models.PhotoEvent{}
    if err = cursor.All(ctx, &events); err != nil {
      panic(err)
    }
    for index, event := range events {
      events[index].Location = filemanager.LocationForClient(event.Photo_id)
    }

    c.JSON(
      http.StatusOK,
      responses.Response{
        Status: http.StatusOK,
        Message: "success",
        Data: map[string]interface{}{"events": events},
      },
    )
  }
}
