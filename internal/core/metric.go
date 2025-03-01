package core

import (
    "time"
    "context"
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

func (mc *MetricCollector) Start(ctx context.Context) error {
    ticker := time.NewTicker(mc.RefreshRate)

    defer ticker.Stop()

    for {
        select {
            case <-ctx.Done():
                return ctx.Err()
            case <-ticker.C:
                data, err := mc.Collector.Collect()
                if err != nil {
                    return err
                }
                mc.EventBus.Publish(mc.Collector.GetTopic(), &Message{
                    Type: mc.Collector.GetTopic(),
                    Data: data,
                })
        }
    }
}


