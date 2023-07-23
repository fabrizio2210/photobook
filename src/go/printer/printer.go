package main

import (
  "context"
  "Printer/db"
  "Printer/models"
  "log"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()
var EventCollection *mongo.Collection


func convertEventsToLayout(events []models.PhotoEvent, layout [][2]models.PhotoEvent) {
  for _, event := range events {
   log.Printf("%+v", event) 
  }
}

func main() {
  db.DB = db.ConnectDB()
  EventCollection := db.GetCollection("events")

  opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
  cursor, err := EventCollection.Find(ctx, bson.D{}, opts)
  if err != nil {
    panic(err)
  }
  events := []models.PhotoEvent{}
  if err = cursor.All(ctx, &events); err != nil {
    panic(err)
  }

  layout := [][2]models.PhotoEvent{}

  convertEventsToLayout(events, layout)

  cs, err := EventCollection.Watch(ctx, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
  if err != nil {
    panic(err)
  }
  defer cs.Close(ctx)

  for {
    if ok := cs.Next(ctx) ; ok {
      next := cs.Current
      log.Printf("%v", next)
      // convertEventsToLayout([]models.PhotoEvent{next["fullDocument"]}, layout)
    }
  }
}

