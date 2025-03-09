package collector

import(
    "go-dashboard/pkg/stats"
    "go-dashboard/internal/core"
    "go-dashboard/internal/logger"
)

type IOCollector struct{
    statsCollector     *stats.StatsManager
}

func NewIOCollector(statsCollector *stats.StatsManager) *IOCollector{
    return &IOCollector{
        statsCollector: statsCollector,
    }
}

func (m *IOCollector) Collect() (core.Storable, error){
    ioData, err := m.statsCollector.GetIO()
    if err != nil{
        logger.GetLogger().Error("Error reading io data: ", err)
        return nil, err
    }
    return ioData, nil
}

func (c *IOCollector) GetTopic() string{
    return "io"
}
