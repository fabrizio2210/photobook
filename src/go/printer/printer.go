package main

import (
  "context"
  "log"
  "os"

  "Printer/db"
  "Lib/models"
  "Lib/filemanager"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()
var EventCollection *mongo.Collection

func removePhoto(photos []models.PhotoEvent, event models.PhotoEvent) []models.PhotoEvent {
  for i := range photos {
    if photos[i].Photo_id == event.Photo_id {
      return append(photos[:i], photos[i+1:]...)
    }
  }
  return photos
}

func substitutePhoto(photos []models.PhotoEvent, event models.PhotoEvent) []models.PhotoEvent {
  for i := range photos {
    if photos[i].Photo_id == event.Photo_id {
      photos[i] = event
      return photos
    }
  }
  return photos
}
  

func convertEventsToLayout(events []models.PhotoEvent) []*[2]models.PhotoEvent {
  var photos []models.PhotoEvent
  layout := []*[2]models.PhotoEvent{}

  for _, event := range events {
    switch event.Event {
    case "creation":
     photos = append(photos, event)
    case "deletion":
     photos = removePhoto(photos, event)
    case "edit":
     photos = substitutePhoto(photos, event)
    }
  }

  var page *[2]models.PhotoEvent
  for i, event := range photos {
    log.Printf("%+v", event) 
    if i%2 == 0 {
      page = &[2]models.PhotoEvent{}
      layout = append(layout, page)
    }
    page[i%2] = event
  }
  return layout
}

func main() {
  db.DB = db.ConnectDB()
  EventCollection := db.GetCollection("events")
  filemanager.Init()

  opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
  cursor, err := EventCollection.Find(ctx, bson.D{}, opts)
  if err != nil {
    panic(err)
  }
  events := []models.PhotoEvent{}
  if err = cursor.All(ctx, &events); err != nil {
    panic(err)
  }


  layout := convertEventsToLayout(events)
  log.Printf("Layout:%+v", layout)
  for _, page := range layout{
    log.Printf("Page:%+v", *page)
  }
  printToPDF(os.Getenv("STATIC_FILES_PATH") + "/download.pdf", layout)
}
