package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s call on endpoint %s", r.Method, r.URL)
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
