package main

import (
    "fmt"
    "os"

    "golang-system-monitor/internal/collector"
)

func main(){
    //test disks info
    smart, err := collector.ReadSmart("/dev/nvme0n1")

    if err != nil{
        fmt.Println("Error reading smart data: ", err)
        os.Exit(1)
    }

    fmt.Println("Device:", smart.Device)


    //print the map[string]string
    for key, value := range smart.Data{
        fmt.Println(key, ":", value)
    }

}
