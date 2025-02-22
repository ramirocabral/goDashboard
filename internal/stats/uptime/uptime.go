package uptime

import (
	"errors"
	"strings"

	"golang-system-monitor/internal/utils"
)


type Uptime struct{
    Uptime uint64   `json:"uptime"`
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
