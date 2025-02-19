package collector

import (
	"errors"
	"strings"

	"golang-system-monitor/internal/utils"
)

// Memory struct

type Uptime struct{
    Uptime uint64 //uptime in seconds
}

const UPTIME_PATH = "/proc/uptime"

// ReadUptime function
func ReadUptime() (Uptime, error){
    output := Uptime{}

    uptimeData, err := utils.ReadFile(UPTIME_PATH)

    if err != nil{
        return Uptime{}, errors.New("error reading uptime data")
    }

    uptimeDataSplit := strings.Fields(string(uptimeData))

    output.Uptime = utils.StrToUint64(uptimeDataSplit[0])

    return output, nil
}
