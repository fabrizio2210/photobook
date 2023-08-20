package models

type CoverInputForm struct {
  Author_id   string `form:"author_id" validate:"required"`
}
