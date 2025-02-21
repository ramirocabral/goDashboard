package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang-system-monitor/internal/collector"

	"github.com/gorilla/websocket"
        "github.com/google/uuid"
)

type Message struct{
    Type            string
    Timestamp       time.Time
    Data            interface{}
}

type Publisher interface{
    Publish(Message)
    Start(context.Context)  error
    Stop()      error
}

type Subscriber interface{
    ID()                    string
    Handle(Message)
    Subscribe(...string)    error
    Unsubscribe(...string)  error
}

type MetricCollector struct{
    refreshRate time.Duration
    metricType  string
    eventBus    *EventBus
    stopChan    chan struct{}
}

type EventBus struct{
    topics  map[string]*Topic
    mu      sync.RWMutex
}

func (eb *EventBus) CreateTopic(name string) *Topic{
    eb.mu.Lock()
    defer eb.mu.Unlock()

    topic := &Topic{
        name: name,
        subscribers: make(map[string]Subscriber),
        messages: make(chan Message, 1000),
    }

    eb.topics[name] = topic

    go topic.dispatch()

    return topic
}

type Topic struct{
    name        string
    subscribers map[string]Subscriber
    mu          sync.RWMutex
    messages    chan Message
}

func (t *Topic) AddSubscriber(sub Subscriber){
    t.mu.Lock()
    defer t.mu.Unlock()

    t.subscribers[sub.ID()] = sub
}


type WebSocketSubscriber struct{
    id          string
    conn        *websocket.Conn
    eventBus    *EventBus
    topics      map[string]struct{}
    mu          sync.RWMutex        //websocket mutex for safety
}

type DatabaseSubscriber struct{
    db      *sql.DB
    eventBus    *EventBus
}

func NewEventBus() *EventBus{
    return &EventBus{
        topics: make(map[string]*Topic),
    }
}

func (eb *EventBus) AddTopic(name string){
    eb.mu.Lock()
    defer eb.mu.Unlock()

    if _, ok := eb.topics[name]; !ok{
        eb.topics[name] = &Topic{
            name: name,
            subscribers: make(map[string]Subscriber),
            messages: make(chan Message),
        }
    }
}

func (t *Topic) dispatch(){
    for msg := range t.messages{
        t.mu.RLock()
        for _, sub := range t.subscribers{
            go sub.Handle(msg)
        }
        t.mu.RUnlock()
    }
}

type CPUCollector struct{
    MetricCollector
}

//publisher
func (c *CPUCollector) Start(ctx context.Context) error{
    ticker := time.NewTicker(c.refreshRate)
    defer ticker.Stop()

    for {
        select{
            case <-ctx.Done():
                return ctx.Err()
            case <-ticker.C:
                // get cpu data
                cpuData, err := collector.ReadCPU()
                if err != nil{
                    log.Fatal("Error reading cpu data: ", err)
                    continue
                }

                c.eventBus.topics["cpu"].messages <- Message{
                    Type: "cpu",
                    Timestamp: time.Now(),
                    Data: cpuData,
                }
        }
    }
}

func NewCpuCollector(refreshRate time.Duration, eb *EventBus) *CPUCollector{
    return &CPUCollector{
        MetricCollector: MetricCollector{
            refreshRate: refreshRate,
            metricType: "cpu",
            eventBus: eb,
            stopChan: make(chan struct{}),
        },
    }
}

func (ws *WebSocketSubscriber) ID() string{
    return ws.id
}

func (ws *WebSocketSubscriber) Handle(msg Message){
    ws.mu.Lock()
    defer ws.mu.Unlock()

    if _, ok := ws.topics[msg.Type]; !ok{
        return
    }

    //enviar al websocket
    err := ws.conn.WriteJSON(msg)
    if err != nil{
        log.Fatal("Error writing to websocket: ", err)
    }
}

func (ws *WebSocketSubscriber) Subscribe(topics ...string) error{
    ws.mu.Lock()
    defer ws.mu.Unlock()

    for _, topic := range topics{
        if _, ok := ws.topics[topic]; ok{
            continue
        }

        ws.topics[topic] = struct{}{}

        ws.eventBus.topics[topic].mu.Lock()
        ws.eventBus.topics[topic].subscribers[ws.conn.RemoteAddr().String()] = ws
        ws.eventBus.topics[topic].mu.Unlock()
    }

    return nil
}

func (ws *WebSocketSubscriber) Unsubscribe(topics ...string) error{
    ws.mu.Lock()
    defer ws.mu.Unlock()

    for _, topic := range topics{
        if _, ok := ws.topics[topic]; !ok{
            continue
        }

        delete(ws.topics, topic)

        ws.eventBus.topics[topic].mu.Lock()
        delete(ws.eventBus.topics[topic].subscribers, ws.conn.RemoteAddr().String())
        ws.eventBus.topics[topic].mu.Unlock()
    }

    return nil
}

func NewWebSocketSubscriber(conn *websocket.Conn, eb *EventBus) *WebSocketSubscriber{
    return &WebSocketSubscriber{
        id: uuid.New().String(),
        conn: conn,
        eventBus: eb,
        topics: make(map[string]struct{}),
    }
}

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}

func main(){

    ctx := context.Background()
    eb := NewEventBus()

    cpuTopic := eb.CreateTopic("cpu")

    fmt.Println(cpuTopic.name)

    cpuCollector := NewCpuCollector(2*time.Second, eb)

    go cpuCollector.Start(ctx)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil{
            log.Fatal("Error upgrading connection: ", err)
            return
        }

        ws := NewWebSocketSubscriber(conn, eb)

        ws.topics["cpu"] = struct{}{}

        go ws.Subscribe("cpu")

        defer ws.Unsubscribe("cpu")
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
