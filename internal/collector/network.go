package collector

import (
	"log"
	"strings"

	"golang-system-monitor/internal/utils"
)

type Network struct {
    Interface string        `json:"interface"`
    Ip        string        `json:"ip"`
    Usage     NetworkUsage  `json:"usage"`
}

type NetworkUsage struct {
    State        string `json:"state"`
    RxBytes      string `json:"rx_per_second"`
    TxBytes      string `json:"tx_per_second"`
}


func GetNetworks() ([]Network, error){
    output := []Network{}

    command := "ip -o addr show scope global | awk '{split($4, a, \"/\"); print $2\" : \"a[1]}'"

    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        return []Network{}, err
    }

    dataSplit := strings.Split(string(data), "\n")

    for _, iface := range dataSplit{

        if strings.HasPrefix(iface, "lo") || strings.HasPrefix(iface, "docker") || strings.HasPrefix(iface, "br") || iface == ""{
            continue
        }

        net := Network{}

        ifaceSplit := strings.Split(iface, " : ")
        net.Interface = ifaceSplit[0]
        net.Ip = ifaceSplit[1]
        net.Usage = getNetworkUsage(net.Interface)

        output = append(output, net)

    }

    return output, nil
}

func getNetworkUsage(interfaceName string) NetworkUsage{
    output := NetworkUsage{}

    command := "cat /proc/net/dev | grep " + interfaceName + " | awk '{print $1\" \"$2\" \"$10}'"
    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        log.Println(err)
    }

    dataSplit := strings.Split(string(data), " ")
    output.RxBytes = dataSplit[1]
    output.TxBytes = dataSplit[2]

    return output
}

