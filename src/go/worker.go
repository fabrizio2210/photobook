package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
    "net/http"
    "io/ioutil"

    "github.com/go-redis/redis/v8"
)

type Data struct {
    photo  string `json:"photo"`
}
type NudityResponse struct {
    description  string `json:"description"`
    value string `json:"value"`
}

var ctx = context.Background()

var  nudity_api_url = "https://api.apilayer.com/nudity_detection/upload"

var redisClient = redis.NewClient(&redis.Options{
    Addr: "redis:6379",
})

func isNudity(payload string) bool{
  client := &http.Client {}
  reader:= strings.NewReader(payload)
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
  nudity_response := NudityResponse{}
  if err := json.Unmarshal([]byte(body), &nudity_response); err != nil {
      panic(err)
  }
  value, err := strconv.ParseInt(nudity_response.value, 10, 0)
  return int(value) > 4
}

func main() {
    fmt.Println("NUDITY_APILAYER_KEY:", os.Getenv("NUDITY_APILAYER_KEY"))
    subscriber := redisClient.Subscribe(ctx, "in_photos")

    data := Data{}

    for {
        msg, err := subscriber.ReceiveMessage(ctx)
        if err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(msg.Payload), &data); err != nil {
            panic(err)
        }


        fmt.Println("Received message from " + msg.Channel + " channel.")
	fmt.Println("IsNudity:", isNudity(data.photo))
        fmt.Printf("%+v\n", data)
    }
}
