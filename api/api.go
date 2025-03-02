package api

import (
	"golang-system-monitor/internal/configuration"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/logger"
	"golang-system-monitor/internal/storage"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type app struct{
    cfg	    configuration.Config
    store   storage.Storage
    eb      *core.EventBus
}

func NewApp(cfg configuration.Config,store storage.Storage, eb *core.EventBus) *app{
    return &app{
	cfg: cfg,
	store: store,
	eb: eb,
    }
}

// create the handlers for the mux
func (app *app) Mount() http.Handler{
    r := mux.NewRouter()

    //set custom middlewares
    r.Use(app.LoggingMiddleware)
    r.Use(app.CORSMiddleware)

    r.HandleFunc("/ws/cpu", app.wsCPUHandler)

    //add the routes
    return r
}

//start the server
func (app *app) Run(mux http.Handler) error {
    srv := &http.Server{
	    Addr:         app.cfg.APIPort,
	    Handler:      mux,
	    WriteTimeout: time.Second * 30,
	    ReadTimeout:  time.Second * 10,
	    IdleTimeout:  time.Minute,
    }

    logger.GetLogger().Infof("Starting server on port %s, env %s", app.cfg.APIPort, app.cfg.Env)

    err := srv.ListenAndServe()

    return err
}
