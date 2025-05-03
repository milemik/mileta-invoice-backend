package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/milemik/mileta-invoice-backend/config"
	"github.com/milemik/mileta-invoice-backend/internal/db"
)

// HttpMessage is a simple structure for returning HttpMessage messages
// This structure should be used for success and error messages
type HttpMessage struct {
	Message string `json:"message"`
}

// updateResponseHeaders sets the response headers for the HTTP response
func updateResponseHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	config := config.ApiConfig{Port: ":8080", MongoDBUri: "mongodb://localhost:27017"}
	mongoClient := db.ConnectToMongoDB(config.MongoDBUri)
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/api/work/list/", func(w http.ResponseWriter, r *http.Request) {
		updateResponseHeaders(w, r)
		workDays := db.GetWorkDays(mongoClient)
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
		var workDay db.WorkDay
		// get data from request
		err = json.Unmarshal(requestBody, &workDay)
		if err != nil {
			log.Fatal(err)
		}
		db.AddWorkDay(mongoClient, workDay)
		updateResponseHeaders(w, r)
		io.Writer.Write(w, requestBody)
	})

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(config.Port, nil)
	log.Fatal(err)
}
