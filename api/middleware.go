package api
//logging middleware

import (
    "net/http"
    "golang-system-monitor/internal/logger"
)

func (app *app) LoggingMiddleware(next http.Handler) http.Handler {
    //type  conversion from a function to a handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.GetLogger().Infof("Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
        next.ServeHTTP(w, r)
    })
}

func (app *app) CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        next.ServeHTTP(w, r)
    })
}
