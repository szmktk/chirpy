package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (srv *Server) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (srv *Server) HandlerMetrics(w http.ResponseWriter, _ *http.Request) {
	template := `<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>`
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(template, srv.fileserverHits.Load())))
}

func (srv *Server) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if srv.cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Operation not allowed on platform: '%s'", srv.cfg.Platform))
		return
	}

	srv.fileserverHits.Store(0)
	if err := srv.db.DeleteUsers(r.Context()); err != nil {
		srv.logger.Error("Error deleting users: %s", "err", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hits reset to 0 and database reset to initial state"})
}
