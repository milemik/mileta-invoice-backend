package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world")
	})

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
