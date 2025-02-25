package core

import (
    "context"
    "time"
    "sync"
)

type Publisher interface{
    Publish(Message)
    Start(context.Context)  error
    Stop()                  error
}

type Subscriber interface{
    ID()                    string
    Handle(*Message)
    Subscribe(*Topic)       error
    Unsubscribe(*Topic)     error
}

type Storable interface{
    ToMap() map[string]interface{}
}

type Message struct{
    Type            string
    Timestamp       time.Time
    Data            Storable
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

// adds a subscriber to the topic
func (t *Topic) AddSubscriber(sub Subscriber){
    t.Mu.Lock()
    defer t.Mu.Unlock()
    t.Subscribers[sub.ID()] = sub
}

// removes a subscriber from the topic
func (t *Topic) RemoveSubscriber(sub Subscriber){
    t.Mu.Lock()
    defer t.Mu.Unlock()

    delete(t.Subscribers, sub.ID())
}

//constructor
func NewEventBus() *EventBus{
    return &EventBus{
        Topics: make(map[string]*Topic),
    }
}

// dispatches every message received to all of its Subscribers
func (t *Topic) dispatch(){
    for msg := range t.Messages{
        t.Mu.RLock()
        for _, sub := range t.Subscribers{
            go sub.Handle(msg)
        }
        t.Mu.RUnlock()
    }
}

// add a new topic to the event bus, so the Subscribers can subscribe to it
func (eb *EventBus) AddTopic(name string){
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    if _, ok := eb.Topics[name]; !ok{
        eb.Topics[name] = &Topic{
            Name: name,
            Subscribers: make(map[string]Subscriber),
            Messages: make(chan *Message),
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
        Messages: make(chan *Message, 1000),
    }

    eb.Topics[name] = topic

    go topic.dispatch()

    return topic
}
