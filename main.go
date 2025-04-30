package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/szmktk/chirpy/internal/config"
	"github.com/szmktk/chirpy/internal/database"
	"github.com/szmktk/chirpy/internal/server"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	srv, err := server.NewServer(cfg, dbQueries, logger)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/app/", http.StripPrefix("/app", srv.MiddlewareMetricsInc(http.FileServer(http.Dir(cfg.FilePathRoot)))))
	mux.HandleFunc("GET /admin/metrics", srv.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", srv.HandlerReset)
	mux.HandleFunc("POST /api/users", srv.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users", srv.AuthMiddleware(srv.HandlerUpdateUser))
	mux.HandleFunc("POST /api/login", srv.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", srv.HandlerRefresh)
	mux.HandleFunc("POST /api/revoke", srv.HandlerRevoke)
	mux.HandleFunc("POST /api/chirps", srv.AuthMiddleware(srv.HandlerCreateChirp))
	mux.HandleFunc("GET /api/chirps/{chirpID}", srv.HandlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", srv.AuthMiddleware(srv.HandlerDeleteChirp))
	mux.HandleFunc("GET /api/chirps", srv.HandlerGetAllChirps)
	mux.HandleFunc("POST /api/polka/webhooks", srv.HandlerUpgradeUserWebhook)
	mux.HandleFunc("GET /api/healthz", srv.HandlerHealth)

	var server *http.Server
	server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	logger.Info("Starting the server", "port", cfg.Port, "server_dir", cfg.FilePathRoot)
	log.Fatal(server.ListenAndServe())
}
