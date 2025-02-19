package utils

import (
	"os/exec"
)

func ExecuteCommand(command string, args ...string) (string, error){
    cmd := exec.Command(command, args...)
    stdout, err := cmd.Output()

    if err != nil{
	return "", err
    }

    return string(stdout), nil
}
