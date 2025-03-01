package stats

import (
    "golang-system-monitor/pkg/stats/container"
    "golang-system-monitor/pkg/stats/cpu"
    "golang-system-monitor/pkg/stats/disk"
    "golang-system-monitor/pkg/stats/host"
    "golang-system-monitor/pkg/stats/io"
    "golang-system-monitor/pkg/stats/memory"
    "golang-system-monitor/pkg/stats/network"
    "golang-system-monitor/pkg/stats/smart"
    "golang-system-monitor/pkg/stats/uptime"
)

type StatsManager struct {}

func NewStatsManager() *StatsManager {
    return &StatsManager{}
}

func (sm *StatsManager) GetContainers() (container.Containers, error) {
    return container.ReadContainers()
}

func (sm *StatsManager) GetCPU() (cpu.CPU, error) {
    return cpu.ReadCPU()
}

func (sm *StatsManager) GetDisk() (disk.Disks, error) {
    return disk.ReadDisks()
}

func (sm *StatsManager) GetHost() (host.Host, error) {
    return host.ReadHost()
}

func (sm *StatsManager) GetIO() (io.DiskIO, error) {
    return io.ReadDiskIO()
}

func (sm *StatsManager) GetMemory() (memory.Memory, error) {
    return memory.ReadMemory()
}

func (sm *StatsManager) GetNetwork() (network.Networks, error) {
    return network.ReadNetworks()
}

func (sm *StatsManager) GetSMART() (smart.Smart , error) {
    return smart.ReadSmart()
}

func (sm *StatsManager) GetUptime() (uptime.Uptime, error) {
    return uptime.ReadUptime()
}
