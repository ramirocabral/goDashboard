package influxdb

import (
	"context"
	"fmt"
	"time"

	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/storage"
	"golang-system-monitor/internal/utils"

	// "golang-system-monitor/internal/logger"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct{
    client	influxdb2.Client 
    writeAPI	api.WriteAPI
    queryAPI	api.QueryAPI
    org		string
    bucket	string
}

func New(addr, token, org, bucket string) (storage.Storage, error) {
    client := influxdb2.NewClient(addr, token)
    writeAPI := client.WriteAPI(org,bucket)
    queryAPI := client.QueryAPI(org)

    influx := &InfluxStore{
        client: client,
        writeAPI: writeAPI,
        queryAPI: queryAPI,
	bucket: bucket,
    }

    alive, err := client.Ping(context.Background())
    _ = alive

    if err != nil {
	return nil, err
    }

    return influx, nil
}

func (s *InfluxStore) Close(){
    s.client.Close()
}

func (s *InfluxStore) ID() string{
    return s.client.ServerURL()
}

func (s *InfluxStore) WriteStats(points []*core.Point){
    for _, point := range points{
        p := influxdb2.NewPoint(
            point.Measurement,
            point.Tags,
            point.Fields,
            point.Timestamp)
        s.writeAPI.WritePoint(p)
    }
}

func (s *InfluxStore) ReadCpuStats(startTime, endTime time.Time, interval string) (storage.CPUResponse, error){

    query := fmt.Sprintf(`
	from(bucket: "%s")
	    |> range(start: %s, stop: %s)
	    |> filter(fn: (r) => r._measurement == "cpu")
	    |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`,
	s.bucket,
	startTime.UTC().Format(time.RFC3339),
	endTime.UTC().Format(time.RFC3339),
	interval,
    )

    result, err := s.queryAPI.Query(context.Background(), query)

    if err != nil {
        return storage.CPUResponse{}, err
    }

    stats := parseCpuStats(result)

    return stats, nil
}

func parseCpuStats(result *api.QueryTableResult) storage.CPUResponse{

    stats := storage.CPUResponse{}

    data := make([]storage.CPUPoint, 0)

    for result.Next() {
	record := result.Record()
	stats.ModelName = record.ValueByKey("model_name").(string)
	stats.Frequency = utils.StrToUint64(record.ValueByKey("frequency").(string))
	stats.Family = record.ValueByKey("family").(string)
	stats.Cores = utils.StrToUint64(record.ValueByKey("cores").(string))
	stats.Threads = utils.StrToUint64(record.ValueByKey("threads").(string))

	timestamp := record.Time()
	usagePercentage := utils.RoundFloat(record.ValueByKey("usage_percentage").(float64),2)
	idlePercentage := utils.RoundFloat(record.ValueByKey("idle_percentage").(float64),2)

	point := storage.CPUPoint{
	    Timestamp: timestamp,
	    UsagePercentage: usagePercentage,
	    IdlePercentage: idlePercentage,
	}

	data = append(data, point)
    }

    stats.Data = data

    return stats
}

func (s *InfluxStore) ReadIOStats(startTime, endTime time.Time, interval string) (storage.IOResponse, error){

    query := fmt.Sprintf(`
	from(bucket: "%s")
	    |> range(start: %s, stop: %s)
	    |> filter(fn: (r) => r._measurement == "io")
	    |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`,
	s.bucket,
	startTime.UTC().Format(time.RFC3339),
	endTime.UTC().Format(time.RFC3339),
	interval,
    )
    
        result, err := s.queryAPI.Query(context.Background(), query)
        if err != nil {
            return storage.IOResponse{}, err
        }
    
        stats := parseIOStats(result)
    
        return stats, nil
}

func parseIOStats(result *api.QueryTableResult) storage.IOResponse{
	deviceMap := make(map[string][]storage.IOPoint)

	for result.Next() {
		record := result.Record()
		device := record.ValueByKey("device").(string)
		timestamp := record.Time()
		readBytes := uint64(record.ValueByKey("read_bytes").(float64))
		writeBytes := uint64(record.ValueByKey("write_bytes").(float64))

		point := storage.IOPoint{
			Timestamp:        timestamp,
			KBReadPerSecond:  readBytes,
			KBWritePerSecond: writeBytes,
		}

		deviceMap[device] = append(deviceMap[device], point)
	}

	var response storage.IOResponse
	for device, stats := range deviceMap {
		response.Devices = append(response.Devices, storage.IOStats{
			Device: device,
			Stats:  stats,
		})
	}
	return response
}

func (s *InfluxStore) ReadMemoryStats(startTime, endTime time.Time, interval string) (storage.MemoryResponse, error){
    query := fmt.Sprintf(`
	from(bucket: "%s")
	    |> range(start: %s, stop: %s)
	    |> filter(fn: (r) => r._measurement == "memory")
	    |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`,
	s.bucket,
	startTime.UTC().Format(time.RFC3339),
	endTime.UTC().Format(time.RFC3339),
	interval,
    )

    result, err := s.queryAPI.Query(context.Background(), query)

    if err != nil {
        return storage.MemoryResponse{}, err
    }

    stats := parseMemoryStats(result)

    return stats, nil
}

func parseMemoryStats(result *api.QueryTableResult) storage.MemoryResponse{
    var response storage.MemoryResponse

    data := make([]storage.MemoryPoint, 0)

    for result.Next() {
	record := result.Record()
	response.Type = record.ValueByKey("type").(string)
	response.Frequency = utils.StrToUint64(record.ValueByKey("frequency").(string))
	
	timestamp := record.Time()
	usedPercentage := utils.RoundFloat(record.ValueByKey("used_percentage").(float64),2)
	total := uint64(record.ValueByKey("total").(float64))
	used := uint64(record.ValueByKey("used").(float64))
	free := uint64(record.ValueByKey("free").(float64))
	active := uint64(record.ValueByKey("active").(float64))
	inactive := uint64(record.ValueByKey("inactive").(float64))
	buffers := uint64(record.ValueByKey("buffers").(float64))
	cached := uint64(record.ValueByKey("cached").(float64))

	point := storage.MemoryPoint{
	    Timestamp: timestamp,
	    UsedPercentage: usedPercentage,
	    Total: total,
	    Used: used,
	    Free: free,
	    Active: active,
	    Inactive: inactive,
	    Buffers: buffers,
	    Cached: cached,
	}
    
	data = append(data, point)

    }
    
    response.Data = data

    return response
}

func (s *InfluxStore) ReadNetworkStats(startTime, endTime time.Time, interval string) (storage.NetworkResponse, error){

    query := fmt.Sprintf(`
	from(bucket: "%s")
	    |> range(start: %s, stop: %s)
	    |> filter(fn: (r) => r._measurement == "network")
	    |> aggregateWindow(every: %s, fn: mean, createEmpty: false)
	    |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	`,
	s.bucket,
	startTime.UTC().Format(time.RFC3339),
	endTime.UTC().Format(time.RFC3339),
	interval,
    )

    result, err := s.queryAPI.Query(context.Background(), query)

    if err != nil {
        return storage.NetworkResponse{}, err
    }

    stats := parseNetworkStats(result)

    return stats, nil
}

func parseNetworkStats(result *api.QueryTableResult) storage.NetworkResponse{

	interfaceMap := make(map[string]*storage.NetworkStats)

	for result.Next() {
		record := result.Record()
		timestamp := record.Time()
		interfaceName := record.ValueByKey("interface").(string)
		rxBytesPS := uint64(record.ValueByKey("rx_bytes_ps").(float64))
		txBytesPS := uint64(record.ValueByKey("tx_bytes_ps").(float64))

		point := storage.NetworkPoint{
			Timestamp: timestamp,
			RxBytesPS: rxBytesPS,
			TxBytesPS: txBytesPS,
		}

		key := interfaceName

		if _, exists := interfaceMap[key]; !exists {
			interfaceMap[key] = &storage.NetworkStats{
				Interface: interfaceName,
				Data:      []storage.NetworkPoint{},
			}
		}

		interfaceMap[key].Data = append(interfaceMap[key].Data, point)
	}

	var response storage.NetworkResponse
	for _, stats := range interfaceMap {
		response.Interfaces = append(response.Interfaces, *stats)
	}

	return response
}
