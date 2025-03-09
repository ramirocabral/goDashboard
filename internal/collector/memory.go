package collector

import(
    "go-dashboard/pkg/stats"
    "go-dashboard/internal/core"
    "go-dashboard/internal/logger"
)

type MemoryCollector struct{
    statsCollector     *stats.StatsManager
}

func NewMemoryCollector(statsCollector *stats.StatsManager) *MemoryCollector{
    return &MemoryCollector{
        statsCollector: statsCollector,
    }
}

func (m *MemoryCollector) Collect() (core.Storable, error){
    memData, err := m.statsCollector.GetMemory()
    if err != nil{
        logger.GetLogger().Error("Error reading memory data: ", err)
        return nil, err
    }
    return memData, nil
}

func (c *MemoryCollector) GetTopic() string{
    return "memory"
}
