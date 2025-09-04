package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"nb_import/internal/device"
	"nb_import/internal/device/cisco"
	"nb_import/internal/device/huawei"
	nbcvs "nb_import/internal/netbox/cvs"
	nb "nb_import/internal/netbox/v4"
)

func main() {
	devices, err := ImportDevices()

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, dev := range devices {
		wg.Add(1)
		go func() {
			fmt.Println(dev.GetHostname())
			err := dev.Ping()
			if err != nil {
				dev.SetStatus(false)
			} else {
				dev.SetStatus(true)
			}

			println(dev.GetHostname(), dev.GetStatus())

			if dev.GetStatus() {
				dev.SetStatus(true)
				fmt.Println(dev.GetInterfaces())
				fmt.Println(dev.GetLLDPNeigbours())
				dev.ShowDeviceInfo()
				err := dev.SetInterfaceDescription()
				if err != nil {
					log.Println(err)
				}
			}
			defer wg.Done()
		}()

	}
	wg.Wait()

	importer := &nb.NetboxDevice{}

	err = nbcvs.ImportNetboxDevices(devices, importer)
	if err != nil {
		log.Fatal(err)
	}
}

func ImportDevices() ([]device.AbstractDevice, error) {
	file, _ := os.ReadFile("inventory.json")
	var devices []*device.Device
	err := json.Unmarshal(file, &devices)
	if err != nil {
		return nil, err
	}

	var ret []device.AbstractDevice
	for _, dev := range devices {
		if dev.Vendor == "Cisco" {
			ret = append(ret, &cisco.CiscoDevice{
				Device: dev,
			})
		}
		if dev.Vendor == "Huawei" {
			ret = append(ret, &huawei.VRPDevice{
				Device: dev,
			})
		}
	}

	return ret, nil
}
