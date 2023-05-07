package db

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    photopb "github.com/fabrizio2210/photobook"
)

func ConnectDB() *mongo.Client  {
    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
    if err != nil {
        log.Fatal(err)
    }

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    //ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB")
    return client
}

//Client instance
var DB *mongo.Client

func GetCollection(collectionName string) *mongo.Collection {
    collection := DB.Database(os.Getenv("DB_NAME")).Collection(collectionName)
    return collection
}

func AcceptPhoto(photo_in *photopb.PhotoIn){
  insertPhoto(photo_in, "events")
}

func DiscardPhoto(photo_in *photopb.PhotoIn){
  insertPhoto(photo_in, "discarded")
}

func insertPhoto(photo_in *photopb.PhotoIn, collection_name string) {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URL")))
  defer func() {
    if err = client.Disconnect(ctx); err != nil {
        panic(err)
    }
  }()
  collection := client.Database(os.Getenv("DB_NAME")).Collection(collection_name)
  id, err := collection.InsertOne(ctx,
                                  bson.D{
                                          {"id", photo_in.Id},
                                          {"photo_id", photo_in.PhotoId},
                                          {"order", photo_in.Order},
                                          {"author_id", photo_in.AuthorId},
                                          {"author", photo_in.Author},
                                          {"timestamp", photo_in.Timestamp},
                                          {"description", photo_in.Description},
                                          {"event", "creation"},
                                        })
  if  err != nil {
    panic(err)
  }
  fmt.Println(id)
}