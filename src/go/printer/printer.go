package main

import (
  "context"
  "log"
  "os"
  "sort"

  "Printer/db"
  "Lib/models"
  "Lib/filemanager"
  "Lib/rediswrapper"

  "github.com/rwcarlsen/goexif/exif"
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
  
func collapseEventsToPhotos(events []models.PhotoEvent) []models.PhotoEvent{
  var photos []models.PhotoEvent

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
  // Add location
  for i := range photos {
    photos[i].Location = filemanager.PathToFullQualityFolder(photos[i].Photo_id)
  }
  return photos
}

func readExifTimestamp(location string) int64 {
  log.Printf("Reading EXIF from %s", location)
  f, err := os.Open(location)
	if err != nil {
		log.Printf("Impossible to open %s: %s", location, err)
	}
  x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Impossible to decode EXIF for %s: %s", location, err)
    return 0
	}
  tm, _ := x.DateTime()
  return tm.Unix()
}

func orderPhotos(photos []models.PhotoEvent) {
  // Populate timestamp based on EXIF
  for i := range photos {
    date := readExifTimestamp(photos[i].Location)
    if date != 0 { 
      photos[i].Timestamp = date
    }
  }
  // order by Timestamp
  sort.Slice(photos, func(i, j int) bool {
    return photos[i].Timestamp < photos[j].Timestamp
  })
}

func convertEventsToLayout(photos []models.PhotoEvent) []*[2]models.PhotoEvent {
  layout := []*[2]models.PhotoEvent{}
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
  rediswrapper.RedisClient = rediswrapper.ConnectRedis(os.Getenv("REDIS_HOST") + ":6379")
  db.DB = db.ConnectDB()
  EventCollection := db.GetCollection("events")
  filemanager.Init()

  for {
    _, err := rediswrapper.WaitFor("in_print")
    if err != nil {
        panic(err)
    }
    opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
    cursor, err := EventCollection.Find(ctx, bson.D{}, opts)
    if err != nil {
      panic(err)
    }
    events := []models.PhotoEvent{}
    if err = cursor.All(ctx, &events); err != nil {
      panic(err)
    }

    photos := collapseEventsToPhotos(events)
    orderPhotos(photos)
    layout := convertEventsToLayout(photos)
    log.Printf("Layout:%+v", layout)
    for _, page := range layout{
      log.Printf("Page:%+v", *page)
    }
    printToPDF(os.Getenv("STATIC_FILES_PATH") + "/download.pdf", layout)
  }
}
