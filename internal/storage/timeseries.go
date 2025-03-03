package storage

import (
    "time"
    "golang-system-monitor/internal/core"
)

// we do the parsing on the repository layer because of the simplicity of the project and the raw format of the data returned by influxdb
type Storage interface{
    ID() string
    WriteStats(points []*core.Point)
    ReadCpuStats(startTime, endTime time.Time, interval string) (CPUResponse, error)
    ReadIOStats(startTime, endTime time.Time, interval string) (IOResponse, error)
    ReadMemoryStats(startTime, endTime time.Time, interval string) (MemoryResponse, error)
    ReadNetworkStats(startTime, endTime time.Time, interval string) (NetworkResponse, error)
    Close()
}
