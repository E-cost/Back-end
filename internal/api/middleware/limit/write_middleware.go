package limit

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

func WriteMiddleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(1, 5)

	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Write request failed",
				Body:   "Limit Exceeded",
			}

			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(&message); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		} else {
			next.ServeHTTP(w, r)
		}
	}
}
