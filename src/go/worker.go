package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "bytes"
    "net/http"
    "io/ioutil"

    "github.com/go-redis/redis/v8"
    "github.com/golang/protobuf/proto"
    photopb "github.com/fabrizio2210/photobook"
)

type NudityResponse struct {
    Value int
    Description  string
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

  fmt.Println(string(body))
  var nudity_response NudityResponse
  if err := json.Unmarshal(body, &nudity_response); err != nil {
      panic(err)
  }
  fmt.Printf("Response:%v\n", nudity_response)
  return int(nudity_response.Value) > 3
}

func main() {
    fmt.Println("NUDITY_APILAYER_KEY:", os.Getenv("NUDITY_APILAYER_KEY"))


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
        fmt.Printf("Photo length: %+v\n", len(photo_in.Photo))
	// fmt.Println("IsNudity:", isNudity(photo_in.Photo))
    }
}
