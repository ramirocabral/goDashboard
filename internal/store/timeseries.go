package store

import (
    // "context"
    // "database/sql"
    // "fmt"
    // "log"
    // "time"
)
//it should write cpu usage, memory usage, disk i/o usage and network usage
// just the interfaces

type Store interface{
    WriteCPUUsage(cpuUsage float64) error
    WriteMemoryUsage(memoryUsage float64) error
    WriteDiskIOUsage(diskIOUsage float64) error
    WriteNetworkUsage(networkUsage float64) error
}
