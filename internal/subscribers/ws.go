package subscribers

import (
	"sync"

	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/logger"
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
        Mu: &sync.RWMutex{},
    }
}

func (ws *WebSocketSubscriber) ID() string{
    return ws.Id
}

// handle function, executes when a message is received
func (ws *WebSocketSubscriber) Handle(msg *core.Message){
    ws.Mu.Lock()

    if _, ok := ws.Topics[msg.Type]; !ok{
        ws.Mu.Unlock()
        return
    }

    err := ws.Conn.WriteJSON(msg)
    ws.Mu.Unlock()
    if err != nil{
        logger.GetLogger().Errorf("Error writing to websocket: %s", err)
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
    topics := ws.Topics
    ws.Mu.Unlock()

    logger.GetLogger().Debugf("Disconnecting websocket: %s", ws.Id)
    for _, topic := range topics{
        ws.Unsubscribe(topic)
    }

    ws.Conn.Close()
}
