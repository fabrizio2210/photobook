package models

type PhotoEvent struct {
  Id          string `json:"id"`
  Description string `json:"description"`
  Photo_id    string `json:"photo_id"`
  Order       int64  `json:"order"`
  Author      string `json:"author"`
  Author_id    string `json:"author_id"`
  Event       string `json:"event"`
  Timestamp   int64  `json:"timestamp"`
  Location    string `json:"location"`
}

type PhotoEdit struct {
  Author      string `json:"author,omitempty"`
  Description string `json:"description,omitempty"`
}
