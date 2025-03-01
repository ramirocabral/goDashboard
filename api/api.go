package api

import (
	"golang-system-monitor/internal/collector/cpu"
	"golang-system-monitor/internal/configuration"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/storage"

	"go.uber.org/zap"
)

type app struct{
    cfg	    configuration.Config
    logger  *zap.SugaredLogger
    store   storage.Storage
    eb      *core.EventBus
}

func NewApp(cfg configuration.Config, logger *zap.SugaredLogger, store storage.Storage, eb *core.EventBus, cpuCollector *cpu.CPUCollector) *app{
    return &app{
    cfg: cfg,
    logger: logger,
    store: store,
    eb: eb,
    }
}
