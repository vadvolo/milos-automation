package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/annetutil/gnetcli/pkg/cmd"
	dcreds "github.com/annetutil/gnetcli/pkg/credentials"
	"github.com/annetutil/gnetcli/pkg/device/cisco"
	"github.com/annetutil/gnetcli/pkg/device/genericcli"
	"github.com/annetutil/gnetcli/pkg/device/huawei"
	"github.com/annetutil/gnetcli/pkg/streamer/ssh"
	"github.com/annetutil/gnetcli/pkg/streamer/telnet"
	"go.uber.org/zap"
)

type AbstractDevice interface {
	_Hostname() string
	_Vendor() string
	_Model() string
	_Address() string
	_Interfaces() []*Interface
	ShowInterfaces() []string
	ShowDeviceInfo()
	GetSoftware() error
	GetInterfaces() error
	GetNeigbours() error
	SetInterfaceDescription() error
	Ping() error
	GetStatus() bool
	SetStatus(s bool)
}

type Device struct {
	Hostname   string       `json:"hostname"`
	Login      string       `json:"login"`
	Password   string       `json:"password"`
	Address    string       `json:"address"`
	Vendor     string       `json:"vendor"`
	Model      string       `json:"model"`
	Breed      string       `json:"breed"`
	Interfaces []*Interface `json:"interfaces"`
	Active     bool         `json:"active"`
	Connector  *ssh.Streamer
}

func NewDeivce(hostname, login, address, breed string) *Device {
	return &Device{
		Hostname: hostname,
		Login:    login,
		Address:  address,
		Breed:    breed,
	}
}

func (d *Device) _Hostname() string {
	return d.Hostname
}

func (d *Device) _Vendor() string {
	return d.Vendor
}

func (d *Device) _Model() string {
	return d.Model
}

func (d *Device) _Address() string {
	return d.Address
}

func (d *Device) _Interfaces() []*Interface {
	return d.Interfaces
}

func (d *Device) ShowInterfaces() []string {
	var ret []string
	for _, iface := range d.Interfaces {
		ret = append(ret, iface.Name)
	}
	return ret
}

func (d *Device) ShowDeviceInfo() {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (d *Device) GetConnector() *ssh.Streamer {
	logger := zap.Must(zap.NewDevelopmentConfig().Build())
	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
		dcreds.WithLogger(logger),
	)
	return ssh.NewStreamer(d.Address, creds, ssh.WithLogger(logger))
}

func (d *Device) TelnetConnector() *telnet.Streamer {
	logger := zap.Must(zap.NewDevelopmentConfig().Build())
	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
		dcreds.WithLogger(logger),
	)
	return telnet.NewStreamer(d.Address, creds, telnet.WithLogger(logger))
}

func (d *Device) SendCommand(command string) (cmd.CmdRes, error) {
	res, err := d.SendCommands(command)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("empty results")
	}
	return res[0], nil
}

func (d *Device) SendCommands(commands ...string) ([]cmd.CmdRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dev := genericcli.GenericDevice{}
	switch d.Vendor {
	case "Cisco":
		dev = cisco.NewDevice(d.GetConnector())
	case "Huawei":
		dev = huawei.NewDevice(d.GetConnector())
	default:
		return nil, errors.New("unknown vendor")
	}
	err := dev.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer dev.Close()
	reses, _ := dev.ExecuteBulk(cmd.NewCmdList(commands))
	for _, res := range reses {
		if res.Status() == 0 {
			// fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return reses, nil
}

func ImportDevices() ([]AbstractDevice, error) {
	file, _ := os.ReadFile("inventory.json")
	devices := []*Device{}
	err := json.Unmarshal([]byte(file), &devices)
	if err != nil {
		return nil, err
	}

	var ret []AbstractDevice
	for _, device := range devices {
		if device.Vendor == "Cisco" {
			ret = append(ret, &CiscoDevice{
				Device: device,
			})
		}
		if device.Vendor == "Huawei" {
			ret = append(ret, &HuaweiDevice{
				Device: device,
			})
		}
	}

	return ret, nil
}

func (d *Device) GetInterfaceByName(name string) *Interface {
	for _, iface := range d.Interfaces {
		if iface.Name == name {
			return iface
		}
		if iface.ShortName == name {
			return iface
		}
	}
	return nil
}

func (d *Device) Ping() error {
	var cmd *exec.Cmd

	// Checking the type of OS because the ping command varies in structure according to the OS type
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", d.Address)
	} else {
		cmd = exec.Command("ping", "-c", "1", d.Address)
	}

	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("there was an error pinging the host: %e", err)
	}

	outStr := string(out)
	if strings.Contains(outStr, "Request timeout") || strings.Contains(outStr, "Destination Host Unreachable") || strings.Contains(outStr, "100% packet loss") {
		return fmt.Errorf("the host is not reachable")
	} else {
		return nil
	}
}

func (d *Device) GetStatus() bool {
	return d.Active
}

func (d *Device) SetStatus(s bool) {
	d.Active = s
}

type IfaceNeighbor struct {
	Neighbor     string `json:"neighbor"`
	NeighborPort string `json:"neighbor_port"`
}

type IfaceStatus struct {
	Operational    string `json:"operational"`
	Administrative string `json:"administrative"`
}

type Interface struct {
	Name         string          `json:"name"`
	ShortName    string          `json:"shortname"`
	IfaceType    string          `json:"type"`
	Description  string          `json:"description"`
	IPAddress    string          `json:"ipaddress"`
	Status       IfaceStatus     `json:"status"`
	Neighbors    []IfaceNeighbor `json:"neighbors"`
	Neighbor     string          `json:"neighbor"`
	NeighborPort string          `json:"neighbor_port"`
}

type CiscoDevice struct {
	*Device
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
	if strings.Contains(name, "Fas") {
		r := regexp.MustCompile(`Fas`)
		return r.ReplaceAllString(name, "Fa")
	}
	if strings.Contains(name, "Gig") {
		r := regexp.MustCompile(`Gig`)
		return r.ReplaceAllString(name, "Gi")
	}
	return ""
}

func (d *CiscoDevice) GetIfaceType(name string) string {
	if strings.Contains(name, "FastEthernet") {
		return "100base-tx"
	}
	if strings.Contains(name, "GigabitEthernet") {
		return "1000base-t"
	}
	return ""
}

func (d *CiscoDevice) GetSoftware() error {
	data, err := d.Device.SendCommand("show ver")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data.Output()))
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Cisco IOS Software") {
			start_line = i
			break
		}
	}
	model := strings.Split(txtlines[start_line], " ")[3]
	d.Model = model
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
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Interface") {
			start_line = i + 1
			break
		}
	}
	if start_line == 0 {
		return nil
	}
	for i := start_line; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		split_line := strings.Split(line, " ")
		d.Interfaces = append(d.Interfaces, &Interface{
			Name:      split_line[0],
			ShortName: d.CutIfaceName(split_line[0]),
			IfaceType: d.GetIfaceType(split_line[0]),
			IPAddress: split_line[1],
			Status: IfaceStatus{
				Operational:    split_line[5],
				Administrative: split_line[4],
			},
		})
	}
	return nil
}

func (d *CiscoDevice) GetNeigbours() error {
	// err := d.GetLLDPNeigbours()
	// if err != nil {
	// 	return err
	// }
	return d.GetCDPNeigbours()
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
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Device ID") {
			start_line = i + 1
			break
		}
	}
	if start_line == 0 {
		return nil
	}
	for i := start_line; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		split_line := strings.Split(line, " ")

		iface := d.GetInterfaceByName(split_line[1])
		if iface != nil {
			iface.Neighbor = split_line[0]
			iface.NeighborPort = split_line[len(split_line)-1]
		}
	}
	return nil
}

func (d *CiscoDevice) GetCDPNeigbours() error {
	data, err := d.Device.SendCommand("show cdp neighbors")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(data.Output()))
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Device ID") {
			start_line = i + 1
			break
		}
	}
	if start_line == 0 {
		return nil
	}
	for i := start_line; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		split_line := strings.Split(line, " ")

		ifaceName := d.CutIfaceName(split_line[1] + split_line[2])

		fmt.Println(ifaceName)

		iface := d.GetInterfaceByName(ifaceName)
		var ifaceNeighbour IfaceNeighbor
		if iface != nil {
			ifaceNeighbour.Neighbor = split_line[0]
			ifaceNameNeighbour := d.CutIfaceName(split_line[len(split_line)-2] + split_line[len(split_line)-1])
			ifaceNeighbour.NeighborPort = ifaceNameNeighbour
			iface.Neighbors = append(iface.Neighbors, ifaceNeighbour)
		}
	}
	return nil
}

func (d *CiscoDevice) GenInterfaceDescription() []string {
	ret := []string{}
	ret = append(ret, "en", "conf t")
	for _, iface := range d.Interfaces {
		if len(iface.Neighbors) > 0 {
			description := new(strings.Builder)
			for _, neighbor := range iface.Neighbors {
				description.WriteString(neighbor.Neighbor + "_" + neighbor.NeighborPort + ",")
			}
			ret = append(ret, "interface "+iface.Name)
			ret = append(ret, "description "+description.String()[:len(description.String())-1])
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

type HuaweiDevice struct {
	*Device
}

func (d *HuaweiDevice) CutIfaceName(name string) string {
	if strings.Contains(name, "FastEthernet") {
		r := regexp.MustCompile(`FastEthernet`)
		return r.ReplaceAllString(name, "Fa")
	}
	return ""
}

func (d *HuaweiDevice) GetSoftware() error {
	return nil
}

func (d *HuaweiDevice) GetInterfaces() error {
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
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Interface") {
			start_line = i + 1
			break
		}
	}
	if start_line == 0 {
		return nil
	}
	for i := start_line; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		split_line := strings.Split(line, " ")
		d.Interfaces = append(d.Interfaces, &Interface{
			Name:      split_line[0],
			ShortName: d.CutIfaceName(split_line[0]),
		})
	}
	return nil
}

func (d *HuaweiDevice) GetNeigbours() error {
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
	start_line := 0
	for i, line := range txtlines {
		if strings.Contains(line, "Local Interface") {
			start_line = i + 2
			break
		}
	}
	if start_line == 0 {
		return nil
	}
	for i := start_line; i < len(txtlines); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(txtlines[i], " ")
		if len(line) == 0 {
			break
		}
		split_line := strings.Split(line, " ")

		iface := d.GetInterfaceByName(split_line[0])
		if iface != nil {
			iface.Neighbor = split_line[3]
			iface.NeighborPort = split_line[2]
		}
	}
	return nil
}

func (d *HuaweiDevice) GenInterfaceDescription() []string {
	ret := []string{}
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

func (d *HuaweiDevice) SetInterfaceDescription() error {
	cmds := d.GenInterfaceDescription()
	_, err := d.SendCommands(cmds...)
	return err
}
