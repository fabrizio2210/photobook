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

type PhotoInputJson struct {
  Author      string `json:"author,omitempty"`
  Description string `json:"description,omitempty"`
}

type PhotoInputForm struct {
  Author      string `form:"author"`
  Author_id   string `form:"author_id" validate:"required"`
  Description string `form:"description"`
}

func (e *PhotoEvent) StripPrivateInfo() {
  e.Author_id = ""
}

type MessageEvent struct {
  Message string `json:"message"`
  Type string `json:"type"`
}
