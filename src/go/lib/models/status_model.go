package models

type Status struct{
  Id string `bson:"id,omitempty"`
  Value interface{} `bson:"value,omitempty"`
}
