package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// var profileCollection *mongo.Collection = *mongo.Client.Database("kimchi").Collection("profile")

// Profile - yet to be defined
type Profile struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Age       int                `json:"age,omitempty" bson:"age,omitempty"`
}

// CreateProfile - create a single profile
func CreateProfile(resp http.ResponseWriter, request *http.Request) {
	resp.Header().Add("content-type", "application/json")
	var profile Profile
	err := json.NewDecoder(request.Body).Decode(&profile)

	profileCollection := client.Database("kimchi").Collection("profile")

	if err != nil {
		log.Panicf("Could not decode request: %v", err)
		return
	}
	insertResult, err := profileCollection.InsertOne(context.TODO(), profile)
	if err != nil {
		log.Panicf("Failed to insert document into collection: %v", err)
	} else {
		log.Printf("Successfully inserted document into collection, id: %v", insertResult.InsertedID)
		json.NewEncoder(resp).Encode(insertResult)
	}
}

// GetAllProfiles - return all documents from profile collection - TODO finish up
func GetAllProfiles(resp http.ResponseWriter, request *http.Request) {
	profileCollection := client.Database("kimchi").Collection("profile")
	resp.Header().Add("context-type", "application/json")
	cursor, err := profileCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Panicf("Could not get db objects: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{ "message"` + err.Error() + `"}`))
		return
	} else {
		resp.WriteHeader(http.StatusOK)

	}
	defer cursor.Close(context.TODO())

	var profiles []Profile
	for cursor.Next(context.TODO()) {
		var profile Profile
		cursor.Decode(&profile)
		profiles = append(profiles, profile)
	}
	if err := cursor.Err(); err != nil {
		log.Panicf("Could not parse db objects: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{ "message"` + err.Error() + `"}`))
	}
	json.NewEncoder(resp).Encode(profiles)
}

func main() {
	log.Println("Start of program execution")
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	err := client.Connect(ctx)

	defer client.Disconnect(context.TODO())

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
	router.HandleFunc("/profiles", CreateProfile).Methods("POST")
	router.HandleFunc("/profiles", GetAllProfiles).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
