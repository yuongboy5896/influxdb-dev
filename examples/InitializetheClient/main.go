package main

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "zwvS0JXTQU2LUiEnWCmLjr6mq_E1UPJagrpePLalFO-SvsmVxKFoC-f1oDZDTU_PTuIGKiVuseFQIn2OR9YFvw=="
	const bucket = "devops"
	const org = "devops"

	client := influxdb2.NewClient("http://192.168.2.60:8086", token)
	// always close client at the en
	defer client.Close()
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -5h) |> filter(fn: (r) => r._measurement == \"stat\") |> filter(fn: (r) => r[\"_measurement\"] == \"stat\")  |> filter(fn: (r) => r[\"_field\"] == \"cpu\")    |> movingAverage(n: 5)", bucket)
	// Get query client
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("1 value: %v\n", result.Record().Value())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"cpu": 30.0, "max": 45.0},
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	writeAPI.Flush()
	// create point using fluent style
	//	p = influxdb2.NewPointWithMeasurement("stat").
	//		AddTag("unit", "temperature").
	//		AddField("cpu", 23.2).
	//		AddField("max", 45).
	//		SetTime(time.Now())
	// write point asynchronously
	///	writeAPI.WritePoint(p)
	// Flush writes
	//	writeAPI.Flush()

}
