package main

import (
    "fmt"
    "os"

    "golang-system-monitor/internal/collector"
)

func main(){
    containers, err := collector.GetContainers()

    if err != nil{
        fmt.Println("Error getting containers: ", err)
        os.Exit(1)
    }

    for _, container := range containers{
        fmt.Println(container.Name)
        fmt.Println(container.Uptime)
        fmt.Println(container.Image)
        fmt.Println(container.Status)
    }


    // networks, err := collector.GetNetworks()
    //
    // if err != nil{
    //     fmt.Println("Error getting networks: ", err)
    //     os.Exit(1)
    // }
    //
    // for _, network := range networks{
    //     fmt.Println(network.Interface)
    //     fmt.Println(network.Ip)
    //     fmt.Println(network.Usage.RxBytes)
    //     fmt.Println(network.Usage.TxBytes)
    //     fmt.Println(network.Usage.State)
    // }

}
