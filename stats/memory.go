package stats

import(
    "log"
    "strings"

    "golang-system-monitor/internal/utils"
)

type Memory struct{
    PercentageUsed float64
    Total uint64
    Used uint64
    Free uint64
    Active uint64
    Inactive uint64
    Buffers uint64
    Cached uint64
}

const MEMORY_PATH = "/proc/meminfo"

func ReadMemory() (Memory, error){
    output := Memory{}

    memData , err := utils.ReadFile(MEMORY_PATH)

    if err != nil{
        log.Println("Error reading memory data: ", err)
        return Memory{}, err
    }

    memDataSplit := strings.Split(string(memData), "\n")

    for _, line := range memDataSplit{
        fields := strings.Fields(line)

        if len(fields) == 0{
            continue
        }

        switch fields[0]{
        case "MemTotal:": 
            output.Total = utils.StrToUint64(fields[1])
        case "MemFree:": 
            output.Free = utils.StrToUint64(fields[1])
        case "Buffers:":
            output.Buffers = utils.StrToUint64(fields[1]) 
        case "Cached:":
            output.Cached = utils.StrToUint64(fields[1])
        case "Active:":
            output.Active = utils.StrToUint64(fields[1])
        case "Inactive:":
            output.Inactive = utils.StrToUint64(fields[1])
        }

        output.Used = output.Total - (output.Free + output.Cached + output.Buffers)
        output.PercentageUsed = (float64(output.Used)/ float64(output.Total)) * 100
    }

    return output, nil
}
