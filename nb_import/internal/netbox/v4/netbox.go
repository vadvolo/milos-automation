package v4

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"nb_import/internal/netbox/v4/device"
	"os"
	"reflect"
	"strings"
)

const InventoryFileName = "exportInventory.csv"

type NetboxDevice struct{}

func (nb *NetboxDevice) WriteInventoryToCSV(devices []*device.NetboxDevice) error {
	csvFile, err := os.Create(InventoryFileName)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer func() { _ = csvFile.Close() }()

	csvwriter := csv.NewWriter(csvFile)
	defer csvwriter.Flush()

	if len(devices) == 0 {
		return fmt.Errorf("no devices to write")
	}

	devType := reflect.TypeOf(*devices[0])
	headers := make([]string, devType.NumField())
	for i := 0; i < devType.NumField(); i++ {
		tag := devType.Field(i).Tag.Get("json")
		if tag == "" {
			headers[i] = devType.Field(i).Name
		} else {
			headers[i] = strings.Split(tag, ",")[0] // отрезаем ",omitempty" и т.п.
		}
	}

	if err := csvwriter.Write(headers); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

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

func (nb *NetboxDevice) WriteInventoryToJSON(devices []*device.NetboxDevice) error {
	if len(devices) == 0 {
		return fmt.Errorf("no devices to write")
	}

	jsonFile, err := os.Create(InventoryFileName)
	if err != nil {
		return fmt.Errorf("failed creating file: %w", err)
	}
	defer func() { _ = jsonFile.Close() }()

	// Преобразуем devices в JSON (pretty)
	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(devices); err != nil {
		return fmt.Errorf("error encoding to json: %w", err)
	}
	return nil
}
