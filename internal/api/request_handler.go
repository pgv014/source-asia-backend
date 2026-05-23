package api

import (
	"encoding/json"
	"net/http"

	"source-asia/internal/ratelimit"
)

type RequestBody struct {
	UserID string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

func RequestHandler(limiter *ratelimit.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var body RequestBody

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid JSON",
			})

			return
		}

		if body.UserID == "" {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "user_id is required",
			})

			return
		}

		if body.Payload == nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "payload is required",
			})

			return
		}

		allowed := limiter.Allow(body.UserID)

		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})

			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "accepted",
			"user_id": body.UserID,
		})
	}
} 

