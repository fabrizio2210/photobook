package filemanager

import (
  "Api/models"
)

var staticPathUrl = "/static/resized/"

func GetFileName(id string) string {
  return id + ".jpg"
}

func LocationForClient(id string) string {
  return staticPathUrl + GetFileName(id)
}

func PhotoToClient(event models.PhotoEvent) models.PhotoEvent {
  populatedEvent := event
  populatedEvent.Location = LocationForClient(event.Photo_id)

  return populatedEvent
}
