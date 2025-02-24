package collector

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Container struct{
    Name    string  `json:"name"`
    Status  string  `json:"status"`
    Uptime  string  `json:"uptime"`
    Image   string  `json:"image"`  
}

func (c *Container) ToMap() map[string]interface{}{
    return map[string]interface{}{
        "name": c.Name,
        "status": c.Status,
        "uptime": c.Uptime,
        "image": c.Image,
    }
}


func ReadContainers() ([]Container, error){
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
            Status: container.Status,
            Uptime: container.State,
            Image: container.Image,
        })
    }

    return output, nil
}
