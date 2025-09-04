package v3

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strings"

	"nb_import/internal/netbox"
)

const InventoryFileName = "exportInventory.csv"

type NetboxDeviceV3 struct{}

func (nb *NetboxDeviceV3) WriteInventoryToCSV(devices []*netbox.NetboxDevice) error {
	csvFile, err := os.Create(InventoryFileName)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer func() { _ = csvFile.Close() }()

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
