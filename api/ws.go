package api

import (
    "net/http"
    "github.com/gorilla/websocket"
    "golang-system-monitor/internal/logger"
    "golang-system-monitor/internal/subscribers"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func (a *app) wsCPUHandler(w http.ResponseWriter, r *http.Request){
    conn, err := upgrader.Upgrade(w, r, nil)

    if err != nil{
	logger.GetLogger().Errorf("Error upgrading connection: %s", err)
	return 
    }

    ws := subscribers.NewWebSocketSubscriber(conn)

    go ws.Subscribe(a.eb.Topics["cpu"])
}
