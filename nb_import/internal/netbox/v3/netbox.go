package v3

import (
	"encoding/csv"
	"fmt"
	"nb_import/internal/device"
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

func ImportInventoryDevices(devices []device.AbstractDevice) []*InventoryDevice {
	var inventoryDevices []*InventoryDevice
	for _, d := range devices {
		fmt.Println(d.GetHostname())
		if d.GetStatus() {
			inventoryDevice := NewInventoryDevice(
				d.GetHostname(),
				WithStatus("ACTIVE"),
				WithManufacturer(d.GetVendor()),
				WithIPv4Address(d.GetAddress()),
				WithInterfaces(d.ShowInterfaces()),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		} else {
			inventoryDevice := NewInventoryDevice(
				d.GetHostname(),
				WithStatus("NOTACTIVE"),
				WithIPv4Address(d.GetAddress()),
				WithManufacturer(d.GetVendor()),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		}
	}
	return inventoryDevices
}

const InventoryFileName = "exportInventory.csv"

func WriteInventoryToCSV(devices []*InventoryDevice) error {
	csvFile, err := os.Create(InventoryFileName)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	for _, dev := range devices {
		v := reflect.Indirect(reflect.ValueOf(dev))
		row := make([]string, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			switch field.Kind() {
			case reflect.Slice:
				var parts []string
				for j := 0; j < field.Len(); j++ {
					parts = append(parts, fmt.Sprint(field.Index(j)))
				}
				row[i] = strings.Join(parts, ";")
			default:
				row[i] = fmt.Sprint(field.Interface())
			}
		}

		if err := csvwriter.Write(row); err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}
	}
	return csvwriter.Error()
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
