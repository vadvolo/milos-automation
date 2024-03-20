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
		if device.GetHostname() == "switch01" {
			// fmt.Println(device.GetInterfaces())
			fmt.Println(device.GetLLDPNeigbours())
			// device.ShowDeviceInfo()
			// device.SetInterfaceDescription()
		}
	}
	inventoryDevices := ImportInventoryDevices(devices)
	WriteInventoryToCVS(inventoryDevices)
}
