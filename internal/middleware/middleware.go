package middleware

import (
	"log"
	"net/http"
)

func SetLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("IP: %-22s Endpoint: %s\n", r.RemoteAddr, r.URL)
		handler(w, r)
	}
}

func SetHeader(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
}

func HandleError(w http.ResponseWriter, err error, statusCode int, message string) {
	log.Printf("%s\n", err)
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
