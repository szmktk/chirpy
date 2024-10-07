package main

import (
	"fmt"
	"net/http"
	"os"

	"log/slog"
)

const port = 8080

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	mux := http.NewServeMux()
	// commented out because the first exercise is about returning 404 response
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "OK")
	// })

	logger.Info("Starting the server", "port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
