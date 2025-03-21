package core

import (
    "sync"
    "time"
    "go-dashboard/internal/logger"
)

type Subscriber interface{
    ID()                    string
    Handle(*Message)
    Subscribe(*Topic)       error
    Unsubscribe(*Topic)     error
}

type Point struct{
    Timestamp       time.Time
    Measurement     string
    Tags            map[string]string
    Fields          map[string]interface{}
}

// not sure if it is the best implementation, but it solves my several dependency issues
type Storable interface{
    ToPoint() []*Point
}

type Message struct{
    Type            string      `json:"type"`
    Data            Storable    `json:"data"`
}

// a topic representes a type of message that can be published to the event bus, and eventually dispatched to all of its subscribers
type Topic struct{
    Name        string
    Subscribers map[string]Subscriber
    Mu          sync.RWMutex               
    Messages    chan *Message
}

// the eventbus is the main component of the system, it holds all the topics and dispatches the messages to the subscribers
type EventBus struct{
    Topics  map[string]*Topic
    Mu      sync.RWMutex
}

func (t *Topic) AddSubscriber(sub Subscriber){
    t.Mu.Lock()
    defer t.Mu.Unlock()
    t.Subscribers[sub.ID()] = sub
}

func (t *Topic) RemoveSubscriber(sub Subscriber){
    t.Mu.Lock()
    defer t.Mu.Unlock()

    delete(t.Subscribers, sub.ID())
}

func NewEventBus() *EventBus{
    return &EventBus{
        Topics: make(map[string]*Topic),
    }
}

// dispatches every message received to all of its Subscribers
func (t *Topic) dispatch(){
    for msg := range t.Messages{
        logger.GetLogger().Debugf("Dispatching message to topic: %s", t.Name)
        t.Mu.RLock()
        for _, sub := range t.Subscribers{
            logger.GetLogger().Debugf("Dispatching message to subscriber: %s", sub.ID())
            go sub.Handle(msg)
        }
        t.Mu.RUnlock()
    }
}

// add a new topic to the event bus, so the subscribers can subscribe to it
func (eb *EventBus) AddTopic(name string){
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    if _, ok := eb.Topics[name]; !ok{
        eb.Topics[name] = &Topic{
            Name: name,
            Subscribers: make(map[string]Subscriber),
            Messages: make(chan *Message, 100),
        }
    }
}

// creates a new topic and starts the dispatching process
func (eb *EventBus) CreateTopic(name string) *Topic{
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    topic := &Topic{
        Name: name,
        Subscribers: make(map[string]Subscriber),
        Messages: make(chan *Message, 100),
    }

    eb.Topics[name] = topic

    go topic.dispatch()

    return topic
}

//publishes a message to a specific topic
func (eb *EventBus) Publish(name string, msg *Message){
    eb.Mu.RLock()
    defer eb.Mu.RUnlock()

    if topic, ok := eb.Topics[name]; ok{
        select{
            case topic.Messages <- msg:
            //if the channel is full
            default:
            logger.GetLogger().Errorf("Message queue is full for topic: %s", name)
        }
    }
}
