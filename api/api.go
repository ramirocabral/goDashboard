package api

import (
	"net/http"
	"time"
	"github.com/gorilla/mux"

	"golang-system-monitor/internal/configuration"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/logger"
	"golang-system-monitor/internal/storage"
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
    r := mux.NewRouter().PathPrefix("/v1").Subrouter()

    //set custom middlewares
    r.Use(app.LoggingMiddleware)
    r.Use(app.CORSMiddleware)

    ws := r.PathPrefix("/ws").Subrouter()
    // stats := r.PathPrefix("/stats").Subrouter()
    // health := r.PathPrefix("/health").Subrouter()
    // 
    ws.HandleFunc("/cpu", app.wsCPUHandler)
    ws.HandleFunc("/memory", app.wsMemoryHandler)
    ws.HandleFunc("/io", app.wsIOHandler)
    ws.HandleFunc("/container", app.wsContainerHandler)
    ws.HandleFunc("/network", app.wsNetworkHandler)
    ws.HandleFunc("/uptime", app.wsUptimeHandler)

    // stats.HandleFunc("/cpu", app.statsCPUHandler)
    // stats.HandleFunc("/memory", app.statsMemoryHandler)
    // stats.HandleFunc("/io", app.statsIOHandler)

    // health.HandleFunc("/cpu", app.healthCPUHandler)

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
