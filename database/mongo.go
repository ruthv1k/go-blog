package mongoconnect

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func InitializeDbConnection(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Println("Error connecting to db. Err: ", err.Error())
		return client, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Unable to ping to primary. Err: ", err.Error())
		return client, err
	}

	log.Println("Connected to db.")

	return client, nil
}

func DisconnectDb() {
    if client == nil {
        log.Fatal("Error accessing db client, db client is nil")
        return
    }

	if err := client.Disconnect(context.TODO()); err != nil {
        log.Fatal(err.Error())
		return
	}

	log.Println("Disconnecting from db.")
}

func GetClient() *mongo.Client {
    if client != nil {
        return client
    }

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	MONGO_URI := os.Getenv("MONGODB_URI")

	client, dbErr := InitializeDbConnection(MONGO_URI)

	if dbErr != nil {
        log.Fatal(dbErr.Error())
	}

	return client
}

func GetDatabase() *mongo.Database {
    if client != nil {
        return client.Database("relay_development")
    }

	return GetClient().Database("relay_development")
}