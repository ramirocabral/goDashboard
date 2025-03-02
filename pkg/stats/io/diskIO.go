package io

//get write/read per second on disks
import (
	"errors"
	"strings"
	"time"

	"golang-system-monitor/internal/utils"
	"golang-system-monitor/internal/core"
)

type DiskIO []Disk

type Disk struct {
    Device          string      `json:"device"`
    ReadPerSecond   uint64      `json:"kb_read_per_second"`
    WritePerSecond  uint64      `json:"kb_write_per_second"`
}

type BytesStore struct{
    ReadBytes   uint64
    WriteBytes  uint64
}

func (d DiskIO) ToPoint() []*core.Point{

    var points []*core.Point

    for _, disk := range d{
        point := &core.Point{
            Timestamp: time.Now(),
            Measurement: "io",
            Tags: map[string]string{
                "device": disk.Device,
            },
            Fields: map[string]interface{}{
                "read_bytes": disk.ReadPerSecond,
                "write_bytes": disk.WritePerSecond,
            },
        }
        points = append(points, point)
    }
    return points
}


var lastDiskData = map[string]BytesStore{}

//this function is called every 1 second so the stats are actually accurate
func ReadDiskIO() (DiskIO, error) {
    diskData, err := utils.ExecuteCommand("iostat", "-d", "-x")

    if err != nil {
        return nil, errors.New("error reading disk data")
    }
    diskDataSplit := strings.Split(string(diskData), "\n")[3:]

    var disks DiskIO

    for _, line := range diskDataSplit {

        fields := strings.Fields(line)

        // if the line is invalid, or the device contains "zram"
        if len(fields) == 0 || strings.HasPrefix(fields[0], "zram") || strings.HasPrefix(fields[0], "Device") {
            continue
        }
        devName := fields[0]
        rxBytes := utils.StrToUint64(fields[2])
        txBytes := utils.StrToUint64(fields[8])

        disk := Disk{}

        disk.Device = devName
        //get the bytes per second
        
        //if it's the first iterarion, we can't calculate the iops
        if lastDiskData[devName].ReadBytes == 0{
            disk.ReadPerSecond = 0
            disk.WritePerSecond = 0
        }else{
            disk.ReadPerSecond = (rxBytes - lastDiskData[devName].ReadBytes)
            disk.WritePerSecond = (txBytes - lastDiskData[devName].WriteBytes)
        }

        lastDiskData[devName] = BytesStore{ReadBytes: rxBytes, WriteBytes: txBytes}

        disks = append(disks, disk)
    }

    return disks, nil
}

