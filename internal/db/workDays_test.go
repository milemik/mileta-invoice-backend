package db

import (
	"context"
	"testing"
	"time"
)

type mockInsertResult struct{ id interface{} }

func (m mockInsertResult) InsertedID() interface{} { return m.id }

type mockCollection struct {
	inserted   []WorkDay
	shouldFail bool
}

func (m *mockCollection) InsertOne(ctx context.Context, doc interface{}) (interface{ InsertedID() interface{} }, error) {
	if m.shouldFail {
		return nil, context.DeadlineExceeded
	}
	wd, ok := doc.(WorkDay)
	if !ok {
		return nil, nil
	}
	m.inserted = append(m.inserted, wd)
	return mockInsertResult{id: "mockid"}, nil
}

func (m *mockCollection) Find(ctx context.Context, filter interface{}, _ ...interface{}) (CursorAPI, error) {
	if m.shouldFail {
		return nil, context.DeadlineExceeded
	}
	return &mockCursor{data: m.inserted}, nil
}

type mockCursor struct {
	data []WorkDay
	idx  int
}

func (c *mockCursor) Next(ctx context.Context) bool {
	if c.idx < len(c.data) {
		return true
	}
	return false
}

func (c *mockCursor) Decode(val interface{}) error {
	if c.idx >= len(c.data) {
		return context.DeadlineExceeded
	}
	ptr, ok := val.(*WorkDay)
	if !ok {
		return nil
	}
	*ptr = c.data[c.idx]
	c.idx++
	return nil
}

func (c *mockCursor) Close(ctx context.Context) error { return nil }
func (c *mockCursor) Err() error                      { return nil }

func TestAddWorkDayWithColl(t *testing.T) {
	mockColl := &mockCollection{}
	workDay := WorkDay{
		WorkDate:    time.Now(),
		HourWorked:  5,
		Description: "Mock test",
	}
	insertedID := AddWorkDayWithColl(mockColl, workDay)
	if insertedID != "mockid" {
		t.Errorf("Expected InsertedID to be 'mockid', got %v", insertedID)
	}
	if len(mockColl.inserted) != 1 {
		t.Errorf("Expected 1 inserted workday, got %d", len(mockColl.inserted))
	}
}

func TestGetWorkDaysWithColl(t *testing.T) {
	mockColl := &mockCollection{
		inserted: []WorkDay{
			{WorkDate: time.Now(), HourWorked: 8, Description: "A"},
			{WorkDate: time.Now(), HourWorked: 6, Description: "B"},
		},
	}
	workDays := GetWorkDaysWithColl(mockColl)
	if len(workDays) != 2 {
		t.Errorf("Expected 2 workdays, got %d", len(workDays))
	}
	if workDays[0].Description != "A" || workDays[1].Description != "B" {
		t.Errorf("Descriptions do not match: %+v", workDays)
	}
}
