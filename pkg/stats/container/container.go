package container

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"golang-system-monitor/internal/core"
)

type Containers []Container     

type Container struct{
    Name    string  `json:"name"`
    Status  string  `json:"status"`
    Uptime  string  `json:"uptime"`
    Image   string  `json:"image"`  
}

func (c Containers) ToPoint() []*core.Point{
    output := []*core.Point{}

    for _, container := range c{
        output = append(output, &core.Point{
            Timestamp: time.Now(),
            Measurement: "container",
            Tags: map[string]string{},
            Fields: map[string]interface{}{
                "name": container.Name,
                "status": container.Status,
                "uptime": container.Uptime,
                "image": container.Image,
            },
        })
    }
    return output
}

func ReadContainers() (Containers, error){
    output := Containers{}

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
            Status: container.Status,
            Uptime: container.State,
            Image: container.Image,
        })
    }

    return output, nil
}
