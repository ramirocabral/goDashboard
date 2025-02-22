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
    Handle(Message)
    Subscribe(string)       error
    Unsubscribe(string)     error
}

type Message struct{
    Type            string
    Timestamp       time.Time
    Data            interface{}
}


type Topic struct{
    Name        string
    Subscribers map[string]Subscriber
    Mu          sync.RWMutex
    Messages    chan Message
}


type EventBus struct{
    Topics  map[string]*Topic
    Mu      sync.RWMutex
}

func (t *Topic) AddSubscriber(sub Subscriber){
    t.Mu.Lock()
    defer t.Mu.Unlock()

    t.Subscribers[sub.ID()] = sub
}

func NewEventBus() *EventBus{
    return &EventBus{
        Topics: make(map[string]*Topic),
    }
}

//dispatches every message received to all of its Subscribers
func (t *Topic) dispatch(){
    for msg := range t.Messages{
        t.Mu.RLock()
        for _, sub := range t.Subscribers{
            go sub.Handle(msg)
        }
        t.Mu.RUnlock()
    }
}

//add a new topic to the event bus, so the Subscribers can subscribe to it
func (eb *EventBus) AddTopic(name string){
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    if _, ok := eb.Topics[name]; !ok{
        eb.Topics[name] = &Topic{
            Name: name,
            Subscribers: make(map[string]Subscriber),
            Messages: make(chan Message),
        }
    }
}

//remove a subscriber from every topic on the event bus
func (eb *EventBus) RemoveSubscriber(id string){
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    for _, topic := range eb.Topics{
        topic.Mu.Lock()
        delete(topic.Subscribers, id)
        topic.Mu.Unlock()
    }
}

func (eb *EventBus) CreateTopic(name string) *Topic{
    eb.Mu.Lock()
    defer eb.Mu.Unlock()

    topic := &Topic{
        Name: name,
        Subscribers: make(map[string]Subscriber),
        Messages: make(chan Message, 1000),
    }

    eb.Topics[name] = topic

    go topic.dispatch()

    return topic
}
