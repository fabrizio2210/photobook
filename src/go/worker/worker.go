package main

import (
    "context"
    "Lib/db"
    "Lib/models"
    "Lib/rediswrapper"
    "encoding/json"
    "log"
    "os"
    "bytes"
    "math/rand"
    "net/http"
    "io/ioutil"
    "time"

    "github.com/golang/protobuf/proto"
    photopb "github.com/fabrizio2210/photobook"
)

type NudityResponse struct {
    Value int
    Description  string
}

var  nudity_api_url = "https://api.apilayer.com/nudity_detection/upload"

var ctx = context.Background()

func isNudity(payload []byte) bool{
  client := &http.Client {}
  reader:= bytes.NewReader(payload)
  req, err := http.NewRequest("POST", nudity_api_url, reader)
  req.Header.Set("apikey", os.Getenv("NUDITY_APILAYER_KEY"))
  if err != nil {
    log.Printf("Error in creating the request to the nudity API server: %s", err)
    return false
  }
  res, err := client.Do(req)
  if err != nil {
    log.Printf("Error in running the request to the nudity API server: %s", err)
    return false
  }
	if res.Body != nil {
    defer res.Body.Close()
  }
  body, err := ioutil.ReadAll(res.Body)

  var nudity_response NudityResponse
  if err := json.Unmarshal(body, &nudity_response); err != nil {
      log.Printf("Error while parsing the json from the nudity API server: %s. Body: \"%s\"", err, body)
      return false
  }
  log.Printf("Response:%v\n", nudity_response)
  return int(nudity_response.Value) > 3
}

func main() {
  log.Printf("NUDITY_APILAYER_KEY:%v", os.Getenv("NUDITY_APILAYER_KEY"))
  db.DB = db.ConnectDB()
  rediswrapper.RedisClient = rediswrapper.ConnectRedis(os.Getenv("REDIS_HOST") + ":6379")

  rand.Seed(time.Now().UnixNano())

  for {
    msg, err := rediswrapper.WaitFor("in_photos")
    if err != nil {
        panic(err)
    }

    log.Printf("Received message from: %+v\n", msg[0])
    
    photo_in := &photopb.PhotoIn{}
    if err := proto.Unmarshal([]byte(msg[1]), photo_in); err != nil {
        panic(err)
    }
    log.Printf("Author: %s\n", *photo_in.AuthorId)
    log.Printf("Order: %v\n", *photo_in.Order)
    log.Printf("Photo length: %+v\n", len(photo_in.Photo))

    // Insted of using the API, emulate it with a sleep.
    nudityAnswer := isNudity(photo_in.Photo)
    log.Printf("IsNudity: %v", nudityAnswer)

    if (nudityAnswer) {
      db.DiscardPhoto(photo_in)
      wrappedJson, err := json.Marshal(models.MessageEvent{
        Message: "Your photo contained nudity, it was discarded.",
        Type: "error",
        Channel: *photo_in.AuthorId,
      })
      if err != nil {
        panic(err)
      }
      // Notify just that client.
      if err := rediswrapper.Publish("sse", wrappedJson); err != nil {
          panic(err)
      }
    } else {
      db.AcceptPhoto(photo_in)
      // Omitting AuthorId on purpose.
      encodedJson, err := json.Marshal(models.PhotoEvent{
        Id: *photo_in.Id,
        Description: *photo_in.Description,
        Photo_id: *photo_in.PhotoId,
        Order: *photo_in.Order,
        Author: *photo_in.Author,
        Event: "creation",
        Timestamp: *photo_in.Timestamp,
        Location: *photo_in.Location,
        })
      if err != nil {
        panic(err)
      }
      wrappedJson, err := json.Marshal(models.MessageEvent{
        Message: string(encodedJson),
        Type: "photo",
      })
      // Notify all the clients.
      if err := rediswrapper.Publish("sse", wrappedJson); err != nil {
          panic(err)
      }
    }
  }
}

