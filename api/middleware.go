package api
//logging middleware

import (
    "net/http"
)

func (app *app) LoggingMiddleware(next http.Handler) http.Handler {
    //type  conversion from a function to a handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.logger.Infow(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
        next.ServeHTTP(w, r)
    })
}
