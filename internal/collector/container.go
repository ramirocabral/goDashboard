package collector

import (
	"context"
	"log"
	"strings"

	"golang-system-monitor/internal/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// type Container struct {
// 	ID         string `json:"Id"`
// 	Names      []string
// 	Image      string
// 	ImageID    string
// 	Command    string
// 	Created    int64
// 	Ports      []Port
// 	SizeRw     int64 `json:",omitempty"`
// 	SizeRootFs int64 `json:",omitempty"`
// 	Labels     map[string]string
// 	State      string
// 	Status     string
// 	HostConfig struct {
// 		NetworkMode string            `json:",omitempty"`
// 		Annotations map[string]string `json:",omitempty"`
// 	}
// 	NetworkSettings *SummaryNetworkSettings
// 	Mounts          []MountPoint
// }

//get containers info

type Container struct{
    Name    string  `json:"name"`
    Uptime  string  `json:"uptime"`
    Image   string  `json:"image"`  
    Status  string  `json:"status"`
}


func GetContainers() ([]Container, error){
    output := []Container{}

    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil{
        log.Println("Error creating docker client: ", err)
        return []Container{}, err
    }


    containers, err := cli.ContainerList(context.Background(), container.ListOptions{All:true})
    if err != nil{
        log.Println("Error listing containers: ", err)
        return []Container{}, err
    }

    for _, container := range containers{
        output = append(output, Container{
            Name: container.Names[0],
            Status: container.State,
            Uptime: container.Status,
            Image: container.Image,
        })
    }

    return output, nil
}
