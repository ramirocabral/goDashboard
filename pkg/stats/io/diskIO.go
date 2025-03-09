package io

import (
	"strings"
	"time"

	"go-dashboard/internal/utils"
	"go-dashboard/internal/core"
	"go-dashboard/internal/logger"
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

const DISK_STATS_PATH = "/host/proc/diskstats"


var lastDiskData = map[string]BytesStore{}

//this function is called every 1 second so the stats are actually accurate
func ReadDiskIO() (DiskIO, error) {
    var disks DiskIO

    data, err := utils.ReadFile(DISK_STATS_PATH)

    if err != nil {
        logger.GetLogger().Error("Error reading disk stats: ", err)
        return nil, err
    }

    lines := strings.Split(string(data), "\n")

    for _, line := range lines {
        fields := strings.Fields(line)

        if len(fields) < 14  || fields[1] != "0" || strings.HasPrefix(fields[2],"zram") || strings.HasPrefix(fields[2],"loop") || strings.HasPrefix(fields[2],"ram") {
            continue
        }

        devName := fields[2]

        readSectors := utils.StrToUint64(fields[5])
        writeSectors := utils.StrToUint64(fields[9])

        readKb := readSectors * 512 / 1024
        writeKb := writeSectors * 512 / 1024


        disk := Disk{Device: devName}

        if prev, exists := lastDiskData[devName]; exists {
            disk.ReadPerSecond = readKb - prev.ReadBytes
            disk.WritePerSecond = writeKb - prev.WriteBytes
        }else{
            disk.ReadPerSecond = 0
            disk.WritePerSecond = 0
        }

        lastDiskData[devName] = BytesStore{ReadBytes: readKb, WriteBytes: writeKb}

        disks = append(disks, disk)

    }
    return disks, nil
}
