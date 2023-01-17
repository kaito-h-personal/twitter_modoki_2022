package main

import (
	"log"
	"net/http"
)

func main() {
	indexFunc := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Write([]byte("{\"hello\": \"world in go!\"}"))
	}

	http.HandleFunc("/", indexFunc)
	log.Fatal(http.ListenAndServe(":8007", nil))
}
