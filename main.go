package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log/slog"
)

const port = 8080

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	mux := http.NewServeMux()
	var server *http.Server
	server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	// commented out because the first exercise is about returning 404 response
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "OK")
	// })

	logger.Info("Starting the server", "port", port)
	log.Fatal(server.ListenAndServe())
}
