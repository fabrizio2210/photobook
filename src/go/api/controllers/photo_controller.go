package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"Api/responses"
	"Lib/db"
	"Lib/filemanager"
	"Lib/models"
	"Lib/rediswrapper"

	photopb "github.com/fabrizio2210/photobook"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	orientation "github.com/takumakei/exif-orientation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GuestApiURL string
var validate = validator.New()

type GuestStruct struct {
	Editor string `json:"editor"`
	Name   string `json:"nome"`
}

type GuestApiResponse struct {
	Guest GuestStruct `json="guest"`
}

func truncateText(str string, length int) string {
	if length <= 0 {
		return ""
	}

	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}

var allowed_extensions = [3]string{"jpg", "jpeg", "png"}

func allowedExtensions(filename string) bool {
	result := false
	for _, ext := range allowed_extensions {
		res := strings.HasSuffix(strings.ToLower(filename), ext)
		result = result || res
	}
	return result
}

func returnEvent(c *gin.Context, event *models.PhotoEvent) {
	c.JSON(
		http.StatusOK,
		responses.Response{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"event": event},
		},
	)
}

func maybeGetJson(c *gin.Context, data *models.PhotoInputJson) bool {
	if err := c.BindJSON(&data); err != nil {
		log.Printf("Error in parsing json: %v", err.Error())
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			},
		)
		return false
	}

	if validationErr := validate.Struct(data); validationErr != nil {
		log.Printf("Error in validating json: %v", validationErr.Error())
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Error %s", validationErr.Error()),
			},
		)
		return false
	}
	return true
}

func blockUpload(c *gin.Context) bool {
	if os.Getenv("BLOCK_UPLOAD") != "" || db.IsUploadBlocked() {
		msg := os.Getenv("BLOCK_UPLOAD_MSG")
		if msg == "" {
			msg = "The admin inhibited the upload."
		}
		c.JSON(
			http.StatusUnauthorized,
			responses.Response{
				Status:  http.StatusUnauthorized,
				Message: msg,
			},
		)
		return true
	}
	return false
}

func getGuest(ctx context.Context, id string) (*http.Response, error) {
	httpClient := http.Client{}
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		GuestApiURL+"/"+id, nil)
	if err != nil {
		log.Printf("Error in creating request:%v", err)
		return &http.Response{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error in doing request: %v", err)
		return &http.Response{}, err
	}
	return res, nil
}

func isEditor(ctx context.Context, id string) bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := getGuest(ctx, id)
	if err != nil {
		return false
	}
	switch res.StatusCode {
	case http.StatusOK:
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, _ := io.ReadAll(res.Body)

		var guest_api_response GuestApiResponse
		if err := json.Unmarshal(body, &guest_api_response); err != nil {
			log.Printf("Error in unmarshaling the request from the API: %v", err)
			return false
		}
		log.Printf("%+v", guest_api_response.Guest)
		if guest_api_response.Guest.Editor == "sì" {
			log.Printf("%s (id: %s) is an editor.",
				guest_api_response.Guest.Name,
				id)
			return true
		}
		log.Printf("%s (id: %s) is NOT an editor.",
			guest_api_response.Guest.Name,
			id)

	case http.StatusNotFound:
		log.Printf("Guest id=%s not found at %s.", id, GuestApiURL)

	default:
		log.Printf("%v: %s", res.StatusCode, res.Body)
	}
	return false
}

func maybeGetPhoto(ctx context.Context, c *gin.Context) *models.PhotoEvent {
	if ctx.Value("write") == true {
		// Block if edit/deletion.
		if blockUpload(c) {
			return nil
		}
	}

	var photo *models.PhotoEvent
	photoId := c.Param("photoId")
	opts := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	err := db.EventCollection.FindOne(ctx,
		bson.M{"photo_id": photoId},
		opts).Decode(&photo)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			responses.Response{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("Error: %s", err.Error()),
			},
		)
		return nil
	}
	if ctx.Value("private") == true {
		// Do not authorize if is not the author or editor.
		if (!isEditor(ctx, c.Query("author_id"))) && c.Query("author_id") != photo.Author_id {
			c.JSON(
				http.StatusUnauthorized,
				responses.Response{
					Status:  http.StatusUnauthorized,
					Message: "Not authorized.",
				},
			)
			return nil
		}
	}
	return photo
}

func insertEventDBAndPublish(ctx context.Context, c *gin.Context, event *models.PhotoEvent) {
	id := uuid.New()
	event.Id = id.String()
	event.Timestamp = rediswrapper.GetCounter("events_count")

	db.EventCollection.InsertOne(ctx, event)

	// Preparing for the public audience.
	event.StripPrivateInfo()
	event.Location = filemanager.LocationForClient(event.Photo_id)

	encodedEvent, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}
	rediswrapper.Publish("sse", encodedEvent)
	log.Printf("Photo edited/deleted:%v", event)

	returnEvent(c, event)
}

func GetPhotoLatestEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		ctx = context.WithValue(ctx, "private", false)
		ctx = context.WithValue(ctx, "write", false)
		event := maybeGetPhoto(ctx, c)
		if event == nil {
			// Photo not found.
			return
		}
		event.Location = filemanager.LocationForClient(event.Photo_id)
		event.Author_id = ""
		returnEvent(c, event)
	}
}

func DeletePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		ctx = context.WithValue(ctx, "private", true)
		ctx = context.WithValue(ctx, "write", true)
		new_event := maybeGetPhoto(ctx, c)
		if new_event == nil {
			// Photo not found or not authorized.
			return
		}

		new_event.Event = "deletion"
		insertEventDBAndPublish(ctx, c, new_event)

	}
}

func EditPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var data models.PhotoInputJson
		if !maybeGetJson(c, &data) {
			return
		}

		ctx = context.WithValue(ctx, "private", true)
		ctx = context.WithValue(ctx, "write", true)
		new_event := maybeGetPhoto(ctx, c)
		if new_event == nil {
			// Photo not found or not authorized.
			return
		}

		if data.Author != "" {
			new_event.Author = data.Author
		}
		if data.Description != "" {
			new_event.Description = data.Description
		}
		new_event.Event = "edit"

		insertEventDBAndPublish(ctx, c, new_event)
	}
}

func GetAllPhotoEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		filter := bson.D{}
		if c.Query("author_id") != "" {
			// Editors can see all the photos.
			if !isEditor(ctx, c.Query("author_id")) {
				filter = bson.D{{"author_id", c.Query("author_id")}}
			}
		}
		opts := options.Find().SetSort(bson.D{{"timestamp", 1}})
		cursor, err := db.EventCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(
				http.StatusNotFound,
				responses.Response{
					Status:  http.StatusNotFound,
					Message: fmt.Sprintf("Error: %s", err.Error()),
				},
			)
			return
		}
		events := []models.PhotoEvent{}
		if err = cursor.All(ctx, &events); err != nil {
			panic(err)
		}
		for index, event := range events {
			events[index].Location = filemanager.LocationForClient(event.Photo_id)
		}

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"events": events},
			},
		)
	}
}

func GetNewPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		ticket_id := uuid.New()

		c.JSON(
			http.StatusOK,
			responses.Response{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"ticket_id": ticket_id.String()},
			},
		)
	}
}

func PutNewPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Validation.
		if blockUpload(c) {
			log.Printf("Upload denied by environment variable")
			return
		}

		ticket_id := ticketID(c)
		if ticket_id == "" {
			log.Println("No ticket_id specified")
			return
		}

		waiting_ticket := "waiting_ticket:" + ticket_id
		var data models.MetadataInputForm
		if !maybeGetForm(c, &data) {
			log.Printf("Wrong parsing of the post data.")
			return
		}

		truncateAuthor := truncateText(data.Author, 20)
		truncateDescription := truncateText(data.Description, 200)
		// Enque the photo for the worker.
		newPhoto := &photopb.PhotoIn{
			AuthorId:    &data.Author_id,
			Author:      &truncateAuthor,
			Description: &truncateDescription,
		}
		marshalledNewPhoto, err := proto.Marshal(newPhoto)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		}
		log.Printf("Put metadata in \"%s\"", waiting_ticket)
		err = rediswrapper.HSet(waiting_ticket, "metadata", marshalledNewPhoto)
		if err != nil {
			log.Printf("Error with enque the metadata of \"%s\" in the waiting list: %v", waiting_ticket, err.Error())
			c.JSON(
				http.StatusBadRequest,
				responses.Response{
					Status:  http.StatusInternalServerError,
					Message: "Error: not possible to store the request.",
				},
			)
			return
		}
		err = maybeEnquePhotoToWorker(c, waiting_ticket)
		if err != nil {
			log.Printf("Error with enque in the in_photo list: %v", err)
			return
		}
	}
}

func ticketID(c *gin.Context) string {
	ticket_id := c.Query("ticket_id")
	if ticket_id == "" {
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Error: no ticket_id found in the request query.",
			},
		)
	}
	return ticket_id
}

func fileFromForm(c *gin.Context) (*multipart.FileHeader, error) {
	// Image processing and writing.

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Error: no file found in the request.",
			},
		)
		return nil, fmt.Errorf("no multipart form found: %v", err.Error())
	}

	files := form.File["file"]
	if len(files) != 1 {
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: "Error: too many file found in the request.",
			},
		)
		return nil, fmt.Errorf("number of files is different from 1: %d", len(files))
	}
	file := files[0]
	if !allowedExtensions(file.Filename) {
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Error: the \"%s\" file has a bad extension. It should one of the following \"%s\"", file.Filename, allowed_extensions),
				Data:    map[string]interface{}{"event": "bad extension"},
			},
		)
		return nil, fmt.Errorf("bad extension: %s", file.Filename)
	}
	return file, nil
}

func saveOnDiskAnsReturn(c *gin.Context, file *multipart.FileHeader, photo_id_str string) (*bytes.Buffer, error) {
	fl, _ := file.Open()
	flRead, _ := io.ReadAll(fl)
	log.Printf("Writing in: %v", filemanager.PathToFullQualityFolder(photo_id_str))
	err := os.WriteFile(
		filemanager.PathToFullQualityFolder(photo_id_str),
		flRead, os.ModePerm)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to save the image.",
			},
		)
		return nil, err
	}

	originalImage, _, err := image.Decode(bytes.NewReader(flRead))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to decode the image.",
			},
		)
		return nil, fmt.Errorf("error in decoding the image: %v", err.Error())
	}
	o, _ := orientation.Read(bytes.NewReader(flRead))
	originalImage = orientation.Normalize(originalImage, o)
	originalImageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(originalImageBuf, originalImage, nil)
	resizedImage := resize.Thumbnail(900, 600, originalImage, resize.Lanczos3)
	imageBuf := bytes.NewBuffer([]byte{})
	jpeg.Encode(imageBuf, resizedImage, nil)
	log.Printf("Writing in: %v", filemanager.PathToUploadFolder(photo_id_str))
	err = os.WriteFile(
		filemanager.PathToUploadFolder(photo_id_str), imageBuf.Bytes(), os.ModePerm)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to save the image.",
			},
		)
		return nil, err
	}
	return imageBuf, nil
}

func putPhotoInWaitingList(c *gin.Context, newPhoto *photopb.PhotoIn, waiting_ticket string) error {
	marshalledNewPhoto, err := proto.Marshal(newPhoto)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to marshall the request.",
			},
		)
		return fmt.Errorf("failed to encode the photo proto:", err)
	}
	log.Printf("Put \"%s\" photo in \"%s\"", newPhoto.GetPhotoId(), waiting_ticket)
	err = rediswrapper.HSet(waiting_ticket, "photo", marshalledNewPhoto)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to store the request.",
			},
		)
		return fmt.Errorf("error with enque of the photo in the waiting list \"%s\": %v", waiting_ticket, err)
	}
	err = rediswrapper.Expire(waiting_ticket, time.Duration(1)*time.Hour)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to store the photo.",
			},
		)
		return fmt.Errorf("error while setting expiration for \"%s\": %v", waiting_ticket, err)
	}
	return nil
}

func PostNewPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		// Validation.
		if blockUpload(c) {
			log.Printf("Upload denied by environment variable")
			return
		}

		ticket_id := ticketID(c)
		if ticket_id == "" {
			log.Println("No ticket_id specified")
			return
		}

		waiting_ticket := "waiting_ticket:" + ticket_id
		photo_id := uuid.New()
		photo_id_str := photo_id.String()
		location := filemanager.LocationForClient(photo_id_str)

		log.Printf("Reserving photo_id %s in \"%s\"", photo_id_str, waiting_ticket)
		err := rediswrapper.HSet(waiting_ticket, "expecting", []byte(photo_id_str))
		if err != nil {
			log.Printf("Error in deleting \"%s\" for photo, but carrying on: %v", waiting_ticket, err)
		}

		var data models.PhotoInputForm
		if !maybeGetForm(c, &data) {
			log.Printf("Wrong parsing of the post data.")
			return
		}

		file, err := fileFromForm(c)
		if err != nil {
			log.Println("Error with the extaction of the file: %v", err)
			return
		}

		imageBuf, err := saveOnDiskAnsReturn(c, file, photo_id_str)
		if err != nil {
			log.Printf("Error in saving the photo \"%s\": %v:", photo_id_str, err)
			return
		}
		// Put the photo in the waiting list.
		newPhoto := &photopb.PhotoIn{
			AuthorId: &data.Author_id,
			Location: &location,
			PhotoId:  &photo_id_str,
			Photo:    imageBuf.Bytes(),
		}
		err = putPhotoInWaitingList(c, newPhoto, waiting_ticket)
		if err != nil {
			log.Printf("Error in putting the photo in waiting list: %v", err)
			return
		}
		err = maybeEnquePhotoToWorker(c, waiting_ticket)
		if err != nil {
			log.Printf("Error with enque photo in the in_photo list: %v", err)
			return
		}
	}
}

func maybeEnquePhotoToWorker(c *gin.Context, waiting_ticket string) error {
	values, err := rediswrapper.HMGet(waiting_ticket, []string{"metadata", "photo", "expecting"})
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to read the waiting list.",
			},
		)
		return err
	}
	if len(values) < 3 {
		log.Printf("%s does not have both metadata and photo to proceed.\n", waiting_ticket)
		return nil
	}
	metadata_pb := values[0]
	photo_pb := values[1]
	expecting_photo_id := values[2]
	photo := &photopb.PhotoIn{}
	if err := proto.Unmarshal([]byte(photo_pb), photo); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to read the photo from the waiting list.",
			},
		)
		return err
	}
	metadata := &photopb.PhotoIn{}
	if err := proto.Unmarshal([]byte(metadata_pb), metadata); err != nil {
		c.JSON(
			http.StatusBadRequest,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to read the metadata from the waiting list.",
			},
		)
		return err
	}

	if expecting_photo_id != *photo.PhotoId {
		log.Printf("The current photo ID \"%s\" was not expected (%s)", photo.GetPhotoId(), expecting_photo_id)
		return nil
	}

	event_id := uuid.New()
	event_id_str := event_id.String()
	order := rediswrapper.GetCounter("photos_count")
	timestamp := rediswrapper.GetCounter("events_count")

	// Enque the photo for the worker.
	newPhoto := &photopb.PhotoIn{
		AuthorId:    metadata.AuthorId,
		Id:          &event_id_str,
		PhotoId:     photo.PhotoId,
		Author:      metadata.Author,
		Description: metadata.Description,
		Timestamp:   &timestamp,
		Order:       &order,
		Location:    photo.Location,
		Photo:       photo.Photo,
	}
	marshalledNewPhoto, err := proto.Marshal(newPhoto)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to marshall the photo to the worker.",
			},
		)
		return fmt.Errorf("failed to encode photo proto:", err)
	}
	rediswrapper.Enque("in_photos", marshalledNewPhoto)
	err = rediswrapper.Del(waiting_ticket)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			responses.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error: not possible to enque the photo to the worker.",
			},
		)
		return fmt.Errorf("redis gave an error: %v", err)
	}
	log.Printf("Enqueued \"%s\" photo in \"in_photos\"", *photo.PhotoId)
	return nil
}
