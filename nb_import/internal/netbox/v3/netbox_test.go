package v3

import (
	"reflect"
	"testing"

	"nb_import/internal/device"
	"nb_import/internal/device/mocks"
)

func TestImportInventoryDevices(t *testing.T) {
	// Create and configure the first mock
	mockDev1 := new(mocks.AbstractDevice)
	mockDev1.
		On("GetHostname").Return("host1").
		On("GetVendor").Return("Cisco").
		On("GetStatus").Return(device.StatusActive).
		On("GetAddress").Return("192.168.0.1").
		On("ShowInterfaces").Return([]string{"eth0", "eth1"})

	// Create and configure second mock with different values
	mockDev2 := new(mocks.AbstractDevice)
	mockDev2.
		On("GetHostname").Return("host2").
		On("GetVendor").Return("Cisco").
		On("GetStatus").Return(device.StatusActive).
		On("GetAddress").Return("192.168.0.2").
		On("ShowInterfaces").Return([]string{"eth0"})

	// Create and configure inactive mocks
	mockDev3 := new(mocks.AbstractDevice)
	mockDev3.
		On("GetHostname").Return("host3").
		On("GetVendor").Return("HP").
		On("GetStatus").Return("NOTACTIVE").
		On("GetAddress").Return("192.168.0.3")

	mockDev4 := new(mocks.AbstractDevice)
	mockDev4.
		On("GetHostname").Return("host4").
		On("GetVendor").Return("Dell").
		On("GetStatus").Return("NOTACTIVE").
		On("GetAddress").Return("192.168.0.4")

	tests := []struct {
		name     string
		devices  []device.AbstractDevice
		expected []*InventoryDevice
	}{
		{
			name: "all_devices_active",
			devices: []device.AbstractDevice{
				mockDev1,
				mockDev2,
			},
			expected: []*InventoryDevice{
				NewInventoryDevice("host1", WithStatus("ACTIVE"), WithManufacturer("Cisco"), WithIPv4Address("192.168.0.1"), WithInterfaces([]string{"eth0", "eth1"})),
				NewInventoryDevice("host2", WithStatus("ACTIVE"), WithManufacturer("Cisco"), WithIPv4Address("192.168.0.2"), WithInterfaces([]string{"eth0"})),
			},
		},
		{
			name: "all_devices_inactive",
			devices: []device.AbstractDevice{
				mockDev3,
				mockDev4,
			},
			expected: []*InventoryDevice{
				NewInventoryDevice("host3", WithStatus("NOTACTIVE"), WithManufacturer("HP"), WithIPv4Address("192.168.0.3")),
				NewInventoryDevice("host4", WithStatus("NOTACTIVE"), WithManufacturer("Dell"), WithIPv4Address("192.168.0.4")),
			},
		},
		{
			name:     "no_devices",
			devices:  []device.AbstractDevice{},
			expected: nil,
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

	// Verify mock expectations
	mockDev1.AssertExpectations(t)
	mockDev2.AssertExpectations(t)
	mockDev3.AssertExpectations(t)
	mockDev4.AssertExpectations(t)
}
