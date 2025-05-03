package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Lets start simple. Lets create struct for work log
// This structure should have: date, hours worked and description
type WorkDay struct {
	WorkDate    time.Time `json:"work_date"`
	HourWorked  int       `json:"hours_worked"`
	Description string    `json:"description"`
}


type HttpMessage struct {
	Message string `json:"message"`
}


func updateResponseHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}


func connectToMongoDB() *mongo.Client {
	uri := "mongodb://localhost:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	return client
}

func getWorkDaysCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("db").Collection("workdays")
	return collection
}


func addWorkDay(client *mongo.Client, workDay WorkDay) interface{} {
	coll := getWorkDaysCollection(client)
	res, err := coll.InsertOne(context.TODO(), workDay)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID
}

func getWorkDays(client *mongo.Client) []WorkDay {
	coll := getWorkDaysCollection(client)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())
	var workDays []WorkDay
	for cursor.Next(context.TODO()) {
		var workDay WorkDay
		if err := cursor.Decode(&workDay); err != nil {
			log.Fatal(err)
		}
		workDays = append(workDays, workDay)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return workDays
}



func main() {
	mongoClient := connectToMongoDB()
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	

	http.HandleFunc("/api/work/list/", func(w http.ResponseWriter, r *http.Request) {
		updateResponseHeaders(w, r)
		workDays := getWorkDays(mongoClient)
		data, err := json.Marshal(workDays)
		if err != nil {
			log.Fatal(err)
		}
		io.Writer.Write(w, data)
	})


	http.HandleFunc("/api/work/add/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			message, err := json.Marshal(HttpMessage{Message: "method not allowed"})
			if err != nil {
				log.Fatal(err)
			}
			io.Writer.Write(w, message)
		}
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()
		var workDay WorkDay
		// get data from request
		err = json.Unmarshal(requestBody, &workDay)
		if err != nil {
			log.Fatal(err)
		}
		addWorkDay(mongoClient, workDay)
		updateResponseHeaders(w, r)
		io.Writer.Write(w, requestBody)
	})

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
