package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/annetutil/gnetcli/pkg/cmd"
	dcreds "github.com/annetutil/gnetcli/pkg/credentials"
	"github.com/annetutil/gnetcli/pkg/device/cisco"
	"github.com/annetutil/gnetcli/pkg/device/huawei"
	"github.com/annetutil/gnetcli/pkg/streamer/ssh"
	"go.uber.org/zap"
)

type AbstractDevice interface {
	ShowDeviceInfo()
	GetInterfaces() error
	GetLLDPNeigbours() error
	SetInterfaceDescription() error
}

type Device struct {
	Hostname   string       `json:"hostname"`
	Login      string       `json:"login"`
	Password   string       `json:"password"`
	Address    string       `json:"address"`
	Vendor     string       `json:"vendor"`
	Breed      string       `json:"breed"`
	Interfaces []*Interface `json:"interfaces"`
}

func NewDeivce(hostname, login, address, breed string) *Device {
	return &Device{
		Hostname: hostname,
		Login:    login,
		Address:  address,
		Breed:    breed,
	}
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

type Interface struct {
	Name         string `json:"name"`
	ShortName    string `json:"shortname"`
	Description  string `json:"description"`
	Neighbor     string `json:"neighbor"`
	NeighborPort string `json:"neighbor_port"`
}

type CiscoDevice struct {
	*Device
}

func (d *CiscoDevice) SendCommands(commands ...string) ([]cmd.CmdRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := zap.Must(zap.NewDevelopmentConfig().Build())

	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
		dcreds.WithLogger(logger),
	)
	connector := ssh.NewStreamer(d.Address, creds, ssh.WithLogger(logger))
	dev := cisco.NewDevice(connector)
	err := dev.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer dev.Close()
	reses, _ := dev.ExecuteBulk(cmd.NewCmdList(commands))
	for _, res := range reses {
		if res.Status() == 0 {
			fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return reses, nil
}

func (d *CiscoDevice) ShowDeviceInfo() {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

func (d *CiscoDevice) CutIfaceName(name string) string {
	if strings.Contains(name, "FastEthernet") {
		r := regexp.MustCompile(`FastEthernet`)
		return r.ReplaceAllString(name, "Fa")
	}
	return ""
}

func (d *CiscoDevice) GetInterfaces() error {
	res, err := d.SendCommands("show ip interface brief")
	if err != nil {
		return err
	}
	data := res[0]
	reader := bytes.NewReader(data.Output())
	scanner := bufio.NewScanner(reader)
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

func (d *CiscoDevice) GetLLDPNeigbours() error {
	res, err := d.SendCommands("show lldp neighbors")
	if err != nil {
		return nil
	}
	data := res[0]
	reader := bytes.NewReader(data.Output())
	scanner := bufio.NewScanner(reader)
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
			iface.NeighborPort = split_line[4]
		}
	}
	return nil
}

func (d *CiscoDevice) GenInterfaceDescription() []string {
	ret := []string{}
	ret = append(ret, "en", "conf t")
	for _, iface := range d.Interfaces {
		if len(iface.Neighbor) > 0 {
			description := iface.Neighbor + "_" + iface.NeighborPort
			ret = append(ret, "interface " + iface.Name)
			ret = append(ret, "description " + description)
			ret = append(ret, "!")
			break
		}
	}
	for _, c := range ret {
		fmt.Println(c)
	}
	return ret
}

func (d *CiscoDevice) SetInterfaceDescription() error {
	cmds := d.GenInterfaceDescription()
	d.SendCommands(cmds...)
	return nil
}

type HuaweiDevice struct {
	*Device
}

func (d *HuaweiDevice) SendCommands(commands []string) ([]cmd.CmdRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	logger := zap.Must(zap.NewDevelopmentConfig().Build())

	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
		dcreds.WithLogger(logger),
	)
	connector := ssh.NewStreamer(d.Address, creds, ssh.WithLogger(logger))
	dev := huawei.NewDevice(connector)
	err := dev.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer dev.Close()
	reses, _ := dev.ExecuteBulk(cmd.NewCmdList(commands))
	for _, res := range reses {
		if res.Status() == 0 {
			fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return nil, nil
}

func (d *HuaweiDevice) ShowDeviceInfo() {
	fmt.Println(
		d.Address,
	)
}
func (d *HuaweiDevice) GetInterfaces() error {
	return nil
}

func (d *HuaweiDevice) GetLLDPNeigbours() error {
	return nil
}

func (d *HuaweiDevice) SetInterfaceDescription() error {
	return nil
}
