package collector

import(
    "golang-system-monitor/pkg/stats"
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/logger"
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
