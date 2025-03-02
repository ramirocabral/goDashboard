package main

import (
    "context"
    "time"

    // "github.com/gorilla/mux"
    "github.com/gorilla/websocket"

    "golang-system-monitor/internal/collector"
    "golang-system-monitor/internal/configuration"
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/influxdb"
    "golang-system-monitor/internal/logger"
    "golang-system-monitor/internal/subscribers"
    "golang-system-monitor/pkg/stats"
    "golang-system-monitor/api"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func main(){
    cfg :=  configuration.GetConfig()

    //database
    db, err := influxdb.New(
	cfg.DB.Addr,
	cfg.DB.Token,
	cfg.DB.Org,
	cfg.DB.Bucket,
    )
    if err != nil{
	logger.GetLogger().Fatal("Error connecting to database: ", err)
    }
    logger.GetLogger().Infof("Connected to database: %s", cfg.DB.Addr)
    defer db.Close()

    statsManager := stats.NewStatsManager()
    ctx := context.Background()
    eb := core.NewEventBus()

    cpuTopic := eb.CreateTopic("cpu")
    memTopic := eb.CreateTopic("memory")
    // ioTopic := eb.CreateTopic("io")
    // containerTopic := eb.CreateTopic("container")
    // _ = containerTopic

    dbSubscriber := subscribers.NewStorageSubscriber(db)
    go dbSubscriber.Subscribe(cpuTopic)
    go dbSubscriber.Subscribe(memTopic)
    // go dbSubscriber.Subscribe(ioTopic)

    initCollectors(eb, statsManager, ctx)
    // log.Fatal(http.ListenAndServe(":8080", nil))
    app := api.NewApp(cfg, db, eb)
    
    mux := app.Mount()

    err = app.Run(mux)

    if err != nil{
	logger.GetLogger().Fatal("Error starting server: ", err)
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
	//
 //    ioMetricCollector := core.NewMetricCollector(
	// time.Duration(time.Second*1),
	// eb,
	// collector.NewIOCollector(statsManager),
 //    )
	//
 //    containerMetricCollector := core.NewMetricCollector(
	// time.Duration(time.Second*1),
	// eb,
	// collector.NewContainerCollector(statsManager),
 //    )

    go cpuMetricCollector.Start(ctx)
    go memMetricCollector.Start(ctx)
    // go ioMetricCollector.Start(ctx)
    // go containerMetricCollector.Start(ctx)
}
