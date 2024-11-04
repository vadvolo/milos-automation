package main

import (
	"encoding/csv"
	"log"
	"os"
)

var InterfaceHeader = []string{
	"device",
	"name",
	"type",
}

func ExportIfaceNB(devices []AbstractDevice) {
	var data [][]string
	csvFile, err := os.Create("interfaces.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ','
	defer csvFile.Close()
	data = append(data, InterfaceHeader)
	for _, d := range devices {
		for _, iface := range d._Interfaces() {
			row := []string{d._Hostname() + ".nh.com", iface.Name, iface.IfaceType}
			data = append(data, row)
		}
	}
	csvwriter.WriteAll(data)
}

func ExportIPsNB(devices []AbstractDevice) {

	var header = []string{
		"device",
		"interface",
		"address",
		"status",
	}

	var data [][]string
	csvFile, err := os.Create("ipaddresses.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ','
	defer csvFile.Close()
	data = append(data, header)
	for _, d := range devices {
		for _, iface := range d._Interfaces() {
			row := []string{d._Hostname() + ".nh.com", iface.Name, iface.IPAddress, "active"}
			data = append(data, row)
		}
	}
	csvwriter.WriteAll(data)
}
