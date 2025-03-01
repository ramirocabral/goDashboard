package main

import (
	"context"
	"log"
	"net/http"
	"time"

	// "github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	// "golang-system-monitor/internal/collector/cpu"
	"golang-system-monitor/internal/configuration"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/influxdb"
	"golang-system-monitor/internal/logger"
	"golang-system-monitor/internal/subscribers"

	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func main(){

    cfg :=  configuration.GetConfig()

    // ctx := context.Background()
    eb := core.NewEventBus()
    cpuTopic := eb.CreateTopic("cpu")

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

    dbSubscriber := subscribers.NewStorageSubscriber(db)
    go dbSubscriber.Subscribe(cpuTopic)
    
    // //collectors
    // cpuCollector := cpu.NewCpuCollector(2*time.Second, eb)
    // go cpuCollector.Start(ctx)
    //
    //
    // http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
    //     conn, err := upgrader.Upgrade(w, r, nil)
    //     if err != nil{
    //         log.Fatal("Error upgrading connection: ", err)
    //         return
    //     }
    //
    //     ws := subscribers.NewWebSocketSubscriber(conn)
    //
    //     go ws.Subscribe(cpuTopic)
    //
    //
    //     defer ws.Unsubscribe(cpuTopic)
    // })
    //
    log.Fatal(http.ListenAndServe(":8080", nil))
}
// create event bus
// create topics and assign it to the event bus
// create collectors
