package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"Api/controllers"
	"Api/responses"
	"Lib/db"
	"Lib/models"
	"Lib/rediswrapper"

	photopb "github.com/fabrizio2210/photobook"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

const jpegImage = `
/9j/4AAQSkZJRgABAQEBLAEsAAD//gATQ3JlYXRlZCB3aXRoIEdJTVD/2wBDAP//////////////
////////////////////////////////////////////////////////////////////////2wBD
Af//////////////////////////////////////////////////////////////////////////
////////////wgARCAAKAAoDAREAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAH/xAAUAQEA
AAAAAAAAAAAAAAAAAAAA/9oADAMBAAIQAxAAAAGgA//EABQQAQAAAAAAAAAAAAAAAAAAACD/2gAI
AQEAAQUCH//EABQRAQAAAAAAAAAAAAAAAAAAACD/2gAIAQMBAT8BH//EABQRAQAAAAAAAAAAAAAA
AAAAACD/2gAIAQIBAT8BH//EABQQAQAAAAAAAAAAAAAAAAAAACD/2gAIAQEABj8CH//EABQQAQAA
AAAAAAAAAAAAAAAAACD/2gAIAQEAAT8hH//aAAwDAQACAAMAAAAQkk//xAAUEQEAAAAAAAAAAAAA
AAAAAAAg/9oACAEDAQE/EB//xAAUEQEAAAAAAAAAAAAAAAAAAAAg/9oACAECAQE/EB//xAAUEAEA
AAAAAAAAAAAAAAAAAAAg/9oACAEBAAE/EB//2Q== `

func strPtr(str string) *string {
	return &str
}

func intPtr(i int64) *int64 {
	return &i
}

func TestTicketRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/new_photo", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var res responses.Response
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Regexp(t, regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), res.Data["ticket_id"])

}

func TestUidRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/uid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var res responses.Response
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Regexp(t, regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"), res.Data["uid"])

}

func TestGetEventRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("GET single event success", func(mt *mtest.T) {
		db.DB = mt.Client
		db.EventCollection = mt.Coll
		source := models.PhotoEvent{
			Author:      "author",
			Author_id:   "abc-123-abc",
			Description: "A description",
			Event:       "creation",
			Id:          "abc-123-123",
			Location:    "",
			Order:       1,
			Photo_id:    "abc-123",
			Timestamp:   1,
		}
		want := source
		want.Author_id = ""
		want.Location = "/static/resized/abc-123.jpg"
		first := mtest.CreateCursorResponse(1, "photobook.events", mtest.FirstBatch, bson.D{
			{Key: "Author", Value: source.Author},
			{Key: "Author_id", Value: source.Author_id},
			{Key: "Description", Value: source.Description},
			{Key: "Event", Value: source.Event},
			{Key: "Id", Value: source.Id},
			{Key: "Order", Value: source.Order},
			{Key: "Photo_id", Value: source.Photo_id},
			{Key: "Timestamp", Value: source.Timestamp},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
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

func TestPutSuccessEventRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	uuid.SetRand(rand.New(rand.NewSource(1)))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"guest": {"editor": "no", "nome": "normal user"}}`)
	}))
	controllers.GuestApiURL = srv.URL

	mt.Run("PUT single event success", func(mt *mtest.T) {
		db.DB = mt.Client
		db.EventCollection = mt.Coll
		db.StatusCollection = mt.Coll
		source := models.PhotoEvent{
			Author:      "author",
			Author_id:   "abc-123-abc",
			Description: "A description",
			Event:       "creation",
			Id:          "52fdfc07-2182-454f-963f-5f0f9a621d72",
			Location:    "",
			Order:       1,
			Photo_id:    "abc-123",
			Timestamp:   1,
		}
		want := source
		want.Description = "new_description"
		want.Event = "edit"
		want.Author_id = ""
		want.Timestamp = 3
		want.Location = "/static/resized/abc-123.jpg"
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		photoEvent := mtest.CreateCursorResponse(1, "photobook.events", mtest.FirstBatch, bson.D{
			{Key: "Author", Value: source.Author},
			{Key: "Author_id", Value: source.Author_id},
			{Key: "Description", Value: source.Description},
			{Key: "Event", Value: "random-id-number"},
			{Key: "Id", Value: source.Id},
			{Key: "Order", Value: source.Order},
			{Key: "Photo_id", Value: source.Photo_id},
			{Key: "Timestamp", Value: source.Timestamp},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors, photoEvent, killCursors)
		expectedSse, _ := json.Marshal(want)
		data := []byte(`{"description": "` + want.Description + `"}`)
		var redisMock redismock.ClientMock
		rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
		redisMock.ExpectIncr("events_count").SetVal(3)
		redisMock.ExpectPublish("sse", expectedSse).SetVal(0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/photo/abc-123?author_id="+source.Author_id, bytes.NewBuffer(data))
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

func TestPutSuccessEditorEventRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	uuid.SetRand(rand.New(rand.NewSource(1)))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"guest": {"editor": "sì", "nome": "editor name"}}`)
	}))
	controllers.GuestApiURL = srv.URL

	mt.Run("PUT single event success", func(mt *mtest.T) {
		db.DB = mt.Client
		db.EventCollection = mt.Coll
		db.StatusCollection = mt.Coll
		source := models.PhotoEvent{
			Author:      "author",
			Author_id:   "abc-123-abc",
			Description: "A description",
			Event:       "creation",
			Id:          "52fdfc07-2182-454f-963f-5f0f9a621d72",
			Location:    "",
			Order:       1,
			Photo_id:    "abc-123",
			Timestamp:   1,
		}
		want := source
		want.Description = "new_description"
		want.Event = "edit"
		want.Author_id = ""
		want.Timestamp = 3
		want.Location = "/static/resized/abc-123.jpg"
		first := mtest.CreateCursorResponse(1, "photobook.events", mtest.FirstBatch, bson.D{
			{Key: "Author", Value: source.Author},
			{Key: "Author_id", Value: source.Author_id},
			{Key: "Description", Value: source.Description},
			{Key: "Event", Value: "random-id-number"},
			{Key: "Id", Value: source.Id},
			{Key: "Order", Value: source.Order},
			{Key: "Photo_id", Value: source.Photo_id},
			{Key: "Timestamp", Value: source.Timestamp},
		})
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors, first, killCursors)
		expectedSse, _ := json.Marshal(want)
		data := []byte(`{"description": "` + want.Description + `"}`)
		var redisMock redismock.ClientMock
		rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
		redisMock.ExpectIncr("events_count").SetVal(3)
		redisMock.ExpectPublish("sse", expectedSse).SetVal(0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/photo/abc-123?author_id=editor", bytes.NewBuffer(data))
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

func TestPutForbiddenEventRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	uuid.SetRand(rand.New(rand.NewSource(1)))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"guest": {"editor": "no", "nome": "normal user"}}`)
	}))
	controllers.GuestApiURL = srv.URL

	mt.Run("PUT single event failure", func(mt *mtest.T) {
		db.DB = mt.Client
		db.EventCollection = mt.Coll
		db.StatusCollection = mt.Coll
		source := models.PhotoEvent{
			Author:      "author",
			Author_id:   "abc-123-abc",
			Description: "A description",
			Event:       "creation",
			Id:          "52fdfc07-2182-454f-963f-5f0f9a621d72",
			Location:    "",
			Order:       1,
			Photo_id:    "abc-123",
			Timestamp:   1,
		}
		want := source
		want.Description = "new_description"
		want.Event = "edit"
		want.Author_id = ""
		want.Timestamp = 3
		want.Location = "/static/resized/abc-123.jpg"
		first := mtest.CreateCursorResponse(1, "photobook.events", mtest.FirstBatch, bson.D{
			{Key: "Author", Value: source.Author},
			{Key: "Author_id", Value: source.Author_id},
			{Key: "Description", Value: source.Description},
			{Key: "Event", Value: "random-id-number"},
			{Key: "Id", Value: source.Id},
			{Key: "Order", Value: source.Order},
			{Key: "Photo_id", Value: source.Photo_id},
			{Key: "Timestamp", Value: source.Timestamp},
		})
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors, first, killCursors)
		expectedSse, _ := json.Marshal(want)
		data := []byte(`{"description": "` + want.Description + `"}`)
		var redisMock redismock.ClientMock
		rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
		redisMock.ExpectIncr("events_count").SetVal(3)
		redisMock.ExpectPublish("sse", expectedSse).SetVal(0)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/photo/abc-123?author_id=wrong_author", bytes.NewBuffer(data))
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
	})
}

func TestPostPhotoBeforeRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)
	ticket_id := "123-1234-123"
	want := &photopb.PhotoIn{
		AuthorId: strPtr("abc-123-abc"),
		Location: strPtr("/static/resized/52fdfc07-2182-454f-963f-5f0f9a621d72.jpg"),
		PhotoId:  strPtr("52fdfc07-2182-454f-963f-5f0f9a621d72"),
		Photo:    jpegImageBuf.Bytes(),
	}
	marshaledWant, _ := proto.Marshal(want)
	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *want.AuthorId)
		part, err := writer.CreateFormFile("file", "someimg.jpeg")
		if err != nil {
			t.Error(err)
		}

		part.Write(decodedImage)
		if err != nil {
			t.Error(err)
		}
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "expecting", []byte(want.GetPhotoId())).SetVal(1)
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "photo", marshaledWant).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:"+ticket_id, "metadata", "photo", "expecting").SetVal([]interface{}{nil, string(marshaledWant), *want.PhotoId})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("POST photo before metadata", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPostPhotoAfterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)
	ticket_id := "1223-2345-345"
	metadata := &photopb.PhotoIn{
		Author:      strPtr("author"),
		AuthorId:    strPtr("abc-123-abc"),
		Description: strPtr("A description"),
	}
	marshaledMetadata, _ := proto.Marshal(metadata)
	photo := &photopb.PhotoIn{
		AuthorId: strPtr("abc-123-abc"),
		Location: strPtr("/static/resized/52fdfc07-2182-454f-963f-5f0f9a621d72.jpg"),
		PhotoId:  strPtr("52fdfc07-2182-454f-963f-5f0f9a621d72"),
		Photo:    jpegImageBuf.Bytes(),
	}
	marshaledPhoto, _ := proto.Marshal(photo)
	want := &photopb.PhotoIn{
		Author:      metadata.Author,
		AuthorId:    metadata.AuthorId,
		Description: metadata.Description,
		Id:          strPtr("9566c74d-1003-4c4d-bbbb-0407d1e2c649"),
		Location:    photo.Location,
		Order:       intPtr(23),
		PhotoId:     photo.PhotoId,
		Photo:       photo.Photo,
		Timestamp:   intPtr(3),
	}
	marshaledWant, _ := proto.Marshal(want)
	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *want.AuthorId)
		part, err := writer.CreateFormFile("file", "someimg.jpeg")
		if err != nil {
			t.Error(err)
		}

		part.Write(decodedImage)
		if err != nil {
			t.Error(err)
		}
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "expecting", []byte(photo.GetPhotoId())).SetVal(1)
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "photo", marshaledPhoto).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:"+ticket_id, "metadata", "photo", "expecting").SetVal([]interface{}{string(marshaledMetadata), string(marshaledPhoto), *photo.PhotoId})
	redisMock.ExpectIncr("photos_count").SetVal(23)
	redisMock.ExpectIncr("events_count").SetVal(3)
	redisMock.ExpectLPush("in_photos", marshaledWant).SetVal(0)
	redisMock.ExpectDel("waiting_ticket:" + ticket_id).SetVal(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("POST photo after metadata", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPostPhotoAfterNotExpectedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)
	ticket_id := "1223-2345-345"
	metadata := &photopb.PhotoIn{
		Author:      strPtr("author"),
		AuthorId:    strPtr("abc-123-abc"),
		Description: strPtr("A description"),
	}
	marshaledMetadata, _ := proto.Marshal(metadata)
	photo := &photopb.PhotoIn{
		AuthorId: strPtr("abc-123-abc"),
		Location: strPtr("/static/resized/52fdfc07-2182-454f-963f-5f0f9a621d72.jpg"),
		PhotoId:  strPtr("52fdfc07-2182-454f-963f-5f0f9a621d72"),
		Photo:    jpegImageBuf.Bytes(),
	}
	marshaledPhoto, _ := proto.Marshal(photo)

	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *photo.AuthorId)
		part, err := writer.CreateFormFile("file", "someimg.jpeg")
		if err != nil {
			t.Error(err)
		}

		part.Write(decodedImage)
		if err != nil {
			t.Error(err)
		}
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "expecting", []byte(photo.GetPhotoId())).SetVal(1)
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "photo", marshaledPhoto).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:"+ticket_id, "metadata", "photo", "expecting").SetVal([]interface{}{string(marshaledMetadata), string(marshaledPhoto), "wrong-photo-ID"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("POST photo after metadata", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPutMetadataFirstRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	ticket_id := "1234-1234-1234"

	want := &photopb.PhotoIn{
		Author:      strPtr("author"),
		AuthorId:    strPtr("abc-123-abc"),
		Description: strPtr("A description"),
	}
	marshaledWant, _ := proto.Marshal(want)
	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *want.AuthorId)
		writer.WriteField("description", *want.Description)
		writer.WriteField("author", *want.Author)
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:1234-1234-1234", "metadata", marshaledWant).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:1234-1234-1234", "metadata", "photo", "expecting").SetVal([]interface{}{string(marshaledWant), nil, nil})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("PUT metadata photo first", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPutMetadataAfterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)
	ticket_id := "1234-1234-12345"
	photo := &photopb.PhotoIn{
		AuthorId: strPtr("abc-123-abc"),
		Id:       strPtr("52fdfc07-2182-454f-963f-5f0f9a621d72"),
		Location: strPtr("/static/resized/9566c74d-1003-4c4d-bbbb-0407d1e2c649.jpg"),
		PhotoId:  strPtr("9566c74d-1003-4c4d-bbbb-0407d1e2c649"),
		Photo:    jpegImageBuf.Bytes(),
	}
	marshaledPhoto, _ := proto.Marshal(photo)
	metadata := &photopb.PhotoIn{
		Author:      strPtr("author"),
		AuthorId:    strPtr("abc-123-abc"),
		Description: strPtr("A description"),
	}
	marshaledMetadata, _ := proto.Marshal(metadata)
	want := &photopb.PhotoIn{
		Author:      metadata.Author,
		AuthorId:    metadata.AuthorId,
		Description: metadata.Description,
		Id:          photo.Id,
		Location:    photo.Location,
		Order:       intPtr(23),
		PhotoId:     photo.PhotoId,
		Photo:       photo.Photo,
		Timestamp:   intPtr(3),
	}
	marshaledWant, _ := proto.Marshal(want)
	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *want.AuthorId)
		writer.WriteField("description", *want.Description)
		writer.WriteField("author", *want.Author)
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "metadata", marshaledMetadata).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:"+ticket_id, "metadata", "photo", "expecting").SetVal([]interface{}{string(marshaledMetadata), string(marshaledPhoto), *photo.PhotoId})
	redisMock.ExpectIncr("photos_count").SetVal(23)
	redisMock.ExpectIncr("events_count").SetVal(3)
	redisMock.ExpectLPush("in_photos", marshaledWant).SetVal(0)
	redisMock.ExpectDel("waiting_ticket:" + ticket_id).SetVal(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("PUT photo", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPutMetadataAfterNotExpectedIdRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)
	ticket_id := "1234-1234-12345"
	photo := &photopb.PhotoIn{
		AuthorId: strPtr("abc-123-abc"),
		Id:       strPtr("52fdfc07-2182-454f-963f-5f0f9a621d72"),
		Location: strPtr("/static/resized/9566c74d-1003-4c4d-bbbb-0407d1e2c649.jpg"),
		PhotoId:  strPtr("9566c74d-1003-4c4d-bbbb-0407d1e2c649"),
		Photo:    jpegImageBuf.Bytes(),
	}
	marshaledPhoto, _ := proto.Marshal(photo)
	metadata := &photopb.PhotoIn{
		Author:      strPtr("author"),
		AuthorId:    strPtr("abc-123-abc"),
		Description: strPtr("A description"),
	}
	marshaledMetadata, _ := proto.Marshal(metadata)

	go func() {
		defer writer.Close()
		writer.WriteField("author_id", *metadata.AuthorId)
		writer.WriteField("description", *metadata.Description)
		writer.WriteField("author", *metadata.Author)
	}()
	var redisMock redismock.ClientMock
	rediswrapper.RedisClient, redisMock = redismock.NewClientMock()
	redisMock.ExpectHSet("waiting_ticket:"+ticket_id, "metadata", marshaledMetadata).SetVal(1)
	redisMock.ExpectHMGet("waiting_ticket:"+ticket_id, "metadata", "photo", "expecting").SetVal([]interface{}{string(marshaledMetadata), string(marshaledPhoto), "wrong-Photo-ID"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("PUT photo", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: false},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.NoError(t, redisMock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
}

func TestPostPhotoBlockedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	uuid.SetRand(rand.New(rand.NewSource(1)))

	ticket_id := "123-456-22"
	decodedImage := make([]byte, base64.StdEncoding.DecodedLen(len(jpegImage)))
	base64.StdEncoding.Decode(decodedImage, []byte(jpegImage))
	unflateImage, _, _ := image.Decode(bytes.NewReader(decodedImage))
	jpegImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(jpegImageBuf, unflateImage, nil)

	go func() {
		defer writer.Close()
		writer.WriteField("author_id", "abc-123-abc")
		writer.WriteField("description", "A description")
		writer.WriteField("author", "author")
		part, err := writer.CreateFormFile("file", "someimg.jpeg")
		if err != nil {
			t.Error(err)
		}

		part.Write(decodedImage)
		if err != nil {
			t.Error(err)
		}
	}()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/new_photo?ticket_id="+ticket_id, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	mt.Run("POST photo", func(mt *mtest.T) {
		db.DB = mt.Client
		db.StatusCollection = mt.Coll
		unblocked := mtest.CreateCursorResponse(1, "photobook.status", mtest.FirstBatch, bson.D{
			{Key: "id", Value: "block_upload"},
			{Key: "value", Value: true},
		})
		killCursors := mtest.CreateCursorResponse(0, "photobook.events", mtest.NextBatch)
		mt.AddMockResponses(unblocked, killCursors)

		router.ServeHTTP(w, req)
	})

	assert.Equal(t, 401, w.Code)
}
