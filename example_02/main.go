package main

import "fmt"

func main() {
	m, err := NewMetricPost()
	if err != nil {
		fmt.Println(err)
	}
	err = m.Post()
	if err != nil {
		fmt.Println(err)
	}

	// devices, err := ImportDevices()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, device := range devices {
	// 	if device.GetHostname() == "switch01" {
	// 		// fmt.Println(device.GetInterfaces())
	// 		fmt.Println(device.GetLLDPNeigbours())
	// 		// device.ShowDeviceInfo()
	// 		// device.SetInterfaceDescription()
	// 	}
	// }
	// inventoryDevices := ImportInventoryDevices(devices)
	// WriteInventoryToCVS(inventoryDevices)
}
