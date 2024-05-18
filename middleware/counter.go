package middleware

import (
	"fmt"
	"net/http"
)

// ApiConfig holds configuration data, including the counter
type ApiConfig struct {
	FileserverHits int
}

// MiddlewareMetricsInc is a middleware that increments the request counter
func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

// HandlerMetrics serves the metrics page
func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	template := fmt.Sprintf(
		`<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>`, cfg.FileserverHits)
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template))
}

// HandlerMetricsReset resets the counter
func (cfg *ApiConfig) HandlerMetricsReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits = 0
	w.Write([]byte("Metrics reset to 0"))
}
