package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world in go\"}"))
	})

	// CORS レスポンスヘッダーの追加
	c := cors.Default()
	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe(":8007", handler))
}
