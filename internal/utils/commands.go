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

func ExecuteCommandWithPipe(command string, params ...string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()
	if err != nil {
		return string(stdout), err
	}
	return string(stdout), nil
}
