package core

import (
    "time"
)

// a MetricCollector collects metrics of the system every RefreshRate and sends it to the event bus
type MetricCollector struct{
    RefreshRate time.Duration
    MetricType  string
    EventBus    *EventBus
    StopChan    chan struct{}
}
