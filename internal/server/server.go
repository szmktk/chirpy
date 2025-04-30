package server

import (
	"log/slog"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/szmktk/chirpy/internal/config"
	"github.com/szmktk/chirpy/internal/database"
)

type Server struct {
	cfg            *config.Config
	db             *database.Queries
	logger         *slog.Logger
	fileserverHits atomic.Int32
}

func NewServer(cfg *config.Config, db *database.Queries, logger *slog.Logger) (*Server, error) {
	srv := &Server{
		cfg:            cfg,
		db:             db,
		logger:         logger,
		fileserverHits: atomic.Int32{},
	}

	return srv, nil
}
