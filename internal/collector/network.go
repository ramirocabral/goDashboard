package collector

import(
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
    State       string `json:"state"`
    RxBytesPS   string `json:"rx_per_second"`
    TxBytesPs   string `json:"tx_per_second"`
}

func GetNetworks() ([]Network, error){
    output := Network{}

    command := "ip -o addr show scope global | awk '{split($4, a, \"/\"); print $2\" : \"a[1]}'"

    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        return []Network{}, err
    }
    dataSplit := strings.Split(string(data), "\n")

    for _, iface := range dataSplit{
        interfaceSplit := strings.Split(iface, " : ")
        output.Interface = interfaceSplit[0]
        output.Ip = interfaceSplit[1]

        output.Usage = getNetworkUsage(output.Interface)

    }


    return output, nil
}

func getNetworkUsage(interface string) NetworkUsage{
    
    output := NetworkUsage{}

    command := "cat /proc/net/dev | grep " + interface + " | awk '{print $2\" \"$10}'"
    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        log.Println("Error reading network usage: ", err)
        return NetworkUsage{}
    }

    dataSplit := strings.Split(string(data), " ")

    output.RxPerSecond = dataSplit[0]
    output.TxPerSecond = dataSplit[1]

    return output

}
