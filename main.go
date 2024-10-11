package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"log/slog"
)

const (
	port         = 8080
	filePathRoot = "."
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetMiddlewareMetrics(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("/metrics", apiCfg.middlewareMetrics)
	mux.HandleFunc("/reset", apiCfg.resetMiddlewareMetrics)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
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
