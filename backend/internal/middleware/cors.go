package middleware

import (
	"net/http"
	
	"github.com/rs/cors"
)

func EnableCORS(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: false,
	})
	
	return c.Handler(handler)
}