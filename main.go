package main

import (
	"fmt"
	"log"
	"net/http"
	"nhuxoll/bdd_chirpy/handler"
	"nhuxoll/bdd_chirpy/middleware"
)

func main() {
	fmt.Println("Starting server...")
	const filepathRoot = "."
	const port = "8080"

	apiCfg := middleware.ApiConfig{
		FileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handler.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandlerMetricsReset)
	mux.HandleFunc("POST /api/validate_chirp", handler.HandlerChirpsValidate)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
