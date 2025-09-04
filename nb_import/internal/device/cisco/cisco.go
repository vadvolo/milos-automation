package cisco

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"nb_import/internal/device"
)

type CiscoDevice struct {
	*device.Device
}

func (d *CiscoDevice) CutIfaceName(name string) string {
	if strings.Contains(name, "FastEthernet") {
		r := regexp.MustCompile(`FastEthernet`)
		return r.ReplaceAllString(name, "Fa")
	}
	if strings.Contains(name, "GigabitEthernet") {
		r := regexp.MustCompile(`GigabitEthernet`)
		return r.ReplaceAllString(name, "Gi")
	}
	return ""
}

func (d *CiscoDevice) GetVersion() error {
	_, err := d.Device.SendCommand("show ver")
	if err != nil {
		return err
	}
	return nil
}

func (d *CiscoDevice) GetInterfaces() error {
	data, err := d.Device.SendCommand("show ip interface brief")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data.Output()))
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	startLine := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Interface") {
			startLine = i + 1
			break
		}
	}
	if startLine == 0 {
		return nil
	}
	for i := startLine; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		splitLine := strings.Split(line, " ")
		d.Interfaces = append(d.Interfaces, &device.Interface{
			Name:      splitLine[0],
			ShortName: d.CutIfaceName(splitLine[0]),
		})
	}
	return nil
}

func (d *CiscoDevice) GetLLDPNeigbours() error {
	data, err := d.Device.SendCommand("show lldp neighbors")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data.Output()))
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	startLine := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Device ID") {
			startLine = i + 1
			break
		}
	}
	if startLine == 0 {
		return nil
	}
	for i := startLine; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		splitLine := strings.Split(line, " ")

		iface := d.GetInterfaceByName(splitLine[1])
		if iface != nil {
			iface.Neighbor = splitLine[0]
			iface.NeighborPort = splitLine[len(splitLine)-1]
		}
	}
	return nil
}

func (d *CiscoDevice) GenInterfaceDescription() []string {
	var ret []string
	ret = append(ret, "en", "conf t")
	for _, iface := range d.Interfaces {
		if len(iface.Neighbor) > 0 {
			description := iface.Neighbor + "_" + iface.NeighborPort
			ret = append(ret, "interface "+iface.Name)
			ret = append(ret, "description "+description)
			ret = append(ret, "!")
		}
	}
	for _, c := range ret {
		fmt.Println(c)
	}
	return ret
}

func (d *CiscoDevice) SetInterfaceDescription() error {
	cmds := d.GenInterfaceDescription()
	_, err := d.SendCommands(cmds...)
	return err
}
