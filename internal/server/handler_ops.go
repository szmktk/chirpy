package server

import (
	"fmt"
	"net/http"
)

func (srv *Server) Health(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
	return nil
}

func (srv *Server) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (srv *Server) Metrics(w http.ResponseWriter, _ *http.Request) error {
	template := `<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>`
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(template, srv.fileserverHits.Load())))
	return nil
}

func (srv *Server) Reset(w http.ResponseWriter, r *http.Request) error {
	if srv.cfg.Platform != "dev" {
		return APIError{Status: http.StatusForbidden, Msg: fmt.Sprintf("Operation not allowed on platform: '%s'", srv.cfg.Platform)}
	}

	srv.fileserverHits.Store(0)
	if err := srv.db.DeleteUsers(r.Context()); err != nil {
		srv.logger.Error("Error deleting users", "err", err)
		return APIError{Status: http.StatusInternalServerError, Msg: "Internal Server Error"}
	}
	return respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hits reset to 0 and database reset to initial state"})
}
