package cors

import (
	"net/http"

	"github.com/rs/cors"
)

func CorsSettings() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://172.20.10.6:3000",
		},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Access-Control-Allow-Origin",
			"Content-Type",
			"Authorization",
			"X-Real-Ip",
			"X-Forwarded-For"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{},
		Debug:              false,
	})

	return c
}
