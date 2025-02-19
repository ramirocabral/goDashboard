package utils

import (
    "strconv"
)

func StrToUint64(s string) uint64{
    value, err := strconv.ParseUint(s, 10 , 64)

    if err != nil {
        panic(err)
    }

    return value
}

func StrToFloat64(s string) float64{
    value, err := strconv.ParseFloat(s, 64)

    if err != nil {
        panic(err)
    }

    return value
}
