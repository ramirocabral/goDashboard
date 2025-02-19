package main

import (
    "fmt"
    "os"

    "golang-system-monitor/internal/collector"
)

func main(){
    //test disks info
    disks, err := collector.GetDisksInfo()

    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }

    for _, disk := range disks{
        fmt.Println("Device: ", disk.Device)
        fmt.Println("Type: ", disk.Type)
        fmt.Println("Size: ", disk.Size)
        fmt.Println("Used: ", disk.Used)
        fmt.Println("Free: ", disk.Free)
        fmt.Println("UsedPercentage: ", disk.UsedPercentage)
        fmt.Println("MountPt: ", disk.MountPt)
        fmt.Println("====================================")
    }

}
