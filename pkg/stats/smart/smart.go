package smart

import (
	// "fmt"
	"log"
	"regexp"
	"strings"

	"golang-system-monitor/internal/utils"
	// "golang-system-monitor/internal/logger"

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
        return Smart{}, err
    }

    output := Smart{}

    for _, device := range devices{
        smartData, err := ReadData(device)

        if err != nil{
            return Smart{}, err
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

//returna a slice with the devices to get smart data from
func getDevices() ([]string, error){
    data, err := utils.ExecuteCommand("lsblk", "-d", "-o", "NAME", "-n", "-l")

    if err != nil{
        log.Println("Error getting disks: ", err)
        return nil, err
    }

    dataSplit := strings.Split(string(data), "\n")

    disks := []string{}

    for _, line := range dataSplit{

        if line == "" || strings.HasPrefix(line, "loop") || strings.HasPrefix(line, "sr") || strings.HasPrefix(line, "ram") || strings.HasPrefix(line, "zram"){
            continue
        }

        if strings.HasPrefix(line, "nvme"){
            re := regexp.MustCompile(`n\d+`)
            line = re.ReplaceAllString(line, "")
        }

        // fmt.Println(line)

        disks = append(disks, "/dev/" + line)
    }

    return disks, nil
}
