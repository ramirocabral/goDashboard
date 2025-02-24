package storage

import (
    // "log"

    "golang-system-monitor/internal/core"
)

type StorageSubscriber struct{
    Id          string
    Storage     Storage
    Topics      map[string]*core.Topic
}

func NewStorageSubscriber(storage Storage) *StorageSubscriber{
    return &StorageSubscriber{
        Id: storage.ID(),
        Storage: storage,
        Topics: make(map[string]*core.Topic),
    }
}

func (ss *StorageSubscriber) ID() string{
    return ss.Id
}

func (ss *StorageSubscriber) Handle(msg core.Message){
    if _, ok := ss.Topics[msg.Type]; !ok{
        return
    }

    switch msg.Type{
    case "cpu":
        ss.Storage.WriteCpuStats(msg.Data)
    case "memory":
        ss.Storage.WriteMemoryStats(msg.MemoryStats)
    case "io":
        ss.Storage.WriteIOStats(msg.IOStats)
    case "network":
        ss.Storage.WriteNetworkStats(msg.NetworkStats)
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
