package huawei

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"nb_import/internal/device"
)

type VRPDevice struct {
	*device.Device
}

func (d *VRPDevice) CutIfaceName(name string) string {
	if strings.Contains(name, "FastEthernet") {
		r := regexp.MustCompile(`FastEthernet`)
		return r.ReplaceAllString(name, "Fa")
	}
	return ""
}

func (d *VRPDevice) GetInterfaces() error {
	data, err := d.Device.SendCommand("display interface brief")
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

func (d *VRPDevice) GetLLDPNeigbours() error {
	data, err := d.Device.SendCommand("display lldp neighbor brief")
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
		if strings.Contains(line, "Local Interface") {
			startLine = i + 2
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

		iface := d.GetInterfaceByName(splitLine[0])
		if iface != nil {
			iface.Neighbor = splitLine[3]
			iface.NeighborPort = splitLine[2]
		}
	}
	return nil
}

func (d *VRPDevice) GenInterfaceDescription() []string {
	var ret []string
	ret = append(ret, "system-view")
	for _, iface := range d.Interfaces {
		if len(iface.Neighbor) > 0 {
			description := iface.Neighbor + "_" + iface.NeighborPort
			ret = append(ret, "interface "+iface.Name)
			ret = append(ret, "description "+description)
			ret = append(ret, "#")
		}
	}
	for _, c := range ret {
		fmt.Println(c)
	}
	return ret
}

func (d *VRPDevice) SetInterfaceDescription() error {
	cmds := d.GenInterfaceDescription()
	_, err := d.SendCommands(cmds...)
	return err
}
