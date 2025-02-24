package storage

import (
    "log"
    "fmt"
    "context"
    "os"
    "golang-system-monitor/internal/stats/cpu"
    "golang-system-monitor/internal/core"
    // "golang-system-monitor/internal/stats/io"
    // "golang-system-monitor/internal/stats/memory"
    // "golang-system-monitor/internal/stats/network"
    "time"

    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct{
    client influxdb2.Client 
    writeAPI api.WriteAPI
    queryAPI api.QueryAPI
}

type Storage interface{
    ID() string
    WriteStats(m *core.Message)
    ReadCpuStats(startTime time.Time, endTime time.Time) ([]*cpu.CPU)
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

func (s *InfluxStore) ReadCpuStats(startTime time.Time, endTime time.Time) ([]*cpu.CPU){
    //TODO: use sprintf
    query := `from(bucket:"my-bucket")
            |> range(start:`+ startTime.Format(time.RFC3339) +`, stop:`+ endTime.Format(time.RFC3339) +` )
            |> filter(fn: (r) => r._measurement == "cpu")`

    result, err := s.queryAPI.Query(context.Background(), query)

    if err != nil {
        log.Println("Error reading cpu stats: ", err)
        return nil
    }

    for result.Next() {
        fmt.Println(result.Record().Values()["cores"])
        fmt.Println(result.Record().Values()["usage_percentage"])
    }

    return nil
}
