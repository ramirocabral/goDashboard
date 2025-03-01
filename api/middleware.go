package api
//logging middleware

import (
    "net/http"
    "log"
)

//gorilla mux logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Printf("Request URI: %s", r.RequestURI)
        next.ServeHTTP(w, r) // Call the next handler
    })
}

