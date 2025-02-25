package websockets

import (
	"log"
	"sync"

	"golang-system-monitor/internal/core"

	"github.com/gorilla/websocket"
)

type WebSocketSubscriber struct{
    Id          string              //usually the ip addr of the client
    Conn        *websocket.Conn     //websocket connection
    Topics      map[string]*core.Topic
    Mu          *sync.RWMutex        //mutex for safety, since multiple goroutines can access the same ws
}

func NewWebSocketSubscriber(conn *websocket.Conn) *WebSocketSubscriber{
    return &WebSocketSubscriber{
        Id: conn.RemoteAddr().String(),
        Conn: conn,
        Topics: make(map[string]*core.Topic),
    }
}

func (ws *WebSocketSubscriber) ID() string{
    return ws.Id
}

// handle function, executes when a message is received
func (ws *WebSocketSubscriber) Handle(msg *core.Message){
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
func (ws *WebSocketSubscriber) Subscribe(topic *core.Topic) error{
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    //if already subscribed, return
    if _, ok := ws.Topics[topic.Name]; ok{
        return nil
    }

    ws.Topics[topic.Name] = topic

    //add the ws to the topic's subscribers
    topic.AddSubscriber(ws)

    return nil
}

//unsubscribe from a topic
func (ws *WebSocketSubscriber) Unsubscribe(topic *core.Topic) error{
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    //if not subscribed, return
    if _, ok := ws.Topics[topic.Name]; !ok{
        return nil
    }

    //remove the topics from the ws list
    delete(ws.Topics, topic.Name)

    //remove the ws from the topic's subscribers
    topic.RemoveSubscriber(ws)

    return nil
}

// close the connection and unsubscribe from all topics
func (ws *WebSocketSubscriber) HandleDisconnect(){
    ws.Mu.Lock()
    defer ws.Mu.Unlock()

    for _, topic := range ws.Topics{
        ws.Unsubscribe(topic)
    }

    ws.Conn.Close()
}
