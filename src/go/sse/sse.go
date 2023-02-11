package main


import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/gofiber/adaptor/v2"
    "github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
  Addr: os.Getenv("REDIS_HOST") + ":6379",
})

type Client struct {
  name   string
  events chan string
}

func main() {
  app := fiber.New()
  app.Get("/api/events", adaptor.HTTPHandler(handler(dashboardHandler)))
  app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
  return http.HandlerFunc(f)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
  client := &Client{name: r.RemoteAddr, events: make(chan string, 10)}
  go listenToEvents(client)

  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  w.Header().Set("Content-Type", "text/event-stream")
  w.Header().Set("Cache-Control", "no-cache")
  w.Header().Set("Connection", "keep-alive")

  timeout := time.After(60 * time.Second)
  select {
  case ev := <-client.events:
    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    enc.Encode(ev)
    fmt.Fprintf(w, "data: %v\n\n", ev)
    fmt.Printf("data: %v\n", buf.String())
  case <-timeout:
    fmt.Fprintf(w, ": nothing to sent\n\n")
  }

  if f, ok := w.(http.Flusher); ok {
    f.Flush()
  }
}

func listenToEvents(client *Client) {

  subscriber := redisClient.Subscribe(ctx, "sse")

  for {
    msg, err := subscriber.ReceiveMessage(ctx)
    fmt.Printf("Msg received:%v\n", msg.Payload)
    if err != nil {
      panic(err)
    }
    client.events <- msg.Payload
  }
}
