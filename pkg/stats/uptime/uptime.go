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

func (u *Uptime) ToPoint() []*core.Point{
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
