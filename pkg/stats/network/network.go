package network

import (
	"log"
	"strings"

	"golang-system-monitor/internal/utils"
	"golang-system-monitor/internal/core"
)

type Networks []Network

type Network struct {
    Interface       string        `json:"interface"`
    Ip              string        `json:"ip"`
    Usage           NetworkUsage  `json:"usage"`
}

type NetworkUsage struct {
    RxBytesPS      uint64 `json:"rx_bytes_ps"`
    TxBytesPS      uint64 `json:"tx_bytes_ps"`
}

type ByteStore struct {
    RxBytes uint64      //received bytes on last check
    TxBytes uint64      //transmitted bytes on last check
}

func (n Networks) ToPoint() []*core.Point{
    var points []*core.Point

    for _, network := range n{
        point := &core.Point{
            Measurement: "network",
            Tags: map[string]string{
                "interface": network.Interface,
                "ip": network.Ip,
            },
            Fields: map[string]interface{}{
                "rx_bytes_ps": network.Usage.RxBytesPS,
                "tx_bytes_ps": network.Usage.TxBytesPS,
            },
        }
        points = append(points, point)
    }

    return points
}

var lastNetworkData = map[string]ByteStore{}

func ReadNetworks() (Networks, error){
    output := Networks{}

    command := "ip -o addr show scope global | awk '{split($4, a, \"/\"); print $2\" : \"a[1]}'"

    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        return output, err
    }

    dataSplit := strings.Split(string(data), "\n")

    for _, iface := range dataSplit{

        if strings.HasPrefix(iface, "lo") || strings.HasPrefix(iface, "docker") || strings.HasPrefix(iface, "br") || iface == ""{
            continue
        }

        net := Network{}

        //get the interface name an ip addr
        ifaceSplit := strings.Split(iface, " : ")
        net.Interface = ifaceSplit[0]
        net.Ip = ifaceSplit[1]

        //get the rw bytes per second
        bytes := getNetworkBytes(net.Interface)

        usage := NetworkUsage{}

        //if it's the first iteration, we can't calculate the bytes per second
        if _, ok := lastNetworkData[net.Interface]; !ok{
            usage.RxBytesPS = 0
            usage.TxBytesPS = 0
        } else{
            //get the bytes per second, subtract the last iteration 
            usage.RxBytesPS = bytes.RxBytes - lastNetworkData[net.Interface].RxBytes
            usage.TxBytesPS = bytes.TxBytes - lastNetworkData[net.Interface].TxBytes
        }

        net.Usage = usage

        //update the last network data
        lastNetworkData[net.Interface] = bytes

        //add the actual network to the slice
        output = append(output, net)
    }

    return output, nil
}

//return the bytes of the actual iteration
func getNetworkBytes(interfaceName string) ByteStore{
    output := ByteStore{}

    command := "cat /proc/net/dev | grep " + interfaceName + " | awk '{print $1\" \"$2\" \"$10}'"
    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        log.Println(err)
    }

    dataSplit := strings.Fields(string(data))
    output.RxBytes = utils.StrToUint64(dataSplit[1])
    output.TxBytes = utils.StrToUint64(dataSplit[2])

    return output
}

