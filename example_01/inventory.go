package main

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"strings"
)

type InventoryDevice struct {
	Name               string `json:"name"`
	Status             string
	Tenant             string
	Site               string
	Location           string
	Rack               string
	Role               string
	Manufacturer       string
	Type               string
	IPAddress          string
	TenantGroup        string
	ID                 string
	Platform           string
	SerialNumber       string
	AssetTag           string
	Region             string
	SiteGroup          string
	ParentDevice       string
	Position           string
	RackFace           string
	Latitude           string
	Longitude          string
	Airflow            string
	IPv4Address        string
	IPv6Address        string
	OOBIP              string
	Cluster            string
	VirtualChassis     string
	VCPosition         string
	VCPriority         string
	Description        string
	ConfigTemplate     string
	Contacts           string
	Tags               []string
	Created            string
	LastUpdated        string
	ConsolePorts       string
	ConsoleServerPorts string
	PowerPorts         string
	PowerOutlets       string
	Interfaces         []string
	FrontPorts         string
	RearPorts          string
	DeviceBays         string
	ModuleBays         string
	InventoryItems     []string
}

func NewInventoryDevice(name string, opts ...InventoryDeviceOption) *InventoryDevice {
	device := &InventoryDevice{
		Name: name,
	}
	for _, opt := range opts {
		opt(device)
	}
	return device
}

func ImportInventoryDevices(devices []AbstractDevice) []*InventoryDevice {
	var inventoryDevices []*InventoryDevice
	for _, device := range devices {
		if device.GetStatus() {
			inventoryDevice := NewInventoryDevice(
				device._Hostname(),
				WithStatus("ACTIVE"),
				WithManufacturer(device._Vendor()),
				WithIPv4Address(device._Address()),
				WithInterfaces(device.ShowInterfaces()),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		} else {
			inventoryDevice := NewInventoryDevice(
				device._Hostname(),
				WithStatus("NOTACTIVE"),
				WithIPv4Address(device._Address()),
				WithManufacturer(device._Vendor()),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		}
	}
	return inventoryDevices
}

func WriteInventoryToCSV(devices []*InventoryDevice) error {
	var data [][]string
	csvFile, err := os.Create("exportInventory.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ','
	defer csvFile.Close()
	for _, device := range devices {
		s := reflect.ValueOf(device).Elem()
		row := []string{}
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).Kind() == reflect.Slice {
				val := s.Field(i)
				ret := new(strings.Builder)
				delim := ";"
				for i := 0; i < val.Len(); i++ {
					if val.Index(i).Kind() == reflect.String {
						// fmt.Println(val.Index(i).String())
						ret.WriteString(val.Index(i).String())
						ret.WriteString(delim)
					}
				}
				row = append(row, ret.String())
			} else {
				row = append(row, s.Field(i).String())
			}
		}
		data = append(data, row)
	}
	csvwriter.WriteAll(data)
	return nil
}

type InventoryDeviceOption func(*InventoryDevice)

func WithStatus(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Status = s
	}
}

func WithTenant(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Tenant = s
	}
}

func WithSite(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Site = s
	}
}

func WithLocation(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Location = s
	}
}

func WithRack(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Rack = s
	}
}

func WithRole(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Role = s
	}
}

func WithManufacturer(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Manufacturer = s
	}
}

func WithIPv4Address(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.IPv4Address = s
	}
}

func WithInterfaces(s []string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Interfaces = s
	}
}
