package smart

import (
	"log"
	"strings"

	"golang-system-monitor/internal/utils"

	"github.com/anatol/smart.go"
)

type Smart struct{
    Devices     []SmartData     `json:"devices"`
}

type SmartData struct{
    Device  string              `json:"device"`
    Data    map[string]string   `json:"data"`
}

func ReadSmart() (Smart, error){
    devices, err := getDevices()

    if err != nil{
        log.Println("Error getting devices: ", err)
        return Smart{}, err
    }

    output := Smart{}

    for _, device := range devices{
        smartData, err := ReadData(device)

        if err != nil{
            log.Println("Error reading smart data: ", err)
            continue
        }

        output.Devices = append(output.Devices, smartData)
    }

    return output, nil
}


func ReadData(device string) (SmartData, error){

    dev, err := smart.Open(device)
    if err != nil{
        log.Println("Error opening device: ", err)
        return SmartData{}, err
    }

    defer dev.Close()


    output := SmartData{}

    smartData, err := utils.ExecuteCommand("smartctl", "-A", device)

    if err != nil{
        log.Println("Error reading smart data: ", err)
        return SmartData{}, err
    }

    output.Device = device

    output.Data = make(map[string]string)
    
    smartDataSplit := strings.Split(string(smartData), "\n")[5:]

    switch dev.(type){
        case *smart.NVMeDevice:
            output.Data = readNvmeSmart(smartDataSplit)
        case *smart.SataDevice:
            output.Data = readSataSmart(smartDataSplit)
        default:
            log.Println("Unknown device type")
    }


    return output, nil
}

func readNvmeSmart(dataSplit []string) map[string]string{

    output := make(map[string]string)

    for _, line := range dataSplit{
        fields := strings.Split(string(line), ":")

        if len(fields) < 2{
            continue
        }

        output[fields[0]] = fields[1]
    }

    return output
}

func readSataSmart(dataSplit []string) map[string]string{

    output := make(map[string]string)

    for _, line := range dataSplit{
        fields := strings.Fields(string(line))

        if len(fields) < 10{
            continue
        }

        output[fields[2]] = strings.Join(fields[2:], " ")
    }

    return output
}

func getDevices() ([]string, error){
    data, err := utils.ExecuteCommand("smartctl" , "--scan")

    if err != nil{
        log.Println("Error getting disks: ", err)
        return nil, err
    }

    dataSplit := strings.Split(string(data), "\n")

    disks := []string{}

    for _, line := range dataSplit{
        fields := strings.Fields(line)

        if len(fields) < 2{
            continue
        }

        disks = append(disks, fields[0])
    }

    return disks, nil
}
