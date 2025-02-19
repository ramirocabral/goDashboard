package main

import (
    "fmt"
    "os"

    "golang-system-monitor/internal/utils"
)

func main(){
    // mem, err := stats.ReadMemory("/proc/meminfo")
    //
    // if err != nil{
    //     fmt.Println(err)
    // }
    //
    // fmt.Printf("PercentageUsed: %v\n Total:%v\n Used:%v\n Free:%v\n Active:%v\n Inactive:%v\n Buffers:%v\n Cached:%v\n ", mem.PercentageUsed, mem.Total, mem.Used, mem.Free, mem.Active, mem.Inactive, mem.Buffers, mem.Cached)
    // cpu, err := stats.ReadCpu("/proc/stat")
    cpu_temp_dir, err := utils.FindCPUTempFile()

    fmt.Println(cpu_temp_dir)

    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }

    cpu_temp, err := utils.ReadFile(cpu_temp_dir)

    fmt.Println(string(cpu_temp))
}
