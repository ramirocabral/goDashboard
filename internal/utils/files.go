package utils

import (
	"errors"
	"log"
	"strings"
	"os"
	"path/filepath"
)

func ReadFile(path string) (string, error){
    data, err := os.ReadFile(path)

    if err != nil {
	return "", err
    }

    return string(data), nil
}

// find the file that contains the CPU temp
func FindCPUTempFile() (string, error){
    const INTEL_SENSOR_NAME = "coretemp"
    const AMD_SENSOR_NAME = "k10temp"

    //get list of possible directories
    dirs, err := filepath.Glob("/sys/class/hwmon/hwmon*")

    if err != nil{
	log.Println("Error finding CPU temp file: ", err)
	return "", err
    }

    for _, dir := range dirs{
	//read the name file and check if it contains the sensor nameq
	name, err := ReadFile(filepath.Join(dir, "name"))
	if errors.Is(err, os.ErrNotExist){
	    continue
	}
	if err != nil{
	    log.Println("Error reading name file: ", err)
	    return "", err
	}

	nameStr := strings.TrimSpace(name)

	if nameStr == INTEL_SENSOR_NAME || nameStr == AMD_SENSOR_NAME{
	    return filepath.Join(dir, "temp1_input"), nil
	}
    }
    return "", errors.New("no CPU temp file found")
}
