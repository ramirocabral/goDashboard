package main

import (
    "context"
    "time"

    "golang-system-monitor/internal/collector"
    "golang-system-monitor/internal/configuration"
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/influxdb"
    "golang-system-monitor/internal/logger"
    "golang-system-monitor/internal/subscribers"
    "golang-system-monitor/pkg/stats"
    "golang-system-monitor/api"
)

func main(){
    cfg :=  configuration.GetConfig()
    logger.Init("prod")

    //database
    db, err := influxdb.New(
	"http://influxdb2:8086",
	cfg.DB.Token,
	cfg.DB.Org,
	cfg.DB.Bucket,
    )
    if err != nil{
	logger.GetLogger().Fatal("Error connecting to database: ", err)
    }
    logger.GetLogger().Infof("Connected to database: %s", cfg.DB.Addr)
    defer db.Close()


    eb := core.NewEventBus()

    cpuTopic := eb.CreateTopic("cpu")
    memTopic := eb.CreateTopic("memory")
    ioTopic := eb.CreateTopic("io")
    containerTopic := eb.CreateTopic("container")
    networkTopic := eb.CreateTopic("network")
    uptimeTopic := eb.CreateTopic("uptime")
    _ = containerTopic
    _ = uptimeTopic

    dbSubscriber := subscribers.NewStorageSubscriber(db)
    go dbSubscriber.Subscribe(cpuTopic)
    go dbSubscriber.Subscribe(memTopic)
    go dbSubscriber.Subscribe(ioTopic)
    go dbSubscriber.Subscribe(networkTopic)

    ctx := context.Background()
    statsManager := stats.NewStatsManager()

    initCollectors(eb, statsManager, ctx)


    app := api.NewApp(cfg, db, eb, statsManager)
    
    mux := app.Mount()

    err = app.Run(mux)

    if err != nil{
	logger.GetLogger().Fatal("Error running server: ", err)
    }
}

func initCollectors(eb *core.EventBus, statsManager *stats.StatsManager, ctx context.Context){

    logger.GetLogger().Info("Initializing collectors")
    cpuMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*1),
	eb,
	collector.NewCPUCollector(statsManager),
    )
    memMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*2),
	eb,
	collector.NewMemoryCollector(statsManager),
    )
    ioMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*1),
	eb,
	collector.NewIOCollector(statsManager),
    )
    containerMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*5),
	eb,
	collector.NewContainerCollector(statsManager),
    )
    networkMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*1),
	eb,
	collector.NewNetworkCollector(statsManager),
    )
    uptimeMetricCollector := core.NewMetricCollector(
	time.Duration(time.Second*2),
	eb,
	collector.NewUptimeCollector(statsManager),
    )
    
    go cpuMetricCollector.Start(ctx)
    go memMetricCollector.Start(ctx)
    go ioMetricCollector.Start(ctx)
    go containerMetricCollector.Start(ctx)
    go networkMetricCollector.Start(ctx)
    go uptimeMetricCollector.Start(ctx)
}
