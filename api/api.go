package api

import (
	"net/http"
	"time"
	"github.com/gorilla/mux"

	"go-dashboard/internal/configuration"
	"go-dashboard/internal/core"
	"go-dashboard/internal/logger"
	"go-dashboard/internal/storage"
	"go-dashboard/pkg/stats"
)

type app struct{
    cfg		configuration.Config
    store	storage.Storage
    eb		*core.EventBus
    statsManager *stats.StatsManager
}

func NewApp(cfg configuration.Config,store storage.Storage, eb *core.EventBus, statsManager *stats.StatsManager) *app{
    return &app{
	cfg: cfg,
	store: store,
	eb: eb,
	statsManager: statsManager,
    }
}

// create the routes
func (app *app) Mount() http.Handler {
    mux := mux.NewRouter()

    r := mux.PathPrefix("/api/v1").Subrouter()

    // Set custom middlewares
    r.Use(app.LoggingMiddleware)
    r.Use(app.CORSMiddleware)

    ws := r.PathPrefix("/ws").Subrouter()
    ws.HandleFunc("/cpu", app.wsCPUHandler)
    ws.HandleFunc("/memory", app.wsMemoryHandler)
    ws.HandleFunc("/io", app.wsIOHandler)
    ws.HandleFunc("/container", app.wsContainerHandler)
    ws.HandleFunc("/network", app.wsNetworkHandler)
    ws.HandleFunc("/uptime", app.wsUptimeHandler)

    stats := r.PathPrefix("/stat").Subrouter()
    stats.HandleFunc("/smart", app.smartHandler)
    stats.HandleFunc("/host", app.hostHandler)
    stats.HandleFunc("/disk", app.diskHandler)

    history := r.PathPrefix("/history").Subrouter()
    history.Use(app.ValidateTimeMiddleware)
    history.HandleFunc("/cpu", app.cpuHistoryHandler)
    history.HandleFunc("/memory", app.memoryHistoryHandler)
    history.HandleFunc("/io", app.ioHistoryHandler)
    history.HandleFunc("/network", app.networkHistoryHandler)

    // serve static files under "/assets/"
    mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("/dist/assets"))))

    // serve index.html for frontend routes with react
    mux.HandleFunc("/{_:.*}", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "/dist/index.html")
    })

    return mux
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

    return srv.ListenAndServe()
}
