package main

import (
	"fmt"
	"time"
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
		time.Sleep(10 * time.Second)
		device.SetInterfaceDescription()
	}
	// host := "10.0.10.2"
	// password := "123qwe"
	// login := "vdm"
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// logger := zap.Must(zap.NewDevelopmentConfig().Build())

	// creds := dcreds.NewSimpleCredentials(
	// 	dcreds.WithUsername(login),
	// 	dcreds.WithPassword(dcreds.Secret(password)), // and password
	// 	dcreds.WithLogger(logger),
	// )
	// connector := ssh.NewStreamer(host, creds, ssh.WithLogger(logger))
	// dev := cisco.NewDevice(connector) // huawei CLI upon SSH
	// err := dev.Connect(ctx)           // connection happens here
	// if err != nil {
	// 	panic(err)
	// }
	// defer dev.Close()
	// reses, _ := dev.ExecuteBulk(cmd.NewCmdList([]string{"show lldp neighbors"}))
	// for _, res := range reses {
	// 	if res.Status() == 0 {
	// 		fmt.Printf("Result: %s\n", res.Output())
	// 	} else {
	// 		fmt.Printf("Error: %s\nStatus: %d\n", res.Status(), res.Error())
	// 	}
	// }
}
