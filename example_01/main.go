package main

import (
	"fmt"
)

func main() {
	devices, err := ImportDevices()
	if err != nil {
		panic(err)
	}
	for _, device := range devices {
		fmt.Println(device._Hostname())
		if device.Ping() == nil {
			device.SetStatus(true)
			fmt.Println(device.GetInterfaces())
			fmt.Println(device.GetLLDPNeigbours())
			device.ShowDeviceInfo()
			device.SetInterfaceDescription()
		} else {
			device.SetStatus(false)
		}
	}
	inventoryDevices := ImportInventoryDevices(devices)
	WriteInventoryToCSV(inventoryDevices)
}
