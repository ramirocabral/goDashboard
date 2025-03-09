package collector

import(
    "go-dashboard/pkg/stats"
    "go-dashboard/internal/core"
    "go-dashboard/internal/logger"
)

type UptimeCollector struct{
    statsCollector     *stats.StatsManager
}

func NewUptimeCollector(statsCollector *stats.StatsManager) *UptimeCollector{
    return &UptimeCollector{
        statsCollector: statsCollector,
    }
}

func (m *UptimeCollector) Collect() (core.Storable, error){
    uptimeData, err := m.statsCollector.GetUptime()
    if err != nil{
        logger.GetLogger().Error("Error reading uptime data: ", err)
        return nil, err
    }
    return uptimeData, nil
}

func (c *UptimeCollector) GetTopic() string{
    return "uptime"
}
