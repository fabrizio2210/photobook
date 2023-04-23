package main


import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/go-redis/redis/v8"
    "github.com/alexandrevicenzi/go-sse"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
  Addr: os.Getenv("REDIS_HOST") + ":6379",
})


func main() {
  s := sse.NewServer(nil)
  defer s.Shutdown()
  http.Handle("/api/notifications", s)
  go func() {
    subscriber := redisClient.Subscribe(ctx, "sse")
    for {
      msg, err := subscriber.ReceiveMessage(ctx)
      fmt.Printf("Msg received:%v\n", msg.Payload)
      if err != nil {
        panic(err)
      }
      s.SendMessage("/api/notifications", sse.NewMessage("", msg.Payload, "photo"))
    }
  }()
  log.Println("Listening at :3000")
  http.ListenAndServe(":3000", nil)
}

