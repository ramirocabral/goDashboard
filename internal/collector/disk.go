package collector

import (
    "errors"
    "strings"

    "golang-system-monitor/internal/utils"
)

type Disk struct {
    Device          string
    Type            string
    MountPt         string  
    Size            uint64
    Used            uint64
    Free            uint64  
    UsedPercentage  uint64
}    

func ReadDisks() ([]Disk, error) {
    diskData, err := utils.ExecuteCommand("df","-T","-BG","--exclude-type=tmpfs","--exclude-type=devtmpfs","--exclude-type=cifs","--exclude-type=efivarfs")
    diskDataSplit := strings.Split(string(diskData), "\n")[1:]

    if err != nil {
        return nil, errors.New("error reading disk data")
    }

    var disks []Disk

    for _, line := range diskDataSplit {
        fields := strings.Fields(line)

        if len(fields) == 0 {
            continue
        }

        disk := Disk{
            Device:      fields[0],
            Type:        fields[1],
            Size:       utils.StrToUint64(strings.TrimSuffix(fields[2], "G")),
            Used:        utils.StrToUint64(strings.TrimSuffix(fields[3], "G")),
            Free:        utils.StrToUint64(strings.TrimSuffix(fields[4], "G")),
            UsedPercentage: utils.StrToUint64(strings.TrimSuffix(fields[5], "%")),
            MountPt:     fields[6],
        }

        disks = append(disks, disk)
    }

    return disks, nil
}
