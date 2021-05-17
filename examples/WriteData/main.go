package main

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func main() {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "zwvS0JXTQU2LUiEnWCmLjr6mq_E1UPJagrpePLalFO-SvsmVxKFoC-f1oDZDTU_PTuIGKiVuseFQIn2OR9YFvw=="
	const bucket = "devops"
	const org = "devops"

	client := influxdb2.NewClient("http://192.168.2.60:8086", token)
	// always close client at the end
	defer client.Close()
	///
	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)
	//
	info, err := host.Info()
	if err == nil {
		fmt.Println(info.Hostname)
	} else {
		panic(err)
	}

	// create point using fluent style
	p := influxdb2.NewPointWithMeasurement("computer").
		//设备iD
		AddTag("devid", info.Hostname).
		//cpu
		AddField("cpu", GetCpuPercent()).
		//Men
		AddField("Mem", GetMemPercent()).
		//disk
		AddField("disk", GetDiskPercent()).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

}
