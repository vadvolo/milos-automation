package main

import (
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"strconv"
	// "strings"
)

type InventoryDevice struct {
	Name           string   `json:"name"`
	Role           string   `json:"role"`
	Tenant         string   `json:"tenant"`
	Manufacturer   string   `json:"manufacturer"`
	Type           string   `json:"device_type"`
	Platform       string   `json:"platform"`
	SerialNumber   string   `json:"serial"`
	AssetTag       string   `json:"asset_tag"`
	Status         string   `json:"status"`
	Site           string   `json:"site"`
	Location       string   `json:"location"`
	Rack           string   `json:"rack"`
	Position       string   `json:"position"`
	Face           string   `json:"face"`
	Latitude       string   `json:"latitude"`
	Longitude      string   `json:"longitude"`
	Parent         string   `json:"parent"`
	Bay            string   `json:"device_bay"`
	Airflow        string   `json:"airflow"`
	VirtualChassis string   `json:"virtual_chassis"`
	VCPosition     string   `json:"vc_position"`
	VCPriority     string   `json:"vc_priority"`
	Cluster        string   `json:"cluster"`
	Description    string   `json:"description"`
	ConfigTemplate string   `json:"config_template"`
	Comments       string   `json:"comments"`
	Tags           []string `json:"tags"`
	ID             string   `json:"id"`
	DevLogin       string   `json:"cf_dev_login"`
	DevPassword    string   `json:"cf_dev_password"`
}

var HEADER = []string{
	"name",
	"role",
	"tenant",
	"manufacturer",
	"device_type",
	"platform",
	"serial",
	"asset_tag",
	"status",
	"site",
	"location",
	"rack",
	"position",
	"face",
	"latitude",
	"longitude",
	"parent",
	"device_bay",
	"airflow",
	"virtual_chassis",
	"vc_position",
	"vc_priority",
	"cluster",
	"description",
	"config_template",
	"comments",
	"tags",
	"id",
	"cf_dev_login",
	"cf_dev_password",
}

type NetboxDeviceType struct {
	Manufacturer    string `json:"manufacturer"`
	DefaultPlatform string `json:"default_platform"`
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
				WithStatus("Active"),
				WithSite("lab"),
				WithRole("Router"),
				WithDeviceType(device._Model()),
				WithManufacturer(device._Vendor()),
				WithDevLogin("milos"),
				WithDevPassword("milos"),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		} else {
			inventoryDevice := NewInventoryDevice(
				device._Hostname(),
				WithStatus("Inactive"),
				WithManufacturer(device._Vendor()),
			)
			inventoryDevices = append(inventoryDevices, inventoryDevice)
		}
	}
	return inventoryDevices
}

func WriteInventoryToCSV(devices []*InventoryDevice) error {
	var data [][]string
	csvFile, err := os.Create("devices.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ','
	defer csvFile.Close()
	data = append(data, HEADER)
	for _, device := range devices {
		s := reflect.ValueOf(device).Elem()
		row := []string{}
		for i := 0; i < s.NumField(); i++ {
			if s.Field(i).Kind() == reflect.Slice {
				ret := s.Field(i).Len()
				row = append(row, strconv.Itoa(ret))
				// val := s.Field(i)
				// ret := new(strings.Builder)
				// delim := ";"
				// for i := 0; i < val.Len(); i++ {
				// 	if val.Index(i).Kind() == reflect.String {
				// 		ret.WriteString(val.Index(i).String())
				// 		ret.WriteString(delim)
				// 	}
				// }
				// row = append(row, ret.String())
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

func WithDeviceType(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Type = s
	}
}

func WithManufacturer(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.Manufacturer = s
	}
}

func WithDevLogin(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.DevLogin = s
	}
}

func WithDevPassword(s string) InventoryDeviceOption {
	return func(d *InventoryDevice) {
		d.DevPassword = s
	}
}
