package collector

import(
    "go-dashboard/pkg/stats"
    "go-dashboard/internal/core"
    "go-dashboard/internal/logger"
)

type ContainerCollector struct{
    statsCollector     *stats.StatsManager
}

func NewContainerCollector(statsCollector *stats.StatsManager) *ContainerCollector{
    return &ContainerCollector{
        statsCollector: statsCollector,
    }
}

func (c *ContainerCollector) Collect() (core.Storable, error){
    containerData, err := c.statsCollector.GetContainers()
    if err != nil{
        logger.GetLogger().Error("Error reading container data: ", err)
        return nil, err
    }
    return containerData, nil
}

func (c *ContainerCollector) GetTopic() string{
    return "container"
}
