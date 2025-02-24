package db

import (

    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

func New(addr, token string) (influxdb2.Client, error) {
    client := influxdb2.NewClient(addr, token)

    return client, nil
}

func NewAPI(client influxdb2.Client, org, bucket string) api.WriteAPI {
    return client.WriteAPI(org, bucket)
}
