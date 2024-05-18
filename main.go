package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")

	type apiConfig struct {
		fileserverHits int	
	}

	func (cfg *apiConfig) middleWareMetricsInc(next http.Handler) http.Handler {
		cfg.fileserverHits += 1;
		return next
	}
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Handle root path
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is not the root
		if r.URL.Path != "/healthz" {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		// Provide a response for the root path
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	// Create a new server with the ServeMux as the handler
	server := &http.Server{
		Addr:    ":8080", // Set the server to listen on port 8080
		Handler: mux,
	}

	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
