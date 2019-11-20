package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var profileCollection *mongo.Collection = client.Database("kimchi").Collection("profile")

// Profile - yet to be defined
type Profile struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

// CreateProfile - create a single profile
func CreateProfile(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")
	var profile Profile
	err := json.NewDecoder(request.Body).Decode(&profile)
	if err != nil {
		log.Panicf("Could not decode request: %v", err)
		return
	}
	insertResult, err := profileCollection.InsertOne(context.TODO(), 2*time.Second)
	if err != nil {
		log.Panicf("Failed to insert document into collection: %v", err)
	} else {
		log.Printf("Successfully inserted document into collection, id: %v", insertResult.InsertedID)
	}
}

// GetAllProfiles - return all documents from profile collection - TODO finish up
func GetAllProfiles(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("context-type", "application/json")
	// fetchResult, err := profileCollection.
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Client initialized")
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Client Connected")
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %v", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/profile", CreateProfile).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
