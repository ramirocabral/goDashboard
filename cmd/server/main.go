package main

import (
	"context"
	"database/sql"
	"fmt"

	// "fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"golang-system-monitor/internal/collector/cpu"
	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/websockets"
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

    cpuCollector := cpu.NewCpuCollector(2*time.Second, eb)

    go cpuCollector.Start(ctx)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil{
            log.Fatal("Error upgrading connection: ", err)
            return
        }

        ws := websockets.NewWebSocketSubscriber(conn, eb)

        ws.Topics["cpu"] = struct{}{}

        go ws.Subscribe("cpu")

        defer ws.Unsubscribe("cpu")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
