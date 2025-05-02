package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// Lets start simple. Lets create struct for work log
// This structure should have: date, hours worked and description
type WorkDay struct {
	WorkDate    time.Time `json:"work_date"`
	HourWorked  int       `json:"hours_worked"`
	Description string    `json:"description"`
}


func updateResponseHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func main() {
	// Lets add some hard cored data to be accessed by our API
	workDays := []WorkDay{
		WorkDay{WorkDate: time.Date(2025, time.Month(5), 2, 0, 0, 0, 0, time.UTC), Description: "Second work day", HourWorked: 10},
		WorkDay{WorkDate: time.Date(2025, time.Month(5), 1, 0, 0, 0, 0, time.UTC), Description: "First work day", HourWorked: 9},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		updateResponseHeaders(w, r)
		data, err := json.Marshal(workDays)
		if err != nil {
			log.Fatal(err)
		}
		io.Writer.Write(w, data)
	})

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
