package controllers

import (
  "bytes"
  "context"
  "encoding/json"
  "image"
  "image/jpeg"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"
  "time"

  "Api/filemanager"
  "Api/models"
  "Api/responses"
  "Api/rediswrapper"

  photopb "github.com/fabrizio2210/photobook"
  orientation "github.com/takumakei/exif-orientation"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "github.com/go-playground/validator/v10"
  "github.com/golang/protobuf/proto"
  "github.com/google/uuid"
  "github.com/nfnt/resize"
)

var EventCollection *mongo.Collection
var validate = validator.New()

func allowedExtensions(filename string) bool {
  exts := [3]string{"jpg", "jpeg","png"}
  result := false
  for _, ext := range exts {
    res := strings.HasSuffix(strings.ToLower(filename), ext)
    result = result || res
  }
  return result
}

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

func maybeGetJson(c *gin.Context, data *models.PhotoInputJson) bool {
  if err := c.BindJSON(&data); err != nil {
    log.Printf("Error in parsing json: %v", err.Error())
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

  if validationErr := validate.Struct(data); validationErr != nil {
    log.Printf("Error in validating json: %v", validationErr.Error())
    c.JSON(
      http.StatusBadRequest,
      responses.Response{
        Status: http.StatusBadRequest,
        Message: "error",
        Data: map[string]interface{}{"event": validationErr.Error()},
      },
    )
    return false
  }
  return true
}

func maybeGetForm(c *gin.Context, data *models.PhotoInputForm) bool {
  if err := c.Bind(&data); err != nil {
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

func blockUpload(c *gin.Context) bool {
  if os.Getenv("BLOCK_UPLOAD") != "" {
    c.JSON(
      http.StatusUnauthorized,
      responses.Response{
        Status: http.StatusUnauthorized,
        Message: os.Getenv("BLOCK_UPLOAD_MSG"),
      },
    )
    return true
  }
  return false
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
  if ctx.Value("write") == true {
    // Block if edit/deletion.
    if (blockUpload(c)) {
      return nil
    }
  }
  if ctx.Value("private") == true {
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
    context.WithValue(ctx, "write", false)
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
    context.WithValue(ctx, "write", true)
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
    defer cancel()

    var data models.PhotoInputJson
    if (! maybeGetJson(c, &data)) {
      return
    }

    context.WithValue(ctx, "private", true)
    context.WithValue(ctx, "write", true)
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

func GetAllPhotoEvents() gin.HandlerFunc {
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

func PostNewPhoto() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Validation.
    if (blockUpload(c)) {
      log.Printf("Upload denied by environment variable")
      return
    }
    var data models.PhotoInputForm
    if (! maybeGetForm(c, &data)) {
      log.Printf("Wrong parsing of the post data.")
      return
    }

    event_id := uuid.New()
    photo_id := uuid.New()
    event_id_str := event_id.String()
    photo_id_str := photo_id.String()
    order := rediswrapper.GetCounter("photos_count")
    timestamp := rediswrapper.GetCounter("events_count")
    location := filemanager.LocationForClient(photo_id_str)

    // Image processing and writing.
    form, _ := c.MultipartForm()
    files := form.File["file"]
    if (len(files) != 1) {
      log.Printf("Number of files is different from 1: %d", len(files))
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: "error: no file found or too many",
          Data: map[string]interface{}{"event": "no file found or too many"},
        },
      )
      return
    }
    file := files[0]
    if (! allowedExtensions(file.Filename)) {
      log.Printf("Bad extension: %s", file.Filename)
      c.JSON(
        http.StatusBadRequest,
        responses.Response{
          Status: http.StatusBadRequest,
          Message: "error: bad extension",
          Data: map[string]interface{}{"event": "bad extension"},
        },
      )
      return
    }

    fl, _ := file.Open()
    flRead, _ := ioutil.ReadAll(fl)
    originalImage, _, err := image.Decode(bytes.NewReader(flRead))
    if (err != nil) {
      log.Printf("Error in decoding the image: %v", err.Error())
    }

    o, _ := orientation.Read(bytes.NewReader(flRead))
    originalImage = orientation.Normalize(originalImage, o)
    originalImageBuf := bytes.NewBuffer([]byte{})
    jpeg.Encode(originalImageBuf, originalImage, nil)
    log.Printf("Writing in: %v", filemanager.PathToFullQualityFolder(photo_id_str))
    err = ioutil.WriteFile(
      filemanager.PathToFullQualityFolder(photo_id_str),
      originalImageBuf.Bytes(), os.ModePerm)
    if err != nil {
      log.Fatal(err)
    }
    resizedImage := resize.Thumbnail(900, 600, originalImage, resize.Lanczos3)
    imageBuf := bytes.NewBuffer([]byte{})
    jpeg.Encode(imageBuf, resizedImage, nil)
    log.Printf("Writing in: %v", filemanager.PathToUploadFolder(photo_id_str))
    err = ioutil.WriteFile(
      filemanager.PathToUploadFolder(photo_id_str), imageBuf.Bytes(), os.ModePerm)
    if err != nil {
      log.Fatal(err)
    }

    // Enque the photo for the worker.
    newPhoto := &photopb.PhotoIn{
      AuthorId: &data.Author_id,
      Id: &event_id_str,
      PhotoId: &photo_id_str,
      Author: &data.Author,
      Description: &data.Description,
      Timestamp: &timestamp,
      Order: &order,
      Location: &location,
      Photo: imageBuf.Bytes(),
    }
    marshalledNewPhoto, err := proto.Marshal(newPhoto)
    if err != nil {
        log.Fatalln("Failed to encode address book:", err)
    }
    rediswrapper.Enque("in_photos", marshalledNewPhoto)


  }
}
