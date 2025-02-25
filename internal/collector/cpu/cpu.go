package cpu

import(
    "log"
    "context"
    "time"
    "golang-system-monitor/internal/stats/cpu"
    "golang-system-monitor/internal/core"
)

type CPUCollector struct{
    core.MetricCollector
}

//publisher
func (c *CPUCollector) Start(ctx context.Context) error{
    ticker := time.NewTicker(c.RefreshRate)
    defer ticker.Stop()

    for {
        select{
            case <-ctx.Done():
                return ctx.Err()
            case <-ticker.C:
                // get cpu data
                cpuData, err := cpu.ReadCPU()
                if err != nil{
                    log.Fatal("Error reading cpu data: ", err)
                    continue
                }
                //create messae struct and send it to the event bus
                c.EventBus.Topics["cpu"].Messages <- &core.Message{
                    Type: "cpu",
                    Timestamp: time.Now(),
                    Data: cpuData,
                }
        }
    }
}

func NewCpuCollector(refreshRate time.Duration, eb *core.EventBus) *CPUCollector{
    return &CPUCollector{
        MetricCollector: core.MetricCollector{
            RefreshRate: refreshRate,
            MetricType: "cpu",
            EventBus: eb,
            StopChan: make(chan struct{}),
        },
    }
}
