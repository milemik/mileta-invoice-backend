package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// WorkDay struct for work log
// This structure should have: date, hours worked and description
type WorkDay struct {
	WorkDate    time.Time `json:"work_date"`
	HourWorked  int       `json:"hours_worked"`
	Description string    `json:"description"`
}

// getWorkDaysCollection returns the work days collection from the database
func getWorkDaysCollection(client *mongo.Client) *mongo.Collection {
	collection := client.Database("db").Collection("workdays")
	return collection
}

// addWorkDay adds a work day to the database
func AddWorkDay(client *mongo.Client, workDay WorkDay) interface{} {
	coll := getWorkDaysCollection(client)
	res, err := coll.InsertOne(context.TODO(), workDay)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID
}

// getWorkDays returns all work days from the database
func GetWorkDays(client *mongo.Client) []WorkDay {
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
