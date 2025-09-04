package netbox

import "nb_import/internal/netbox/v4/device"

type NetboxImporter interface {
	WriteInventoryToCSV(devices []*device.NetboxDevice) error
	WriteInventoryToJSON(devices []*device.NetboxDevice) error
}

type NetboxDevice struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
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

func NewNetboxDevice(name string, opts ...NetboxDeviceOption) *NetboxDevice {
	dev := &NetboxDevice{
		Name: name,
	}
	for _, opt := range opts {
		opt(dev)
	}
	return dev
}

type NetboxDeviceOption func(*NetboxDevice)

func WithStatus(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Status = s
	}
}

func WithTenant(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Tenant = s
	}
}

func WithSite(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Site = s
	}
}

func WithLocation(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Location = s
	}
}

func WithRack(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Rack = s
	}
}

func WithRole(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Role = s
	}
}

func WithManufacturer(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Manufacturer = s
	}
}

func WithIPv4Address(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.IPv4Address = s
	}
}

func WithInterfaces(s []string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Interfaces = s
	}
}
