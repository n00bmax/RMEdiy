package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var config Config
var out drivers.Out

type ParameterKey struct {
	Index int
	Name  string
}

var deviceStatus map[ParameterKey]int

type rmediator struct {
	IP       string `yaml:"ip" json:"ip"`
	Port     string `yaml:"port" json:"port"`
	Disabled bool   `yaml:"disabled" json:"disabled"`
}
type ADI struct {
	ID        DeviceID  `yaml:"id" json:"id"`
	InPort    int       `yaml:"midi_port_in" json:"midi_port_in"`
	OutPort   int       `yaml:"midi_port_out" json:"midi_port_out"`
	Name      string    `yaml:"name" json:"name"`
	Rmediator rmediator `yaml:"server" json:"server"`
}

type Config struct {
	ADI  `yaml:"device" json:"device"`
	Sync struct {
		Interval int `yaml:"interval"`
		RGB      struct {
			Enabled     bool `yaml:"enabled"`
			RefreshRate int  `yaml:"refreshRate"`
		} `yaml:"rgb"`
	} `yaml:"sync"`
}

func getConfig() error {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("Failed to read the YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal YAML content: %v", err)
	}

	dev, found := ADIDevices[config.ID]
	if !found {
		return fmt.Errorf("invalid device id: %v; must be one of %v[%v],%v[%v],%v[%v]", config.ID, ADI2DAC, ADIDevices[ADI2DAC], ADI2Pro, ADIDevices[ADI2Pro], ADI24, ADIDevices[ADI24])

	}
	switch config.ID {
	case ADI2DAC, ADI24, ADI2Pro:
	default:
	}
	config.Name = ADIDevices[config.ID]
	klog.Infof("Device %v:%v", dev, config.ID)
	klog.Info("Selected port", config.Rmediator.Port)
	klog.Info("Sync Interval:", config.Sync.Interval)
	klog.Info("RGB Enabled:", config.Sync.RGB.Enabled)
	klog.Info("RGB Refresh Rate:", config.Sync.RGB.RefreshRate)
	return nil
}
