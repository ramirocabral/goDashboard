package subscribers

import (
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/storage"
)

type StorageSubscriber struct{
    Id          string
    Storage     storage.Storage
    Topics      map[string]*core.Topic
}

func NewStorageSubscriber(storage storage.Storage) *StorageSubscriber{
    return &StorageSubscriber{
        Id: storage.ID(),
        Storage: storage,
        Topics: make(map[string]*core.Topic),
    }
}

func (ss *StorageSubscriber) ID() string{
    return ss.Id
}

func (ss *StorageSubscriber) Handle(msg *core.Message){
    if _, ok := ss.Topics[msg.Type]; !ok{
        return
    }

    point := msg.Data.ToPoint()

    err := ss.Storage.WriteStats(point)
    if err != nil{
        return
    }
}

func (ss *StorageSubscriber) Subscribe(topic *core.Topic) error{
    if _, ok := ss.Topics[topic.Name]; ok{
        return nil
    }

    ss.Topics[topic.Name] = topic
    topic.AddSubscriber(ss)
    return nil
}

func (ss *StorageSubscriber) Unsubscribe(topic *core.Topic) error{
    if _, ok := ss.Topics[topic.Name]; !ok{
        return nil
    }

    delete(ss.Topics, topic.Name)
    topic.RemoveSubscriber(ss)
    return nil
}
