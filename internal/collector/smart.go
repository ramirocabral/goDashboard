package collector

import(
    "log"
    "strings"

    "golang-system-monitor/internal/utils"
    "github.com/anatol/smart.go"
)

type Smart struct{
    Device  string              `json:"device"`
    Data    map[string]string   `json:"data"`
}

func ReadSmart(device string) (Smart, error){

    dev, err := smart.Open(device)
    if err != nil{
        log.Println("Error opening device: ", err)
        return Smart{}, err
    }

    defer dev.Close()


    output := Smart{}

    smartData, err := utils.ExecuteCommand("smartctl", "-A", device)

    if err != nil{
        log.Println("Error reading smart data: ", err)
        return Smart{}, err
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

