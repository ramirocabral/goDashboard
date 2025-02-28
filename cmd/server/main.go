package main

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"golang-system-monitor/internal/collector/cpu"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/influxdb"
	"golang-system-monitor/internal/subscribers"
)

type DatabaseSubscriber struct{
    db      *sql.DB
    eventBus    *core.EventBus
}

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

type config struct{
    addr	string
    apiURL	string
    db      	dbConfig
    env		string
}

//influxdb config
type dbConfig struct{
    addr    string
    token   string
    org     string
    bucket  string
}

func main(){

    ctx := context.Background()
    eb := core.NewEventBus()
    cpuTopic := eb.CreateTopic("cpu")
    fmt.Println(cpuTopic.Name)
    db, err := influxdb.New("http://localhost:8086", "mytoken","my-org", "mybucket")
    //create bucket

    if err != nil{
	log.Fatal("Error creating storage: ", err)
    }

    dbSubscriber := subscribers.NewStorageSubscriber(db)

    go dbSubscriber.Subscribe(cpuTopic)
    
    cpuCollector := cpu.NewCpuCollector(2*time.Second, eb)

    go cpuCollector.Start(ctx)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil{
            log.Fatal("Error upgrading connection: ", err)
            return
        }

        ws := subscribers.NewWebSocketSubscriber(conn)

        go ws.Subscribe(cpuTopic)


        defer ws.Unsubscribe(cpuTopic)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
