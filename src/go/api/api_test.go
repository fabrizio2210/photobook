package main

import (
  "encoding/json"
	"net/http"
	"net/http/httptest"
  "regexp"
	"testing"

  "Api/responses"
  "Api/controllers"
  "Api/models"
  "Api/db"

	"github.com/stretchr/testify/assert"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/integration/mtest"
)


func TestUidRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/uid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
  var res responses.Response
  json.Unmarshal(w.Body.Bytes(), &res)
  assert.Regexp(t, regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), res.Data["uid"])

}


func TestEventRoute(t *testing.T) {
	router := setupRouter()
  mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
  defer mt.Close()

  mt.Run("test", func(mt *mtest.T) {
    db.DB = mt.Client
    controllers.EventCollection = mt.Coll
    source := models.PhotoEvent{
      Author:"author",
      Author_id:"abc-123-abc",
      Description:"A description",
      Event:"creation",
      Id:"abc-123-123",
      Location:"",
      Order:1,
      Photo_id:"abc-123",
      Timestamp:1,
    }
    want := source
    want.Author_id = ""
    want.Location = "/static/resized/abc-123.jpg"
    first := mtest.CreateCursorResponse(1, "test.trainers", mtest.FirstBatch, bson.D{
      {Key: "Author", Value: source.Author},
      {Key: "Author_id", Value: source.Author_id},
      {Key: "Description", Value: source.Description},
      {Key: "Event", Value: source.Event},
      {Key: "Id", Value: source.Id},
      {Key: "Order", Value: source.Order},
      {Key: "Photo_id", Value: source.Photo_id},
      {Key: "Timestamp", Value: source.Timestamp},
    })
    killCursors := mtest.CreateCursorResponse(0, "test.trainers", mtest.NextBatch)
    mt.AddMockResponses(first, killCursors)


    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/photo/abc-123", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    d := json.NewDecoder(w.Body)
    d.UseNumber()
    var res responses.Response
    d.Decode(&res)
    jsonData, _ := json.Marshal(res.Data["event"])
    var event models.PhotoEvent
    json.Unmarshal(jsonData, &event)
    assert.EqualValues(t, want, event)
  })

}
