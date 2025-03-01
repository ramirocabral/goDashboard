package container

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Containers []Container

type Container struct{
    Name    string  `json:"name"`
    Status  string  `json:"status"`
    Uptime  string  `json:"uptime"`
    Image   string  `json:"image"`  
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
