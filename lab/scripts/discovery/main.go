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
			device.GetSoftware()
			fmt.Println(device.GetInterfaces())
			fmt.Println(device.GetNeigbours())
			device.ShowDeviceInfo()
			device.SetInterfaceDescription()
		} else {
			device.SetStatus(false)
		}
	}
	ExportIfaceNB(devices)
	ExportIPsNB(devices)
	inventoryDevices := ImportInventoryDevices(devices)
	WriteInventoryToCSV(inventoryDevices)

}
