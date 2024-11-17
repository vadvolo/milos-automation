package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/annetutil/gnetcli/pkg/cmd"
	dcreds "github.com/annetutil/gnetcli/pkg/credentials"
	"github.com/annetutil/gnetcli/pkg/device/cisco"
	"github.com/annetutil/gnetcli/pkg/device/genericcli"
	"github.com/annetutil/gnetcli/pkg/streamer/ssh"
	"github.com/annetutil/gnetcli/pkg/streamer/telnet"
)

type NetworkDevice interface {
	_Address() string
	ShowRun() error
	Ping() (bool, error)
	SSHEnabled() (bool, error)
	SetSSH(keylength string) error
	WaitExec() error
}

type Device struct {
	Hostname  string `json:"hostname"`
	Domain    string `json:"domain"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	Vendor    string `json:"vendor"`
	Breed     string `json:"breed"`
	Protocol  string `json:"protocol"`
	Connector *ssh.Streamer
}

func NewDeivce(hostname, ipdomain, login, password, address, vendor, breed, protocol string) *Device {
	return &Device{
		Hostname: hostname,
		Domain:   ipdomain,
		Login:    login,
		Password: password,
		Address:  address,
		Vendor:   vendor,
		Breed:    breed,
		Protocol: protocol,
	}
}

func (d *Device) _Address() string {
	return d.Address
}

func (d *Device) WaitExec() error {
	count := 24
	for i := 0; i <= count; i++ {
		fmt.Printf("Trying to connect to device %s attempt %d\n", d._Address(), i)
		res, err := d.Ping()
		if err != nil {
			return err
		}
		if res {
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("Device %s is not available", d.Address)
}

func (d *Device) Ping() (bool, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", d.Address)
	} else {
		cmd = exec.Command("ping", "-c", "1", d.Address)
	}

	out, _ := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false, fmt.Errorf("there was an error pinging the host: %e", err)
	// }

	outStr := string(out)
	if strings.Contains(outStr, "Request timeout") || strings.Contains(outStr, "Destination Host Unreachable") || strings.Contains(outStr, "100% packet loss") {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *Device) SSHConnector() *ssh.Streamer {
	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
	)
	return ssh.NewStreamer(d.Address, creds)
}

func (d *Device) TelnetConnector() *telnet.Streamer {
	creds := dcreds.NewSimpleCredentials(
		dcreds.WithUsername(d.Login),
		dcreds.WithPassword(dcreds.Secret(d.Password)),
	)
	return telnet.NewStreamer(d.Address, creds)
}

func (d *Device) SendCommand(command cmd.Cmd) (cmd.CmdRes, error) {
	res, err := d.SendCommands(command)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("empty results")
	}
	return res[0], nil
}

func (d *Device) SendCommands(commands ...cmd.Cmd) ([]cmd.CmdRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dev := genericcli.GenericDevice{}
	switch d.Vendor {
	case "cisco":
		if d.Protocol == "telnet" {
			dev = cisco.NewDevice(d.TelnetConnector())
		} else {
			dev = cisco.NewDevice(d.SSHConnector())
		}
	default:
		return nil, errors.New("unknown vendor")
	}
	err := dev.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer dev.Close()
	reses, err := dev.ExecuteBulk(commands)
	if err != nil {
		return nil, err
	}
	for _, res := range reses {
		if res.Status() == 0 {
			fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return reses, nil
}
