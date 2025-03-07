package uptime

import (
	"errors"
	"strings"
	"time"

	"golang-system-monitor/internal/utils"
	"golang-system-monitor/internal/core"
)


type Uptime struct{
    Uptime uint64   `json:"uptime"`
}

func (u Uptime) ToPoint() []*core.Point{
	return []*core.Point{{
		Timestamp: time.Now(),
		Measurement: "uptime",
		Tags: map[string]string{},
		Fields: map[string]interface{}{
			"uptime": u.Uptime,
		},
	} ,
	}
}

const UPTIME_PATH = "/host/proc/uptime"

func ReadUptime() (Uptime, error){
    output := Uptime{}

    uptimeData, err := utils.ReadFile(UPTIME_PATH)

    if err != nil{
        return Uptime{}, errors.New("error reading uptime data")
    }

    uptimeDataSplit := strings.Fields(string(uptimeData))

    output.Uptime = uint64(utils.RoundFloat(utils.StrToFloat64(uptimeDataSplit[0]),0))

    return output, nil
}
