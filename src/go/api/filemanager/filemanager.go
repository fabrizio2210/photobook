package filemanager

import (
  "Api/models"

  "log"
  "os"
)

var staticPathUrl = "/static/"
var fullQualityFolder = "/tmp/"
var uploadFolder = "/tmp/"

func GetFileName(id string) string {
  return id + ".jpg"
}

func LocationForClient(id string) string {
  return staticPathUrl + "/resized/" + GetFileName(id)
}

func PhotoToClient(event models.PhotoEvent) models.PhotoEvent {
  populatedEvent := event
  populatedEvent.Location = LocationForClient(event.Photo_id)

  return populatedEvent
}

func PathToFullQualityFolder(id string) string {
  return fullQualityFolder + GetFileName(id)
}

func PathToUploadFolder(id string) string {
  return uploadFolder + GetFileName(id)
}

func SetFullQualityFolder(path string) {
  fullQualityFolder = path + "/original/"
  if err := os.MkdirAll(fullQualityFolder, os.ModePerm); err != nil {
    log.Fatal(err)
  }
}

func SetUploadFolder(path string) {
  uploadFolder = path + "/resized/"
  if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
    log.Fatal(err)
  }
}
