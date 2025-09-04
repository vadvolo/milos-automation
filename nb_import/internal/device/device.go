package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/annetutil/gnetcli/pkg/cmd"
	dcreds "github.com/annetutil/gnetcli/pkg/credentials"
	"github.com/annetutil/gnetcli/pkg/device/cisco"
	"github.com/annetutil/gnetcli/pkg/device/genericcli"
	"github.com/annetutil/gnetcli/pkg/device/huawei"
	"github.com/annetutil/gnetcli/pkg/streamer/ssh"
	"github.com/annetutil/gnetcli/pkg/streamer/telnet"
)

const (
	StatusActive  = "ACTIVE"
	StatusOffline = "OFFLINE"
)

//go:generate go run github.com/vektra/mockery/v2 --name=AbstractDevice
type AbstractDevice interface {
	GetHostname() string
	GetVendor() string
	GetAddress() string
	ShowInterfaces() []string
	ShowDeviceInfo()
	GetInterfaces() error
	GetLLDPNeigbours() error
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
	Breed      string       `json:"breed"`
	Interfaces []*Interface `json:"interfaces"`
	Active     bool         `json:"active"`
	Platform   string       `json:"platform"`
	Model      string       `json:"model"`
	Role       string       `json:"role"`
	Connector  *ssh.Streamer
	Logger     *zap.Logger
}

func NewDeivce(hostname, login, address, breed string) *Device {
	return &Device{
		Hostname: hostname,
		Login:    login,
		Address:  address,
		Breed:    breed,
	}
}

func (d *Device) GetHostname() string {
	return d.Hostname
}

func (d *Device) GetVendor() string {
	return d.Vendor
}

func (d *Device) GetAddress() string {
	return d.Address
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
	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
		dcreds.WithLogger(d.Logger),
	)
	return ssh.NewStreamer(d.Address, creds, ssh.WithLogger(d.Logger))
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
			fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return reses, nil
}

func (d *Device) GetStatus() bool {
	return d.Active
}

func (d *Device) SetStatus(s bool) {
	d.Active = s
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
	var command *exec.Cmd

	// Checking the type of OS because the ping command varies in structure according to the OS type
	if runtime.GOOS == "windows" {
		command = exec.Command("ping", "-n", "1", d.Address)
	} else {
		command = exec.Command("ping", "-c", "1", d.Address)
	}

	out, err := command.CombinedOutput()

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
