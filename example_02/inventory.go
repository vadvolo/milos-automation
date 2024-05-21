package main

import (
	"bytes"

	"github.com/mohae/struct2csv"
)

type InventoryDevice struct {
	Name               string
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
		inventoryDevice := NewInventoryDevice(
			device.GetHostname(),
			WithStatus("OPERATIONAL"),
		)
		inventoryDevices = append(inventoryDevices, inventoryDevice)
	}
	return inventoryDevices
}

func WriteInventoryToCVS(devices []*InventoryDevice) error {
	buff := &bytes.Buffer{}
	w := struct2csv.NewWriter(buff)
	err := w.WriteStructs(devices)
	return err
}

type InventoryDeviceOption func(*InventoryDevice)

func WithStatus(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Status = s
	}
}
