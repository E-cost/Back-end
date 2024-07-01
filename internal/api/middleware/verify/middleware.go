package verify

import (
	"github.com/google/uuid"
	"net/http"
	"regexp"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emailID := r.URL.Query().Get("email_id")
		secretCode := r.URL.Query().Get("secret_code")

		if _, err := uuid.Parse(emailID); err != nil {
			http.Error(w, "Invalid email_id", http.StatusBadRequest)
			return
		}

		match, err := regexp.MatchString("^[a-zA-Z]{32}$", secretCode)
		if err != nil {
			http.Error(w, "Invalid secret_code", http.StatusBadRequest)
			return
		}

		if !match {
			http.Error(w, "Invalid secret_code", http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
