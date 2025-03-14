package storage

import(
    "time"
)

//models for getting historical data
type CPUResponse struct{
    ModelName   string          `json:"model_name"`
    Frequency   uint64          `json:"frequency"`
    Family      string          `json:"family"`
    Cores       uint64          `json:"cores"`
    Threads     uint64          `json:"threads"`
    Stats       []CPUPoint      `json:"stats"`
}

type CPUPoint struct{
    Timestamp       time.Time   `json:"timestamp"`
    UsagePercentage float64     `json:"usage_percentage"`
    IdlePercentage  float64     `json:"idle_percentage"`
}

type IOResponse struct{
    Devices    []IOStats       `json:"devices"`
}

type IOStats struct{
    Device          string      `json:"device"`
    Stats           []IOPoint        `json:"stats"`
}

type IOPoint struct{
    Timestamp           time.Time   `json:"timestamp"`
    KBReadPerSecond     uint64      `json:"kb_read_per_second"`
    KBWritePerSecond    uint64      `json:"kb_write_per_second"`
}

type MemoryResponse struct{
    Type        string          `json:"type"`
    Frequency   uint64          `json:"frequency"`
    Stats       []MemoryPoint   `json:"stats"`
}

type MemoryPoint struct{
    Timestamp       time.Time       `json:"timestamp"`
    UsedPercentage  float64         `json:"used_percentage"`
    Total           uint64          `json:"total"`
    Used            uint64          `json:"used"`
    Free            uint64          `json:"free"`
    Active          uint64          `json:"active"`
    Inactive        uint64          `json:"inactive"`
    Buffers         uint64          `json:"buffers"`
    Cached          uint64          `json:"cached"`
}

type NetworkResponse struct{
    Interfaces        []NetworkStats  `json:"interfaces"`
}

type NetworkStats struct{
    Interface       string          `json:"interface"`
    Stats            []NetworkPoint `json:"stats"` 
} 

type NetworkPoint struct{
    Timestamp       time.Time       `json:"timestamp"`
    RxBytesPS       uint64          `json:"rx_bytes_ps"`
    TxBytesPS       uint64          `json:"tx_bytes_ps"`
}
