package dns

import (
	"Ecost/internal/redis/service"
	remoteDNS "Ecost/internal/utils/ip"
	"Ecost/pkg/logging"
	"net/http"
)

type Middleware struct {
	logger       *logging.Logger
	redisService service.RedisIpRequestService
}

func NewIpRequestMiddleware(logger *logging.Logger, redisService service.RedisIpRequestService) *Middleware {
	return &Middleware{
		logger:       logger,
		redisService: redisService,
	}
}

func (m *Middleware) IpMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IpOutput, err := remoteDNS.ReadUserIP(r)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		count, err := m.redisService.GetRequestCount(r.Context(), IpOutput.RealIP, r.URL.Path)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count >= 10 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		err = m.redisService.IncrementRequestCount(r.Context(), IpOutput.RealIP, r.URL.Path)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}
}
