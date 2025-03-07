package network

import (
	"log"
	"strings"

	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/utils"
)

type Networks []Network

type Network struct {
    Interface       string        `json:"interface"`
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

    // command := "ip -o addr show scope global | awk '{split($4, a, \"/\"); print $2\" : \"a[1]}'"
    command := "ls /host/sys/class/net"

    data, err := utils.ExecuteCommandWithPipe(command)
    if err != nil{
        return output, err
    }

    ifaces := strings.Fields(string(data))

    for _, iface := range ifaces{

        if strings.HasPrefix(iface, "docker") || strings.HasPrefix(iface, "br") || iface == "" || strings.HasPrefix(iface, "veth") {
            continue
        }

        net := Network{}

        //get the interface name an ip addr
        net.Interface = iface

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

    command1 := "cat /host/sys/class/net/" + interfaceName + "/statistics/rx_bytes"
    command2 := "cat /host/sys/class/net/" + interfaceName + "/statistics/tx_bytes"

    data, err := utils.ExecuteCommandWithPipe(command1)
    if err != nil{
        log.Println(err)
    }

    data = utils.TrimNewLine(data)
    output.RxBytes = utils.StrToUint64(string(data))

    data, err = utils.ExecuteCommandWithPipe(command2)
    if err != nil{
        log.Println(err)
    }

    data = utils.TrimNewLine(data)
    output.TxBytes = utils.StrToUint64(string(data))

    // dataSplit := strings.Fields(string(data))
    // output.RxBytes = utils.StrToUint64(dataSplit[1])
    // output.TxBytes = utils.StrToUint64(dataSplit[2])

    return output
}

