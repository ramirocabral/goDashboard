package collector

import(
    "go-dashboard/pkg/stats"
    "go-dashboard/internal/core"
    "go-dashboard/internal/logger"
)

type NetworkCollector struct{
    statsCollector     *stats.StatsManager
}

func NewNetworkCollector(statsCollector *stats.StatsManager) *NetworkCollector{
    return &NetworkCollector{
        statsCollector: statsCollector,
    }
}

func (m *NetworkCollector) Collect() (core.Storable, error){
    netData, err := m.statsCollector.GetNetwork()
    if err != nil{
        logger.GetLogger().Error("Error reading network data: ", err)
        return nil, err
    }
    return netData, nil
}

func (c *NetworkCollector) GetTopic() string{
    return "network"
}
