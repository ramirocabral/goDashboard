package memory

import (
	"log"
	"strings"
	"time"

	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/utils"
)

type Memory struct{
    UsedPercentage  float64  `json:"used_percentage"`
    Type            string  `json:"type"`
    Frequency       uint64  `json:"frequency"` 
    Total           uint64  `json:"total"`
    Used            uint64  `json:"used"`
    Free            uint64  `json:"free"`
    Active          uint64  `json:"active"`
    Inactive        uint64  `json:"inactive"`
    Buffers         uint64  `json:"buffers"`
    Cached          uint64  `json:"cached"`
}

func (m Memory) ToPoint() []*core.Point{
    return []*core.Point{{
        Timestamp: time.Now(),
        Measurement: "memory",
        Tags: map[string]string{
            "type": m.Type,
            "frequency": utils.Uint64ToStr(m.Frequency),
        }, 
        Fields: map[string]interface{}{
            "used_percentage": m.UsedPercentage,
            "total": m.Total,
            "used": m.Used,
            "free": m.Free,
            "active": m.Active,
            "inactive": m.Inactive,
            "buffers": m.Buffers,
            "cached": m.Cached,
            },
        },
    }
}

const MEMORY_PATH = "/host/proc/meminfo"

func ReadMemory() (Memory, error){
    output := Memory{}

    memInfo, err := utils.ExecuteCommand("dmidecode", "-t", "17")

    if err != nil{
        log.Println("Error reading memory info: ", err)
        return Memory{}, err
    }

    for _, line := range strings.Split(memInfo, "\n"){
        if strings.Contains(line, "Type:"){
            output.Type = strings.TrimSpace(strings.Split(line, ":")[1])
        }
        if strings.HasPrefix(strings.TrimSpace(line), "Speed"){
            output.Frequency = utils.StrToUint64(strings.Fields(line)[1])
        }
    }

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
        output.UsedPercentage = utils.RoundFloat((float64(output.Used)/ float64(output.Total)) * 100,2)

    }

    return output, nil
}
