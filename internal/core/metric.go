package core

import (
    "time"
    "context"
    "golang-system-monitor/internal/logger"
)

type Collector interface{
    Collect()   (Storable, error)
    GetTopic()  string
}

// a MetricCollector collects metrics of the system every RefreshRate and sends it to the event bus
type MetricCollector struct{
    RefreshRate time.Duration
    EventBus    *EventBus
    Collector   Collector
}

func NewMetricCollector(refreshRate time.Duration, eventBus *EventBus, collector Collector) *MetricCollector{
    return &MetricCollector{
        RefreshRate: refreshRate,
        EventBus: eventBus,
        Collector: collector,
    }
}

func (mc *MetricCollector) Start(ctx context.Context) error {
    ticker := time.NewTicker(mc.RefreshRate)

    defer ticker.Stop()

    logger.GetLogger().Info("Starting metric collector for: ", mc.Collector.GetTopic())
    for {
        select {
            case <-ctx.Done():
                return ctx.Err()
            case <-ticker.C:
                data, err := mc.Collector.Collect()
                if err != nil {
                    logger.GetLogger().Error("Error collecting data: ", err)
                    continue
                }
                mc.EventBus.Publish(mc.Collector.GetTopic(), &Message{
                    Type: mc.Collector.GetTopic(),
                    Data: data,
                })
        }
    }
}
