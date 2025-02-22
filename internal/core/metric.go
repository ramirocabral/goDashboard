package core

import (
    "time"
)

type MetricCollector struct{
    RefreshRate time.Duration
    MetricType  string
    EventBus    *EventBus
    StopChan    chan struct{}
}

