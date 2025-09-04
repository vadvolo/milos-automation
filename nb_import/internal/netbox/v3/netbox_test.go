package v3

import (
	"nb_import/internal/device"
	"reflect"
	"testing"
)

func TestImportInventoryDevices(t *testing.T) {
	tests := []struct {
		name     string
		devices  []device.AbstractDevice
		expected []*InventoryDevice
	}{
		{
			name: "all_devices_active",
			devices: []device.AbstractDevice{
				mockDevice{"host1", "Cisco", true, "192.168.0.1", []string{"eth0", "eth1"}},
				mockDevice{"host2", "Juniper", true, "192.168.0.2", []string{"eth0"}},
			},
			expected: []*InventoryDevice{
				NewInventoryDevice("host1", WithStatus("ACTIVE"), WithManufacturer("Cisco"), WithIPv4Address("192.168.0.1"), WithInterfaces([]string{"eth0", "eth1"})),
				NewInventoryDevice("host2", WithStatus("ACTIVE"), WithManufacturer("Juniper"), WithIPv4Address("192.168.0.2"), WithInterfaces([]string{"eth0"})),
			},
		},
		{
			name: "all_devices_inactive",
			devices: []device.AbstractDevice{
				mockDevice{"host3", "HP", false, "192.168.0.3", nil},
				mockDevice{"host4", "Dell", false, "192.168.0.4", nil},
			},
			expected: []*InventoryDevice{
				NewInventoryDevice("host3", WithStatus("NOTACTIVE"), WithManufacturer("HP"), WithIPv4Address("192.168.0.3")),
				NewInventoryDevice("host4", WithStatus("NOTACTIVE"), WithManufacturer("Dell"), WithIPv4Address("192.168.0.4")),
			},
		},
		{
			name: "mixed_status_devices",
			devices: []device.AbstractDevice{
				mockDevice{"host5", "Cisco", true, "192.168.1.1", []string{"eth0"}},
				mockDevice{"host6", "Juniper", false, "192.168.1.2", nil},
			},
			expected: []*InventoryDevice{
				NewInventoryDevice("host5", WithStatus("ACTIVE"), WithManufacturer("Cisco"), WithIPv4Address("192.168.1.1"), WithInterfaces([]string{"eth0"})),
				NewInventoryDevice("host6", WithStatus("NOTACTIVE"), WithManufacturer("Juniper"), WithIPv4Address("192.168.1.2")),
			},
		},
		{
			name:     "no_devices",
			devices:  []device.AbstractDevice{},
			expected: []*InventoryDevice{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ImportInventoryDevices(tt.devices)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ImportInventoryDevices() = %v, want %v", got, tt.expected)
			}
		})
	}
}

type mockDevice struct {
	hostname   string
	vendor     string
	status     bool
	address    string
	interfaces []string
}

func (m mockDevice) GetHostname() string {
	return m.hostname
}

func (m mockDevice) GetVendor() string {
	return m.vendor
}

func (m mockDevice) GetStatus() bool {
	return m.status
}

func (m mockDevice) GetAddress() string {
	return m.address
}

func (m mockDevice) ShowInterfaces() []string {
	return m.interfaces
}
