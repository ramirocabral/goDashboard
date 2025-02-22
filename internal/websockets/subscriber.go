package websockets

import (
    "log"
    "sync"

    "github.com/gorilla/websocket"
    "golang-system-monitor/internal/core"
)

type WebSocketSubscriber struct{
    Id          string              //usually the ip addr of the client
    Conn        *websocket.Conn     //websocket connection
    EventBus    *core.EventBus           //EventBus assigned
    Topics      map[string]struct{} //topics subscribed to
    Mu          sync.RWMutex        //websocket mutex for safety
}

func (ws *WebSocketSubscriber) ID() string{
    return ws.Id
}

func (ws *WebSocketSubscriber) Handle(msg core.Message){
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    if _, ok := ws.Topics[msg.Type]; !ok{
        return
    }

    err := ws.Conn.WriteJSON(msg)
    if err != nil{
        log.Println("Error writing to websocket: ", err)
        ws.HandleDisconnect()
    }
}

//subscribe to a topic
func (ws *WebSocketSubscriber) Subscribe(topic string) error{
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    if _, ok := ws.Topics[topic]; ok{
        return nil
    }

    ws.Topics[topic] = struct{}{}

    ws.EventBus.Topics[topic].Mu.Lock()
    ws.EventBus.Topics[topic].Subscribers[ws.ID()] = ws
    ws.EventBus.Topics[topic].Mu.Unlock()

    return nil
}

//unsubscribe from a topic
func (ws *WebSocketSubscriber) Unsubscribe(topic string) error{
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    if _, ok := ws.Topics[topic]; !ok{
        return nil
    }

    //remove the topics from the ws list
    delete(ws.Topics, topic)

    //remove the ws from the topic's subscribers
    ws.EventBus.Topics[topic].Mu.Lock()
    delete(ws.EventBus.Topics[topic].Subscribers, ws.ID())
    ws.EventBus.Topics[topic].Mu.Unlock()

    return nil
}

func (ws *WebSocketSubscriber) HandleDisconnect(){
    ws.EventBus.RemoveSubscriber(ws.Id)  
    ws.Close()
}

func (ws *WebSocketSubscriber) Close(){
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    for topic := range ws.Topics{
        ws.Unsubscribe(topic)
    }

    ws.Conn.Close()
}

func NewWebSocketSubscriber(conn *websocket.Conn, eb *core.EventBus) *WebSocketSubscriber{
    return &WebSocketSubscriber{
        Id: conn.RemoteAddr().String(),
        Conn: conn,
        EventBus: eb,
        Topics: make(map[string]struct{}),
    }
}
