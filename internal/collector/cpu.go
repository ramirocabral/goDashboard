package collector

import(
    "golang-system-monitor/pkg/stats"
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/logger"
)

type CPUCollector struct{
    statsCollector     *stats.StatsManager
}

func NewCPUCollector(statsCollector *stats.StatsManager) *CPUCollector{
    return &CPUCollector{
        statsCollector: statsCollector,
    }
}

func (c *CPUCollector) Collect() (core.Storable, error){
    cpuData, err := c.statsCollector.GetCPU()
    if err != nil{
        logger.GetLogger().Error("Error reading cpu data: ", err)
        return nil, err
    }
    return cpuData, nil
}

func (c *CPUCollector) GetTopic() string{
    return "cpu"
}
