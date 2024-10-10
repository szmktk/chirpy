package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log/slog"
)

const (
	port         = 8080
	filePathRoot = "."
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		// w.Write([]byte("OK"))
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	var server *http.Server
	server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	logger.Info("Starting the server", "port", port, "server_dir", filePathRoot)
	log.Fatal(server.ListenAndServe())
}
