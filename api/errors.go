package api

import (
    "net/http"

    "go-dashboard/internal/logger"
)

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error){
    logger.GetLogger().Errorf("Internal error, method: %s, path: %s, error: %s", r.Method, r.URL.Path, err)

    writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *app) badRequest(w http.ResponseWriter, r *http.Request, message string){
	logger.GetLogger().Errorf("Bad request, method: %s, path: %s, message: %s", r.Method, r.URL.Path, message)

	writeJSONError(w, http.StatusBadRequest, message)
}
