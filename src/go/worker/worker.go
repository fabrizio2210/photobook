package main


import (
    "context"
    "db"
    "encoding/json"
    "fmt"
    "os"
    "bytes"
    "math/rand"
    "net/http"
    "io/ioutil"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/golang/protobuf/proto"
    photopb "github.com/fabrizio2210/photobook"
)

type NudityResponse struct {
    Value int
    Description  string
}

type PhotoEvent struct {
  Id string `json:"id"`
  Description string `json:"description"`
  Photo_id string `json:"photo_id"`
  Order int64 `json:"order"`
  Author string `json:"author"`
  Event string `json:"event"`
  Timestamp int64 `json:"timestamp"`
  Location string `json:"location"`
}

var  nudity_api_url = "https://api.apilayer.com/nudity_detection/upload"

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
    Addr: os.Getenv("REDIS_HOST") + ":6379",
})

func isNudity(payload []byte) bool{
  client := &http.Client {}
  reader:= bytes.NewReader(payload)
  req, err := http.NewRequest("POST", nudity_api_url, reader)
  req.Header.Set("apikey", os.Getenv("NUDITY_APILAYER_KEY"))
  if err != nil {
    fmt.Println(err)
  }
  res, err := client.Do(req)
	if res.Body != nil {
    defer res.Body.Close()
  }
  body, err := ioutil.ReadAll(res.Body)

  var nudity_response NudityResponse
  if err := json.Unmarshal(body, &nudity_response); err != nil {
      panic(err)
  }
  fmt.Printf("Response:%v\n", nudity_response)
  return int(nudity_response.Value) > 3
}

func main() {
    fmt.Println("NUDITY_APILAYER_KEY:", os.Getenv("NUDITY_APILAYER_KEY"))

    rand.Seed(time.Now().UnixNano())

    for {
        msg, err := redisClient.BLPop(ctx, 0, "in_photos").Result()
        if err != nil {
            panic(err)
        }

        fmt.Printf("Received message from: %+v\n", msg[0])
        
        photo_in := &photopb.PhotoIn{}
        if err := proto.Unmarshal([]byte(msg[1]), photo_in); err != nil {
            panic(err)
        }
        fmt.Printf("Author: %s\n", *photo_in.AuthorId)
        fmt.Printf("Order: %s\n", *photo_in.Order)
        fmt.Printf("Photo length: %+v\n", len(photo_in.Photo))

        // Insted of using the API, emulate it with a sleep.
        // fmt.Println("IsNudity:", isNudity(photo_in.Photo))
        n := rand.Intn(7)
        time.Sleep(time.Duration(n)*time.Second)

        db.AcceptPhoto(photo_in)
        //db.DiscardPhoto(photo_in)
        // Notify all the clients.
        
        encodedJson, err := json.Marshal(PhotoEvent{
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
        fmt.Printf("json to send: %s\n", encodedJson)

        if err := redisClient.Publish(ctx, "sse", encodedJson).Err(); err != nil {
            panic(err)
        }
    }
}
