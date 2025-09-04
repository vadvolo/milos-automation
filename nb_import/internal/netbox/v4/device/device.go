package device

type NetboxDevice struct {
	Role           string   `json:"role"`
	Manufacturer   string   `json:"manufacturer"`
	DeviceType     string   `json:"device_type"`
	Status         string   `json:"status"`
	Site           string   `json:"site"`
	Name           string   `json:"name,omitempty"`
	Tenant         string   `json:"tenant,omitempty"`
	Platform       string   `json:"platform,omitempty"`
	Serial         string   `json:"serial,omitempty"`
	AssetTag       string   `json:"asset_tag,omitempty"`
	Location       string   `json:"location,omitempty"`
	Rack           string   `json:"rack,omitempty"`
	Position       string   `json:"position,omitempty"`
	Face           string   `json:"face,omitempty"`
	Latitude       string   `json:"latitude,omitempty"`
	Longitude      string   `json:"longitude,omitempty"`
	Parent         string   `json:"parent,omitempty"`
	DeviceBay      string   `json:"device_bay,omitempty"`
	Airflow        string   `json:"airflow,omitempty"`
	VirtualChassis string   `json:"virtual_chassis,omitempty"`
	VCPosition     string   `json:"vc_position,omitempty"`
	VCPriority     string   `json:"vc_priority,omitempty"`
	Cluster        string   `json:"cluster,omitempty"`
	Description    string   `json:"description,omitempty"`
	ConfigTemplate string   `json:"config_template,omitempty"`
	Comments       string   `json:"comments,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	ID             string   `json:"id,omitempty"`
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

func WithDeviceType(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.DeviceType = s
	}
}

func WithSite(s string) NetboxDeviceOption {
	return func(d *NetboxDevice) {
		d.Site = s
	}
}
