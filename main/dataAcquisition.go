package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"runtime"
	"strconv"
	"time"
)

// dataAcquisition 数据采集，
//包含cpu数量，cpu使用情况
//内存信息，内存使用情况
//磁盘信息，磁盘使用情况
//网络信息，当前网络情况/tcp/udp连接数量
//开机时长
func dataAcquisition() string {
	dataJson := "{}"
	//cpu数量
	cpuNum := runtime.NumCPU()
	dataJson, _ = sjson.Set(dataJson, "cpu.num", cpuNum)
	//cpu使用百分比
	CpuPercent, _ := cpu.Percent(time.Second, false)
	CpuCorePercent, _ := cpu.Percent(time.Second, true)
	dataJson, _ = sjson.Set(dataJson, "cpu.percent", CpuPercent[0])
	for i := 0; i < len(CpuCorePercent); i++ {
		dataJson, _ = sjson.Set(dataJson, "cpu.core"+strconv.Itoa(i), CpuCorePercent[i])
	}
	//网络使用
	netInfo, _ := net.IOCounters(false)
	dataJson, _ = sjson.Set(dataJson, "net.bytesSent", netInfo[0].BytesSent)
	dataJson, _ = sjson.Set(dataJson, "net.bytesRecv", netInfo[0].BytesRecv)
	connsTcp, _ := net.Connections("tcp")
	connsUdp, _ := net.Connections("udp")
	dataJson, _ = sjson.Set(dataJson, "net.tcp", len(connsTcp))
	dataJson, _ = sjson.Set(dataJson, "net.udp", len(connsUdp))

	//内存使用
	memInfo, _ := mem.VirtualMemory()
	dataJson, _ = sjson.Set(dataJson, "mem.total", float64(memInfo.Total)/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "mem.free", float64(memInfo.Free)/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "mem.used", float64(memInfo.Used)/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "mem.usedPercent", memInfo.UsedPercent)
	//磁盘使用情况
	diskInfo, _ := disk.Usage("/")
	dataJson, _ = sjson.Set(dataJson, "disk.total", diskInfo.Total/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "disk.free", diskInfo.Free/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "disk.used", diskInfo.Used/1024/1024/1024)
	dataJson, _ = sjson.Set(dataJson, "disk.usage", diskInfo.UsedPercent)
	//主机信息
	hostInfo, _ := host.Info()
	dataJson, _ = sjson.Set(dataJson, "host.name", hostInfo.Hostname)
	dataJson, _ = sjson.Set(dataJson, "host.os", hostInfo.OS)
	dataJson, _ = sjson.Set(dataJson, "host.time", float64(hostInfo.Uptime)/60/60/24)

	return dataJson
}

// configInfo 配置信息读取
func configInfo() map[string]string {
	file, err := ioutil.ReadFile("serverState.json")
	if err != nil {
		fmt.Errorf("配置文件读取错误")
	}
	configJson := string(file)
	mapList := make(map[string]string)
	if gjson.Get(configJson, "mode").Value().(string) == "client" {
		mapList["serverIp"] = gjson.Get(configJson, "client.serverIp").Value().(string)
		mapList["clientKey"] = gjson.Get(configJson, "client.clientKey").Value().(string)
		mapList["mode"] = "client"
	} else {
		mapList["mode"] = "server"
		mapList["port"] = gjson.Get(configJson, "server.port").Value().(string)
		mapList["serverKey"] = gjson.Get(configJson, "server.serverKey").Value().(string)
	}
	mapList["username"] = gjson.Get(configJson, "username").Value().(string)
	return mapList
}
