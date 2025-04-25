package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"
	"github.com/szmktk/chirpy/internal/config"
	"github.com/szmktk/chirpy/internal/database"

	"log/slog"
)

const (
	port         = 8080
	filePathRoot = "."
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	polkaKey       string
	tokenSecret    string
}

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
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       cfg.Platform,
		polkaKey:       cfg.PolkaKey,
		tokenSecret:    cfg.TokenSecret,
	}

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.authMiddleware(apiCfg.handlerUpdateUser))
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("POST /api/chirps", apiCfg.authMiddleware(apiCfg.handlerCreateChirp))
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.authMiddleware(apiCfg.handlerDeleteChirp))
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUserWebhook)
	mux.HandleFunc("GET /api/healthz", handlerHealth)

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
