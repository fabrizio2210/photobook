package main


import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "Lib/models"

    "github.com/go-redis/redis/v8"
    "github.com/alexandrevicenzi/go-sse"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
  Addr: os.Getenv("REDIS_HOST") + ":6379",
})

var notificationRoot = "/api/notifications/"

func main() {
  s := sse.NewServer(nil)
  defer s.Shutdown()
  http.Handle(notificationRoot, s)
  go func() {
    subscriber := redisClient.Subscribe(ctx, "sse")
    for {
      msg, err := subscriber.ReceiveMessage(ctx)
      fmt.Printf("Msg received:%v\n", msg.Payload)
      if err != nil {
        panic(err)
      }
      var m models.MessageEvent
      err = json.Unmarshal([]byte(msg.Payload), &m)
      if err != nil {
        panic(err)
      }
      if m.Channel != "" {
        m.Channel = notificationRoot + m.Channel
      }
      s.SendMessage(m.Channel, sse.NewMessage("", m.Message, m.Type))
    }
  }()
  log.Println("Listening at :3000")
  http.ListenAndServe(":3000", nil)
}

