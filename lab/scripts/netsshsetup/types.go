package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/annetutil/gnetcli/pkg/cmd"
	dcreds "github.com/annetutil/gnetcli/pkg/credentials"
	"github.com/annetutil/gnetcli/pkg/device/cisco"
	"github.com/annetutil/gnetcli/pkg/device/genericcli"
	"github.com/annetutil/gnetcli/pkg/streamer/ssh"
	"github.com/annetutil/gnetcli/pkg/streamer/telnet"
	"go.uber.org/zap"
)

type NetworkDevice interface {
	ShowRun() error
	SSHEnabled() (bool, error)
	SetSSH() error
}

type Device struct {
	Hostname  string `json:"hostname"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	Vendor    string `json:"vendor"`
	Breed     string `json:"breed"`
	Protocol  string `json:"protocol"`
	Connector *ssh.Streamer
}

func NewDeivce(hostname, login, password, address, vendor, breed, protocol string) *Device {
	return &Device{
		Hostname: hostname,
		Login:    login,
		Password: password,
		Address:  address,
		Vendor:   vendor,
		Breed:    breed,
		Protocol: protocol,
	}
}

func (d *Device) SSHConnector() *ssh.Streamer {
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
	reses, _ := dev.ExecuteBulk(commands)
	for _, res := range reses {
		if res.Status() == 0 {
			fmt.Printf("Result: %s\n", res.Output())
		} else {
			fmt.Printf("Error: %d\nStatus: %d\n", res.Status(), res.Error())
		}
	}
	return reses, nil
}
