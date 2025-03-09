package disk

import (
    "errors"
    "strings"

    "go-dashboard/internal/utils"
)

type Disks struct{
    Disks []Disk `json:"disks"`
}

type Disk struct {
    Device          string      `json:"device"`
    Type            string      `json:"type"`
    MountPt         string      `json:"mount_point"`
    UsedPercentage  uint64      `json:"used_percentage"`
    GBSize          uint64      `json:"gb_size"`
    GBUsed          uint64      `json:"gb_used"`
    GBFree          uint64      `json:"gb_free"`
}    

func ReadDisks() (Disks, error) {
    diskData, err := utils.ExecuteCommand("df","-T","-BG","--exclude-type=tmpfs","--exclude-type=devtmpfs","--exclude-type=cifs","--exclude-type=efivarfs","--exclude-type=overlay")
    diskDataSplit := strings.Split(string(diskData), "\n")[1:]

    if err != nil {
        return Disks{}, errors.New("error reading disk data")
    }

    var disks Disks

    for _, line := range diskDataSplit {
        fields := strings.Fields(line)

        if len(fields) == 0 {
            continue
        }

        disk := Disk{
            Device:      fields[0],
            Type:        fields[1],
            UsedPercentage: utils.StrToUint64(strings.TrimSuffix(fields[5], "%")),
            GBSize:       utils.StrToUint64(strings.TrimSuffix(fields[2], "G")),
            GBUsed:        utils.StrToUint64(strings.TrimSuffix(fields[3], "G")),
            GBFree:        utils.StrToUint64(strings.TrimSuffix(fields[4], "G")),
            MountPt:     fields[6],
        }

        disks.Disks = append(disks.Disks, disk)
    }

    return disks, nil
}
