package main

type NBDevice struct {
	Name         string          `json:"name"`
	Role         *NBDeviceRole   `json:"role"`
	Vendor       *NBManufacturer `json:"manufacturer"`
	DeviceType   *NBDeviceType   `json:"device_type"`
	Platform     string          `json:"platform"`
	SerialNumber string          `json:"serial"`
	Site         string          `json:"site"`
	Location     string          `json:"location"`
}

type NBDeviceRole struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}

type NBManufacturer struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type NBDeviceType struct {
	Vendor *NBManufacturer `json:"manufacturer"`
	Model  string          `json:"model"`
	Slug   string          `json:"slug"`
	Units  string          `json:"u_height"`
}

type NBInterface struct {
	Device      *NBDevice `json:"device"`
	Name        string    `json:"name"`
	Type        string    `json:"device_type"`
	Speed       string    `json:"speed"`
	Duplex      string    `json:"duplex"`
	Mac         string    `json:"mac"`
	MTU         string    `json:"mtu"`
	Description string    `json:"description"`
}
