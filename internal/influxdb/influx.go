package influxdb

import(
    "time"
    "context"
    "fmt"

    "golang-system-monitor/internal/core"
    "golang-system-monitor/internal/storage"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct{
    client influxdb2.Client 
    writeAPI api.WriteAPI
    queryAPI api.QueryAPI
}


func New(addr, token string) (storage.Storage, error) {
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

//TODO: Change and receive a point
func (s *InfluxStore) WriteStats(p *core.Point){

    point := influxdb2.NewPoint(
            p.Measurement,
            p.Tags,
            p.Fields,
            p.Timestamp)

    s.writeAPI.WritePoint(point)
}

func (s *InfluxStore) ReadCpuStats(startTime, endTime time.Time) (storage.CPUResponse, error){

        query := fmt.Sprintf(`
        from(bucket: "my-bucket")
            |> range(start: %s, stop: %s)
            |> filter(fn: (r) => r._measurement == "cpu")
            |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
            |> keep(columns: ["_time", "model_name", "temp", "usage_percentage", "idle_percentage"]) `, startTime.UTC().Format(time.RFC3339), endTime.UTC().Format(time.RFC3339))

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
