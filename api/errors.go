package api

import (
    "net/http"

    "golang-system-monitor/internal/logger"
)

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error){
    logger.GetLogger().Errorf("Internal error, method: %s, path: %s, error: %s", r.Method, r.URL.Path, err)

    writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}
