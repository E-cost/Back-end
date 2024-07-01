package errors

import (
	"Ecost/internal/apperror"
	"errors"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func Middleware(next appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var appErr *apperror.AppError
		err := next(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, apperror.ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(apperror.ErrNotFound.Marshal())
					return
				}

				if errors.Is(err, apperror.ForbiddenErr) {
					w.WriteHeader(http.StatusForbidden)
					w.Write(apperror.ForbiddenErr.Marshal())
					return
				}

				if errors.Is(err, apperror.TooManyRequestsErr) {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write(apperror.TooManyRequestsErr.Marshal())
					return
				}

				if errors.Is(err, apperror.ConflictErr) {
					w.WriteHeader(http.StatusConflict)
					w.Write(apperror.ErrNotFound.Marshal())
					return
				}

				if errors.Is(err, apperror.BadGatewayErr) {
					w.WriteHeader(http.StatusBadGateway)
					w.Write(apperror.ErrNotFound.Marshal())
					return
				}

				err = err.(*apperror.AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appErr.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(apperror.SystemError(err).Marshal())
		}
	}
}
