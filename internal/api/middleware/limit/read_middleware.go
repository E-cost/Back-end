package limit

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"net/http"
)

func ReadMiddleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := rate.NewLimiter(1, 600)

	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "Read Request failed",
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
