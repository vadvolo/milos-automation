package cvs

import (
	"fmt"
	"nb_import/internal/device"
	"nb_import/internal/netbox"
	nbdev "nb_import/internal/netbox/v4/device"
)

const (
	StatusActive   = "active"
	StatusDeactive = "deactive"
)

func ImportNetboxDevices(devices []device.AbstractDevice, importer netbox.NetboxImporter) error {
	var NetboxDevices []*nbdev.NetboxDevice
	for _, d := range devices {
		fmt.Println(d.GetHostname())
		if d.GetStatus() {
			NetboxDevice := nbdev.NewNetboxDevice(
				d.GetHostname(),
				nbdev.WithStatus(StatusActive),
				nbdev.WithManufacturer(d.GetVendor()),
				nbdev.WithRole("Router"),
				nbdev.WithDeviceType("7200"),
				nbdev.WithSite("lab"),
			)
			NetboxDevices = append(NetboxDevices, NetboxDevice)
		} else {
			NetboxDevice := nbdev.NewNetboxDevice(
				d.GetHostname(),
				nbdev.WithStatus("NOTACTIVE"),
			)
			NetboxDevices = append(NetboxDevices, NetboxDevice)
		}
	}
	err := importer.WriteInventoryToCSV(NetboxDevices)
	if err != nil {
		return err
	}

	err = importer.WriteInventoryToJSON(NetboxDevices)
	if err != nil {
		return err
	}
	return nil
}
