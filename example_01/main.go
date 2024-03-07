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
		fmt.Println(device.GetInterfaces())
		fmt.Println(device.GetLLDPNeigbours())
		device.ShowDeviceInfo()
		device.SetInterfaceDescription()
	}
}
