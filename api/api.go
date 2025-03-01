package api

import (
	// "golang-system-monitor/internal/collector"
	"golang-system-monitor/internal/configuration"
	"golang-system-monitor/internal/core"
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


