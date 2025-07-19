package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// CollectionAPI abstracts mongo.Collection for testing
type CollectionAPI interface {
	InsertOne(context.Context, interface{}) (interface{ InsertedID() interface{} }, error)
	Find(context.Context, interface{}, ...interface{}) (CursorAPI, error)
}

// CursorAPI abstracts mongo.Cursor for testing
type CursorAPI interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	Err() error
}

// WorkDay struct for work log
// This structure should have: date, hours worked and description
type WorkDay struct {
	WorkDate    time.Time `json:"work_date"`
	HourWorked  int       `json:"hours_worked"`
	Description string    `json:"description"`
}

// getWorkDaysCollection returns the work days collection from the database
func getWorkDaysCollection(client *mongo.Client) CollectionAPI {
	return &collectionAdapter{client.Database("db").Collection("workdays")}
}

// collectionAdapter adapts *mongo.Collection to CollectionAPI
type collectionAdapter struct {
	inner *mongo.Collection
}

func (c *collectionAdapter) InsertOne(ctx context.Context, doc interface{}) (interface{ InsertedID() interface{} }, error) {
	res, err := c.inner.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return insertResultAdapter{res}, nil
}

func (c *collectionAdapter) Find(ctx context.Context, filter interface{}, _ ...interface{}) (CursorAPI, error) {
	cursor, err := c.inner.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &cursorAdapter{cursor}, nil
}

type insertResultAdapter struct {
	inner *mongo.InsertOneResult
}

func (i insertResultAdapter) InsertedID() interface{} {
	return i.inner.InsertedID
}

type cursorAdapter struct {
	inner *mongo.Cursor
}

func (c *cursorAdapter) Next(ctx context.Context) bool {
	return c.inner.Next(ctx)
}

func (c *cursorAdapter) Decode(val interface{}) error {
	return c.inner.Decode(val)
}

func (c *cursorAdapter) Close(ctx context.Context) error {
	return c.inner.Close(ctx)
}

func (c *cursorAdapter) Err() error {
	return c.inner.Err()
}

// addWorkDay adds a work day to the database
func AddWorkDayWithColl(coll CollectionAPI, workDay WorkDay) interface{} {
	res, err := coll.InsertOne(context.TODO(), workDay)
	if err != nil {
		log.Fatal(err)
	}
	return res.InsertedID()
}

// AddWorkDay is the original API
func AddWorkDay(client *mongo.Client, workDay WorkDay) interface{} {
	return AddWorkDayWithColl(getWorkDaysCollection(client), workDay)
}

// getWorkDays returns all work days from the database
func GetWorkDaysWithColl(coll CollectionAPI) []WorkDay {
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

// GetWorkDays is the original API
func GetWorkDays(client *mongo.Client) []WorkDay {
	return GetWorkDaysWithColl(getWorkDaysCollection(client))
}
