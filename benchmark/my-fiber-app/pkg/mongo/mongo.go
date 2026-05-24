package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// uri := config.GetEnv("MONGO_URI", "mongodb://localhost:27017")
	// dbName := config.GetEnv("MONGO_DB", "mydb")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017")) //.SetMonitor(otelmongo.NewMonitor()))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
		panic("Failed to connect to MongoDB" + err.Error())
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping error:", err)
		panic("Failed to connect to MongoDB" + err.Error())
	}

	Client = client
	DB = client.Database("cotador")

	log.Println("✅ Connected to MongoDB")
}
