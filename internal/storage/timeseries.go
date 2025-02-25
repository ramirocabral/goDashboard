package storage

import (
    "fmt"
    "context"
    "os"
    "golang-system-monitor/internal/stats/cpu"
    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/stats/io"
    "golang-system-monitor/internal/stats/memory"
    "golang-system-monitor/internal/stats/network"
    "time"

    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct{
    client influxdb2.Client 
    writeAPI api.WriteAPI
    queryAPI api.QueryAPI
}

type CPUStats struct{
    Timestamp       time.Time   `json:"timestamp"`
    CPUInfo         cpu.CPU     `json:"cpu_info"`
}

type MemoryStats struct{
    Timestamp       time.Time   `json:"timestamp"`
    MemoryInfo      memory.Memory   `json:"memory_info"`
}

type IOStats struct{
    Timestamp       time.Time   `json:"timestamp"`
    IOInfo          io.DiskIO   `json:"io_info"`
}

type NetworkStats struct{
    Timestamp       time.Time   `json:"timestamp"`
    NetworkInfo     network.Network `json:"network_info"`
}

type Storage interface{
    ID() string
    WriteStats(m *core.Message)
    ReadCpuStats(startTime, endTime time.Time) ([]CPUStats, error)
    ReadIOStats(startTime, endTime time.Time) ([]IOStats, error)
}

func New(addr, token string) (Storage, error) {
    client := influxdb2.NewClient(addr, token)
    writeAPI := client.WriteAPI("my-org", "my-bucket")
    queryAPI := client.QueryAPI("my-org")
    return &InfluxStore{
        client: client,
        writeAPI: writeAPI,
        queryAPI: queryAPI, }, nil
}

func (s *InfluxStore) ID() string{
    return s.client.ServerURL()
}

func (s *InfluxStore) WriteStats(m *core.Message){
    point := influxdb2.NewPoint(m.Type,
            map[string]string{
            "host": os.Getenv("HOST"),
            },
            m.Data.ToMap(),
            time.Now(),
    )

    s.writeAPI.WritePoint(point)
}


func (s *InfluxStore) ReadCpuStats(startTime time.Time, endTime time.Time) ([]CPUStats, error){

        query := fmt.Sprintf(`
        from(bucket: "my-bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "cpu")
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> keep(columns: ["_time", "model_name", "cores", "threads", "temp", "usage_percentage", "idle_percentage"]) `, startTime.UTC().Format(time.RFC3339), endTime.UTC().Format(time.RFC3339))

    result, err := s.queryAPI.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }

    stats := parseCpuStats(result)

    if result.Err() != nil {
        return nil, result.Err()
    }

    return stats, nil
}

func parseCpuStats(result *api.QueryTableResult) []CPUStats{

    var stats []CPUStats

    for result.Next() {
        values := result.Record().Values()

        cpu := cpu.CPU{
            ModelName: values["model_name"].(string),
            Cores:     uint64(values["cores"].(uint64)),
            Threads:   uint64(values["threads"].(uint64)),
            Temp:      uint64(values["temp"].(uint64)),
            UsageStatistics: cpu.Usage{
                UsagePercentage: values["usage_percentage"].(float64),
                IdlePercentage:  values["idle_percentage"].(float64),
            },
        }

        stat := CPUStats{
            Timestamp: values["_time"].(time.Time),
            CPUInfo:       cpu,
        }

        stats = append(stats, stat)
    }

    return stats
}

// type DiskIO struct {
//     Device          string      `json:"device"`
//     ReadPerSecond   uint64      `json:"kb_read_per_second"`
//     WritePerSecond  uint64      `json:"kb_write_per_second"`
// }

// func (d *DiskIO) ToMap() map[string]interface{}{
//     return map[string]interface{}{
//         "device": d.Device,
//         "kb_read_per_second": d.ReadPerSecond,
//         "kb_write_per_second": d.WritePerSecond,
//     }
// }

func (s *InfluxStore) ReadIOStats(startTime, endTime time.Time) ([]IOStats, error){
    
        query := fmt.Sprintf(`
        from(bucket: "my-bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "io")
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> keep(columns: ["_time", "device", "kb_read_per_second", "write_bytes"]) `, startTime.UTC().Format(time.RFC3339), endTime.UTC().Format(time.RFC3339))

    
        result, err := s.queryAPI.Query(context.Background(), query)
        if err != nil {
            return nil, err
        }
    
        stats := parseIOStats(result)
    
        if result.Err() != nil {
            return nil, result.Err()
        }
    
        return stats, nil
}

func parseIOStats(result *api.QueryTableResult) []IOStats{
         
    var stats []IOStats 

    for result.Next() {
        values := result.Record().Values()

        io := io.DiskIO{
            Device: values["device"].(string),
            ReadPerSecond: uint64(values["read_bytes"].(uint64)),
            WritePerSecond: uint64(values["write_bytes"].(uint64)),
        }

        stat := IOStats{
            Timestamp: values["_time"].(time.Time),
            IOInfo:       io,
        }
        stats = append(stats, stat)
    }

    return stats
} 
