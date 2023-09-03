package controllers

import (
  "context"
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"
  "time"

  "Api/responses"
  "Lib/db"
  "Lib/filemanager"
  "Lib/models"
  "Lib/rediswrapper"

  "github.com/gin-gonic/gin"
)

func GetUpload() gin.HandlerFunc {
  return func(c *gin.Context) {
    _, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

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

func PostCover() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    var data models.CoverInputForm
    if (! maybeGetForm(c, &data)) {
      log.Printf("Wrong parsing of the post data.")
      return
    }
    if (! isEditor(ctx, data.Author_id)) {
      c.JSON(
        http.StatusUnauthorized,
        responses.Response{
          Status: http.StatusUnauthorized,
          Message: "Not authorized",
        },
      )
      return
    }


    // PDF writing.
    form, err := c.MultipartForm()
    if err != nil {
      log.Printf("No multipart form found: %v", err.Error())
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
    if (!strings.HasSuffix(strings.ToLower(file.Filename), ".pdf")) {
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
    location := filemanager.GetCoverLocation()
    log.Printf("Writing in: %v", location)
    err = ioutil.WriteFile(
      location,
      flRead, os.ModePerm)
    if err != nil {
      log.Fatal(err)
    }
  }
}
