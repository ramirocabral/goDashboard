package api

//logging middleware

import (
	"golang-system-monitor/internal/logger"
	"net/http"
	"time"
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

func (app *app) ValidateTimeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()

        startStr := query.Get("start")
        endStr := query.Get("end")
        intervalStr := query.Get("interval")

        if startStr == "" || endStr == "" || intervalStr == "" {
            app.badRequest(w, r, "Missing query parameters")
            return
        }

        start, err := validateRFC3339(startStr)
        if err != nil {
            app.badRequest(w, r, "Invalid start time format. Use RFC3339")
            return
        }

        end, err := validateRFC3339(endStr)
        if err != nil {
            app.badRequest(w, r, "Invalid end time format. Use RFC3339")
            return
        }

        // check that end time is after start time
        if end.Before(start) {
            app.badRequest(w, r, "End time must be after start time")
            return
        }

        _, err = validateInterval(intervalStr)
        if err != nil {
            app.badRequest(w, r, "Invalid interval format. Use durations like 1s, 5m, 1h")
            return
        }

        next.ServeHTTP(w, r)
    })
}

func validateRFC3339(param string) (time.Time, error) {
    return time.Parse(time.RFC3339, param)
}

func validateInterval(interval string) (time.Duration, error) {
    return time.ParseDuration(interval)
}
