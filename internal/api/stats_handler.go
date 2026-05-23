package api

import (
	"encoding/json"
	"net/http"

	"source-asia/internal/ratelimit"
)

func StatsHandler(limiter *ratelimit.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats := limiter.GetStats()

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(stats)
	}
} 
